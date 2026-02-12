package telegram

import (
	"database/sql"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"life-system-backend/internal/model"
)

// TaskCompleter interface to avoid circular dependency with logic package
type TaskCompleter interface {
	CompleteTask(userID int64, taskID int64) (expGained int, spiritStonesGained int, realmTitle string, spiritStones int, err error)
	DeleteTask(userID int64, taskID int64) error
}

// ServiceContextInterface defines the interface for service context to avoid circular imports
type ServiceContextInterface interface {
	GetDB() *sql.DB
	GetUserModel() *model.UserModel
	GetTaskModel() *model.TaskModel
	GetCharacterModel() *model.CharacterModel
}

type Bot struct {
	api           *tgbotapi.BotAPI
	db            *sql.DB
	userModel     *model.UserModel
	taskModel     *model.TaskModel
	charModel     *model.CharacterModel
	svcCtx        ServiceContextInterface
	taskCompleter TaskCompleter
}

func NewBot(token string, db *sql.DB) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		api:       api,
		db:        db,
		userModel: model.NewUserModel(db),
		taskModel: model.NewTaskModel(db),
		charModel: model.NewCharacterModel(db),
	}

	return bot, nil
}

func (b *Bot) Start() {
	go b.pollUpdates()
}

func (b *Bot) pollUpdates() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			b.handleMessage(update.Message)
		} else if update.CallbackQuery != nil {
			b.handleCallback(update.CallbackQuery)
		}
	}
}

func (b *Bot) handleMessage(message *tgbotapi.Message) {
	chatID := message.Chat.ID

	if message.IsCommand() {
		b.handleCommand(chatID, message.Command(), message.CommandArguments())
	} else if message.Text != "" {
		// Try to handle plain text as a bind code
		b.handlePlainText(chatID, message.Text)
	}
}

func (b *Bot) handlePlainText(chatID int64, text string) {
	// Try to use the text as a bind code
	user, err := b.userModel.FindByBindCode(text)
	if err != nil {
		log.Printf("Error finding user by bind code: %v", err)
		return
	}

	if user != nil {
		// Found a valid bind code, perform binding
		username := ""
		chatConfig := tgbotapi.ChatInfoConfig{ChatConfig: tgbotapi.ChatConfig{ChatID: chatID}}
		if chat, err := b.api.GetChat(chatConfig); err == nil {
			username = chat.UserName
		}

		if err := b.userModel.UpdateTgBinding(user.ID, chatID, username); err != nil {
			b.SendMessage(chatID, "❌ 绑定失败，请重试。")
			log.Printf("Error binding user: %v", err)
			return
		}

		if err := b.userModel.ClearBindCode(user.ID); err != nil {
			log.Printf("Error clearing bind code: %v", err)
		}

		msg := fmt.Sprintf("✅ 绑定成功！你的账号 %s 已关联。\n使用 /tasks 查看任务，/help 查看帮助。", user.Username)
		b.SendMessage(chatID, msg)
	} else {
		b.SendMessage(chatID, "未知命令。使用 /help 查看帮助，或发送有效的绑定码进行账号绑定。")
	}
}

func (b *Bot) handleCommand(chatID int64, command, args string) {
	switch command {
	case "start":
		if args != "" {
			b.handleStartWithCode(chatID, args)
		} else {
			b.handleStartNoCode(chatID)
		}
	case "tasks":
		b.handleTasks(chatID)
	case "help":
		b.handleHelp(chatID)
	default:
		b.SendMessage(chatID, "未知命令。使用 /help 查看帮助。")
	}
}

func (b *Bot) handleStartWithCode(chatID int64, code string) {
	// Find user by bind code
	user, err := b.userModel.FindByBindCode(code)
	if err != nil {
		b.SendMessage(chatID, "错误：数据库查询失败")
		log.Printf("Error finding user by bind code: %v", err)
		return
	}

	if user == nil {
		b.SendMessage(chatID, "❌ 绑定码无效或已过期。请重新生成。")
		return
	}

	// Bind user
	username := ""
	chatConfig := tgbotapi.ChatInfoConfig{ChatConfig: tgbotapi.ChatConfig{ChatID: chatID}}
	if chat, err := b.api.GetChat(chatConfig); err == nil {
		username = chat.UserName
	}

	if err := b.userModel.UpdateTgBinding(user.ID, chatID, username); err != nil {
		b.SendMessage(chatID, "❌ 绑定失败，请重试。")
		log.Printf("Error binding user: %v", err)
		return
	}

	// Clear bind code
	if err := b.userModel.ClearBindCode(user.ID); err != nil {
		log.Printf("Error clearing bind code: %v", err)
	}

	message := fmt.Sprintf("✅ 绑定成功！你的账号 %s 已关联。\n使用 /tasks 查看任务，/help 查看帮助。", user.Username)
	b.SendMessage(chatID, message)
}

