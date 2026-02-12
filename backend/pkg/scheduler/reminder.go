package scheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"life-system-backend/internal/logic"
	"life-system-backend/internal/model"
	"life-system-backend/internal/realm"
	"life-system-backend/internal/svc"
	"life-system-backend/pkg/bark"
	"life-system-backend/pkg/telegram"
)

type Scheduler struct {
	bot           *telegram.Bot
	barkClient    *bark.Client
	taskModel     *model.TaskModel
	charModel     *model.CharacterModel
	svcCtx        *svc.ServiceContext
	interval      time.Duration
	stop          chan struct{}
	running       bool
	lastResetDate string
}

func NewScheduler(bot *telegram.Bot, svcCtx *svc.ServiceContext, interval time.Duration) *Scheduler {
	if interval == 0 {
		interval = 1 * time.Minute
	}

	return &Scheduler{
		bot:           bot,
		barkClient:    svcCtx.BarkClient,
		taskModel:     svcCtx.TaskModel,
		charModel:     svcCtx.CharacterModel,
		svcCtx:        svcCtx,
		interval:      interval,
		stop:          make(chan struct{}),
		running:       false,
		lastResetDate: "",
	}
}

func (s *Scheduler) Start() {
	if s.running {
		return
	}

	s.running = true
	go s.run()
}

func (s *Scheduler) run() {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-s.stop:
			return
		case <-ticker.C:
			s.checkDailyReset()
			s.checkAttributeDecay()
			s.checkExpiredChallengeTasks()
			s.checkTasks()
		}
	}
}

func (s *Scheduler) Stop() {
	if s.running {
		s.running = false
		close(s.stop)
	}
}

func (s *Scheduler) checkTasks() {
	tasksWithUsers, err := s.taskModel.FindTasksNeedingReminder()
	if err != nil {
		log.Printf("Error finding tasks needing reminder: %v", err)
		return
	}

	now := time.Now()

	for _, tw := range tasksWithUsers {
		task := tw.Task
		chatID := tw.TgChatID

		if !task.Deadline.Valid {
			continue
		}

		deadline := task.Deadline.Time
		reminderTime := deadline.Add(-time.Duration(task.RemindBefore) * time.Minute)

		// Skip expired tasks
		if now.After(deadline) {
			continue
		}

		shouldRemind := false

		if !task.LastRemindedAt.Valid {
			if now.After(reminderTime) || now.Equal(reminderTime) {
				shouldRemind = true
			}
		} else {
			if task.RemindInterval > 0 {
				lastReminded := task.LastRemindedAt.Time
				nextReminder := lastReminded.Add(time.Duration(task.RemindInterval) * time.Minute)
				if now.After(nextReminder) || now.Equal(nextReminder) {
					shouldRemind = true
				}
			}
		}

		if shouldRemind {
			remaining := deadline.Sub(now)
			var remainingStr string

			hours := int(remaining.Hours())
			minutes := int(remaining.Minutes()) % 60

			if hours > 0 {
				remainingStr = fmt.Sprintf("%d小时%d分钟", hours, minutes)
			} else {
				remainingStr = fmt.Sprintf("%d分钟", minutes)
			}

			description := ""
			if task.Description != "" {
				description = fmt.Sprintf("\n%s", task.Description)
			}

			message := fmt.Sprintf("⏰ 提醒：任务「%s」还剩 %s 到期！%s",
				task.Title, remainingStr, description)

			// Send Telegram notification
			if s.bot != nil && chatID > 0 {
				completeBtn := tgbotapi.NewInlineKeyboardButtonData("✅ 完成", fmt.Sprintf("complete:%d", task.ID))
				deleteBtn := tgbotapi.NewInlineKeyboardButtonData("🗑 删除", fmt.Sprintf("delete:%d", task.ID))
				keyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(completeBtn, deleteBtn),
				)

				if err := s.bot.SendMessageWithKeyboard(chatID, message, keyboard); err != nil {
					log.Printf("Error sending Telegram reminder: %v", err)
				}
			}

			// Send Bark notification
			if s.barkClient != nil {
				user, err := s.svcCtx.UserModel.FindByID(task.UserID)
				if err == nil && user != nil && user.BarkKey != "" {
					barkTitle := fmt.Sprintf("⏰ 任务提醒 - 还剩%s", remainingStr)
					barkBody := task.Title
					if task.Description != "" {
						barkBody = fmt.Sprintf("%s\n%s", task.Title, task.Description)
					}
					if err := s.barkClient.PushAlarm(user.BarkKey, barkTitle, barkBody); err != nil {
						log.Printf("Error sending Bark reminder: %v", err)
					}
				}
			}

			if err := s.taskModel.UpdateLastReminded(task.ID, now); err != nil {
				log.Printf("Error updating last reminded: %v", err)
			}
		}
	}
}

