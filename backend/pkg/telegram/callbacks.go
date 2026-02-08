package telegram

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleCallback(query *tgbotapi.CallbackQuery) {
	chatID := query.Message.Chat.ID
	messageID := query.Message.MessageID
	data := query.Data

	// Parse callback data
	parts := strings.SplitN(data, ":", 2)
	if len(parts) != 2 {
		b.SendMessage(chatID, "❌ 无效的操作")
		return
	}

	action := parts[0]
	taskIDStr := parts[1]

	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		b.SendMessage(chatID, "❌ 无效的任务ID")
		return
	}

	// Find user by chat ID
	user, err := b.userModel.FindByTgChatID(chatID)
	if err != nil {
		b.SendMessage(chatID, "❌ 数据库查询失败")
		log.Printf("Error finding user: %v", err)
		return
	}

	if user == nil {
		b.SendMessage(chatID, "❌ 账号未绑定")
		return
	}

	switch action {
	case "complete":
		b.handleCompleteCallback(chatID, messageID, user.ID, taskID)
	case "delete":
		b.handleDeleteCallback(chatID, messageID, user.ID, taskID)
	default:
		b.SendMessage(chatID, "❌ 未知操作")
	}

	// Answer callback query (dismiss the loading state)
	answerCfg := tgbotapi.NewCallback(query.ID, "")
	b.api.Request(answerCfg)
}

func (b *Bot) handleCompleteCallback(chatID int64, messageID int, userID int64, taskID int64) {
	// Get task to verify ownership
	task, err := b.taskModel.FindByID(taskID)
	if err != nil {
		b.SendMessage(chatID, "❌ 查询任务失败")
		log.Printf("Error finding task: %v", err)
		return
	}

	if task == nil {
		b.SendMessage(chatID, "❌ 任务不存在")
		return
	}

	if task.UserID != userID {
		b.SendMessage(chatID, "❌ 无权操作此任务")
		return
	}

	// Complete task using the task completer interface
	if b.taskCompleter == nil {
		b.SendMessage(chatID, "❌ 系统未就绪")
		log.Printf("Task completer not set")
		return
	}

	expGained, goldGained, newLevel, newExp, err := b.taskCompleter.CompleteTask(userID, taskID)
	if err != nil {
		b.SendMessage(chatID, fmt.Sprintf("❌ 完成失败：%s", err.Error()))
		log.Printf("Error completing task: %v", err)
		return
	}

	// Send success message
	msg := fmt.Sprintf("✅ 任务「%s」已完成！\n获得 %d经验 %d金币\n\n当前等级：%d | 经验：%d",
		task.Title, expGained, goldGained, newLevel, newExp)
	b.SendMessage(chatID, msg)

	// Update the original task list message to reflect completion
	b.refreshTaskListMessage(chatID, messageID, userID)
}

func (b *Bot) handleDeleteCallback(chatID int64, messageID int, userID int64, taskID int64) {
	// Get task to verify ownership
	task, err := b.taskModel.FindByID(taskID)
	if err != nil {
		b.SendMessage(chatID, "❌ 查询任务失败")
		log.Printf("Error finding task: %v", err)
		return
	}

	if task == nil {
		b.SendMessage(chatID, "❌ 任务不存在")
		return
	}

	if task.UserID != userID {
		b.SendMessage(chatID, "❌ 无权操作此任务")
		return
	}

	// Delete task using the task completer interface
	if b.taskCompleter == nil {
		b.SendMessage(chatID, "❌ 系统未就绪")
		log.Printf("Task completer not set")
		return
	}

	if err := b.taskCompleter.DeleteTask(userID, taskID); err != nil {
		b.SendMessage(chatID, fmt.Sprintf("❌ 删除失败：%s", err.Error()))
		log.Printf("Error deleting task: %v", err)
		return
	}

	b.SendMessage(chatID, fmt.Sprintf("🗑 任务「%s」已删除", task.Title))

	// Update the original task list message to reflect deletion
	b.refreshTaskListMessage(chatID, messageID, userID)
}

// refreshTaskListMessage edits the original /tasks message to show updated task list
func (b *Bot) refreshTaskListMessage(chatID int64, messageID int, userID int64) {
	tasks, err := b.taskModel.FindByUserID(userID, "", "active")
	if err != nil {
		log.Printf("Error refreshing task list: %v", err)
		return
	}

	if len(tasks) == 0 {
		// No more active tasks, update message
		edit := tgbotapi.NewEditMessageText(chatID, messageID, "📋 任务列表：\n\n📭 所有任务已完成！")
		edit.ParseMode = tgbotapi.ModeHTML
		b.api.Send(edit)
		return
	}

	// Rebuild message with remaining tasks
	text := "📋 你的任务列表：\n\n"
	var keyboard [][]tgbotapi.InlineKeyboardButton

	for _, task := range tasks {
		text += fmt.Sprintf("• %s\n", task.Title)
		if task.Description != "" {
			text += fmt.Sprintf("  📝 %s\n", task.Description)
		}

		completeBtn := tgbotapi.NewInlineKeyboardButtonData("✅ 完成", fmt.Sprintf("complete:%d", task.ID))
		deleteBtn := tgbotapi.NewInlineKeyboardButtonData("🗑 删除", fmt.Sprintf("delete:%d", task.ID))
		keyboard = append(keyboard, tgbotapi.NewInlineKeyboardRow(completeBtn, deleteBtn))
		text += "\n"
	}

	markup := tgbotapi.NewInlineKeyboardMarkup(keyboard...)
	edit := tgbotapi.NewEditMessageTextAndMarkup(chatID, messageID, text, markup)
	edit.ParseMode = tgbotapi.ModeHTML
	b.api.Send(edit)
}
