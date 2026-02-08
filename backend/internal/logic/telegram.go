package logic

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"life-system-backend/internal/svc"
	"life-system-backend/internal/types"
)

type TelegramLogic struct {
	svcCtx *svc.ServiceContext
}

func NewTelegramLogic(svcCtx *svc.ServiceContext) *TelegramLogic {
	return &TelegramLogic{
		svcCtx: svcCtx,
	}
}

func (l *TelegramLogic) GenerateBindCode(ctx context.Context, userID int64) (*types.BindCodeResp, error) {
	// Generate random 6-char alphanumeric code
	code := generateRandomCode(6)

	// Set expiry to 5 minutes from now (UTC for SQLite compatibility)
	expire := time.Now().UTC().Add(5 * time.Minute)

	// Store bind code in user
	if err := l.svcCtx.UserModel.SetBindCode(userID, code, expire); err != nil {
		return nil, err
	}

	// Get bot username
	botUsername := "life_system_bot"
	if l.svcCtx.TelegramBot != nil {
		botUsername = l.svcCtx.TelegramBot.GetUsername()
	}

	return &types.BindCodeResp{
		Code:        code,
		BotUsername: botUsername,
		ExpiresIn:   300, // 5 minutes
	}, nil
}

func (l *TelegramLogic) GetStatus(ctx context.Context, userID int64) (*types.TgStatusResp, error) {
	user, err := l.svcCtx.UserModel.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	bound := user.TgChatID > 0
	return &types.TgStatusResp{
		Bound:      bound,
		TgUsername: user.TgUsername,
		TgChatID:   user.TgChatID,
	}, nil
}

func (l *TelegramLogic) Unbind(ctx context.Context, userID int64) error {
	user, err := l.svcCtx.UserModel.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	// Clear telegram binding
	return l.svcCtx.UserModel.UpdateTgBinding(userID, 0, "")
}

func generateRandomCode(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