func (b *Bot) handleStartNoCode(chatID int64) {
	message := `👋 欢迎使用 Life System RPG！

这是一个游戏化的任务管理应用，可以帮助你提升生活质量。

要绑定账号，请：
1. 访问 Web 应用
2. 进入设置页面
3. 生成绑定码
4. 使用 /start <绑定码> 命令

绑定后，你可以：
- 📋 使用 /tasks 查看任务
- ✅ 点击按钮完成任务
- 🗑 删除任务

使用 /help 查看完整帮助。`

	b.SendMessage(chatID, message)
}

func (b *Bot) handleTasks(chatID int64) {
	// Find user by chat ID
	user, err := b.userModel.FindByTgChatID(chatID)
	if err != nil {
		b.SendMessage(chatID, "❌ 数据库查询失败")
		log.Printf("Error finding user by chat ID: %v", err)
		return
	}

	if user == nil {
		b.SendMessage(chatID, "❌ 账号未绑定。请使用 /start <绑定码> 进行绑定。")
		return
	}

	// Get active tasks
	tasks, err := b.taskModel.FindByUserID(user.ID, "", "active")
	if err != nil {
		b.SendMessage(chatID, "❌ 获取任务失败")
		log.Printf("Error getting tasks: %v", err)
		return
	}

	if len(tasks) == 0 {
		b.SendMessage(chatID, "📭 没有活跃任务。")
		return
	}

	// Build message with inline keyboard
	message := "📋 你的任务列表：\n\n"
	var keyboard [][]tgbotapi.InlineKeyboardButton

	for _, task := range tasks {
		message += fmt.Sprintf("• %s\n", task.Title)
		if task.Description != "" {
			message += fmt.Sprintf("  📝 %s\n", task.Description)
		}

		completeBtn := tgbotapi.NewInlineKeyboardButtonData("✅ 完成", fmt.Sprintf("complete:%d", task.ID))
		deleteBtn := tgbotapi.NewInlineKeyboardButtonData("🗑 删除", fmt.Sprintf("delete:%d", task.ID))
		keyboard = append(keyboard, tgbotapi.NewInlineKeyboardRow(completeBtn, deleteBtn))
		message += "\n"
	}

	b.SendMessageWithKeyboard(chatID, message, tgbotapi.NewInlineKeyboardMarkup(keyboard...))
}

func (b *Bot) handleHelp(chatID int64) {
	message := `🆘 帮助菜单

可用命令：
/start <绑定码> - 使用绑定码进行账号绑定
/tasks - 查看你的所有任务
/help - 显示此帮助信息

按钮操作：
✅ 完成 - 标记任务为已完成，获得奖励
🗑 删除 - 删除任务

需要更多帮助，请访问 Web 应用设置。`

	b.SendMessage(chatID, message)
}

func (b *Bot) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeHTML

	_, err := b.api.Send(msg)
	return err
}

func (b *Bot) SendMessageWithKeyboard(chatID int64, text string, keyboard tgbotapi.InlineKeyboardMarkup) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = keyboard

	_, err := b.api.Send(msg)
	return err
}

func (b *Bot) GetUsername() string {
	return b.api.Self.UserName
}

func (b *Bot) Stop() {
	b.api.StopReceivingUpdates()
}

// SetServiceContext sets the service context after bot initialization
func (b *Bot) SetServiceContext(svcCtx ServiceContextInterface) {
	b.svcCtx = svcCtx
}

// SetTaskCompleter sets the task completer to avoid circular dependency
func (b *Bot) SetTaskCompleter(completer TaskCompleter) {
	b.taskCompleter = completer
}