// checkDailyReset resets daily completion counts for repeatable tasks at the start of each day
func (s *Scheduler) checkDailyReset() {
	today := time.Now().Format("2006-01-02")

	if s.lastResetDate == today {
		return
	}

	log.Printf("🔄 Starting daily reset for %s", today)
	if err := s.taskModel.ResetDailyCompletionCounts(today); err != nil {
		log.Printf("Error resetting daily completion counts: %v", err)
		return
	}

	s.lastResetDate = today
	log.Printf("✅ Daily reset completed for %s", today)
}

// checkExpiredChallengeTasks finds expired challenge tasks and applies penalties
func (s *Scheduler) checkExpiredChallengeTasks() {
	tasks, err := s.taskModel.FindExpiredChallengeTasks()
	if err != nil {
		log.Printf("Error finding expired challenge tasks: %v", err)
		return
	}

	if len(tasks) == 0 {
		return
	}

	log.Printf("⚠️  Found %d expired challenge task(s)", len(tasks))

	taskLogic := logic.NewTaskLogic(s.svcCtx)

	for _, task := range tasks {
		reason := fmt.Sprintf("超过截止时间 %s", task.Deadline.Time.Format("2006-01-02 15:04"))

		if err := taskLogic.FailTask(context.Background(), task.ID, reason); err != nil {
			log.Printf("Error failing task #%d: %v", task.ID, err)
			continue
		}

		log.Printf("✅ Task #%d '%s' marked as failed with penalties applied", task.ID, task.Title)

		// Send notification
		if s.bot != nil {
			user, err := s.svcCtx.UserModel.FindByID(task.UserID)
			if err == nil && user != nil && user.TgChatID > 0 {
				message := fmt.Sprintf("❌ 挑战任务「%s」已超过截止时间，自动失败！\n惩罚：-%d灵石",
					task.Title, task.PenaltySpiritStones)
				if err := s.bot.SendMessage(user.TgChatID, message); err != nil {
					log.Printf("Error sending failure notification: %v", err)
				}
			}
		}
	}
}

// checkAttributeDecay applies attribute decay for inactive characters
func (s *Scheduler) checkAttributeDecay() {
	characters, err := s.charModel.FindInactiveCharacters(1)
	if err != nil {
		log.Printf("Error finding inactive characters: %v", err)
		return
	}

	if len(characters) == 0 {
		return
	}

	today := time.Now().Format("2006-01-02")
	const decayRatePerDay = 0.01 // 1% decay per day

	for _, stats := range characters {
		lastActivity, err := time.Parse("2006-01-02", stats.LastActivityDate)
		if err != nil {
			log.Printf("Error parsing last activity date for user %d: %v", stats.UserID, err)
			continue
		}

		now := time.Now()
		daysInactive := int(now.Sub(lastActivity).Hours() / 24)

		if daysInactive < 1 {
			continue
		}

		decayMultiplier := 1.0 - (float64(daysInactive) * decayRatePerDay)
		if decayMultiplier < 0.5 {
			decayMultiplier = 0.5
		}

		// Load attributes and apply decay
		attrs, err := s.charModel.FindAttributesByUserID(stats.UserID)
		if err != nil {
			log.Printf("Error finding attributes for user %d: %v", stats.UserID, err)
			continue
		}

		for _, attr := range attrs {
			if attr.AttrKey == "luck" {
				continue
			}

			minVal := realm.AttrMin(attr.Realm)
			newValue := attr.Value * decayMultiplier
			if newValue < minVal {
				newValue = minVal
			}

			if newValue != attr.Value {
				attr.Value = newValue
				if err := s.charModel.UpdateAttribute(attr); err != nil {
					log.Printf("Error updating attribute %s for user %d: %v", attr.AttrKey, stats.UserID, err)
				}
			}
		}

		// Update last activity date to prevent repeated decay
		stats.LastActivityDate = today
		if err := s.charModel.Update(stats); err != nil {
			log.Printf("Error updating character stats for user %d: %v", stats.UserID, err)
			continue
		}

		log.Printf("⚠️  Applied attribute decay to user %d after %d days of inactivity", stats.UserID, daysInactive)

		// Send notification
		if s.bot != nil {
			user, err := s.svcCtx.UserModel.FindByID(stats.UserID)
			if err == nil && user != nil && user.TgChatID > 0 {
				message := fmt.Sprintf("⚠️ 由于 %d 天未活动，你的属性发生了衰减！\n完成任务来恢复和提升属性吧！",
					daysInactive)
				if err := s.bot.SendMessage(user.TgChatID, message); err != nil {
					log.Printf("Error sending decay notification: %v", err)
				}
			}
		}
	}
}
