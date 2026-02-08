package scheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"life-system-backend/internal/logic"
	"life-system-backend/internal/model"
	"life-system-backend/internal/svc"
	"life-system-backend/pkg/telegram"
)

type Scheduler struct {
	bot           *telegram.Bot
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
			s.checkDailyEnergyRestore()
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

		// Skip expired tasks (handled by checkExpiredChallengeTasks for challenge tasks)
		if now.After(deadline) {
			continue
		}

		// Check if we should send reminder
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

			// Create inline keyboard for quick actions
			completeBtn := tgbotapi.NewInlineKeyboardButtonData("✅ 完成", fmt.Sprintf("complete:%d", task.ID))
			deleteBtn := tgbotapi.NewInlineKeyboardButtonData("🗑 删除", fmt.Sprintf("delete:%d", task.ID))
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(completeBtn, deleteBtn),
			)

			if err := s.bot.SendMessageWithKeyboard(chatID, message, keyboard); err != nil {
				log.Printf("Error sending reminder: %v", err)
				continue
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

	// Only reset once per day
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

// checkDailyEnergyRestore restores energy to max for all characters at the start of each day
func (s *Scheduler) checkDailyEnergyRestore() {
	today := time.Now().Format("2006-01-02")

	// Only restore once per day (uses same lastResetDate as checkDailyReset)
	if s.lastResetDate == today {
		return
	}

	log.Printf("🔋 Daily energy restore is handled through sleep records")
	// Energy is primarily restored through sleep records
	// Players should record their sleep to restore energy
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

		// Send notification to user if bot is available and user has Telegram bound
		if s.bot != nil {
			user, err := s.svcCtx.UserModel.FindByID(task.UserID)
			if err == nil && user != nil && user.TgChatID > 0 {
				message := fmt.Sprintf("❌ 挑战任务「%s」已超过截止时间，自动失败！\n惩罚：-%d经验 -%d金币",
					task.Title, task.PenaltyExp, task.PenaltyGold)
				if err := s.bot.SendMessage(user.TgChatID, message); err != nil {
					log.Printf("Error sending failure notification: %v", err)
				}
			}
		}
	}
}

// checkAttributeDecay applies attribute decay for inactive characters
func (s *Scheduler) checkAttributeDecay() {
	// Find characters inactive for 1+ days
	characters, err := s.charModel.FindInactiveCharacters(1)
	if err != nil {
		log.Printf("Error finding inactive characters: %v", err)
		return
	}

	if len(characters) == 0 {
		return
	}

	today := time.Now().Format("2006-01-02")
	const minAttribute = 5.0
	const decayRatePerDay = 0.01 // 1% decay per day

	for _, stats := range characters {
		// Calculate days of inactivity
		lastActivity, err := time.Parse("2006-01-02", stats.LastActivityDate)
		if err != nil {
			log.Printf("Error parsing last activity date for user %d: %v", stats.UserID, err)
			continue
		}

		now := time.Now()
		daysInactive := int(now.Sub(lastActivity).Hours() / 24)

		if daysInactive < 1 {
			continue // Safety check
		}

		// Apply decay (1% per day)
		decayMultiplier := 1.0 - (float64(daysInactive) * decayRatePerDay)
		if decayMultiplier < 0.5 {
			decayMultiplier = 0.5 // Max 50% decay
		}

		originalStats := *stats

		stats.Strength = maxFloat(minAttribute, stats.Strength*decayMultiplier)
		stats.Intelligence = maxFloat(minAttribute, stats.Intelligence*decayMultiplier)
		stats.Vitality = maxFloat(minAttribute, stats.Vitality*decayMultiplier)
		stats.Spirit = maxFloat(minAttribute, stats.Spirit*decayMultiplier)

		// Recalculate MaxHP based on new attributes
		stats.MaxHP = 100 + int(stats.Strength*2) + int(stats.Vitality*3)
		if stats.HP > stats.MaxHP {
			stats.HP = stats.MaxHP
		}

		// Update last activity date to prevent repeated decay
		stats.LastActivityDate = today

		// Update character
		if err := s.charModel.Update(stats); err != nil {
			log.Printf("Error updating character stats for user %d: %v", stats.UserID, err)
			continue
		}

		log.Printf("⚠️  Applied decay to user %d after %d days of inactivity: Str %.1f→%.1f, Int %.1f→%.1f, Vit %.1f→%.1f, Spr %.1f→%.1f",
			stats.UserID, daysInactive,
			originalStats.Strength, stats.Strength,
			originalStats.Intelligence, stats.Intelligence,
			originalStats.Vitality, stats.Vitality,
			originalStats.Spirit, stats.Spirit)

		// Send notification to user if bot is available and user has Telegram bound
		if s.bot != nil {
			user, err := s.svcCtx.UserModel.FindByID(stats.UserID)
			if err == nil && user != nil && user.TgChatID > 0 {
				message := fmt.Sprintf("⚠️ 由于 %d 天未活动，你的属性发生了衰减！\n"+
					"💪 力量: %.1f → %.1f\n"+
					"🧠 智力: %.1f → %.1f\n"+
					"❤️ 体力: %.1f → %.1f\n"+
					"✨ 精神: %.1f → %.1f\n"+
					"\n完成任务来恢复和提升属性吧！",
					daysInactive,
					originalStats.Strength, stats.Strength,
					originalStats.Intelligence, stats.Intelligence,
					originalStats.Vitality, stats.Vitality,
					originalStats.Spirit, stats.Spirit)

				if err := s.bot.SendMessage(user.TgChatID, message); err != nil {
					log.Printf("Error sending decay notification: %v", err)
				}
			}
		}
	}
}

// maxFloat returns the maximum of two float64 values
func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
