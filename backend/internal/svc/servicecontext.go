package svc

import (
	"database/sql"
	"life-system-backend/internal/config"
	"life-system-backend/internal/model"
	"life-system-backend/pkg/bark"
	"life-system-backend/pkg/ratelimit"
	"life-system-backend/pkg/telegram"
)

type ServiceContext struct {
	Config         config.Config
	DB             *sql.DB
	UserModel      *model.UserModel
	CharacterModel *model.CharacterModel
	TaskModel      *model.TaskModel
	SleepModel     *model.SleepModel
	ShopModel      *model.ShopModel
	TelegramBot    *telegram.Bot
	BarkClient     *bark.Client
	RateLimiter    *ratelimit.Limiter
}

func NewServiceContext(cfg config.Config, db *sql.DB, bot *telegram.Bot) *ServiceContext {
	// Bark client always initialized with official server
	barkClient := bark.NewClient("https://api.day.app")

	// Rate limiter with configurable limits (defaults applied by go-zero)
	rateLimiter := ratelimit.NewLimiter(cfg.RateLimit.MaxLoginFailures, cfg.RateLimit.MaxDailyRegisters)

	ctx := &ServiceContext{
		Config:         cfg,
		DB:             db,
		UserModel:      model.NewUserModel(db),
		CharacterModel: model.NewCharacterModel(db),
		TaskModel:      model.NewTaskModel(db),
		SleepModel:     model.NewSleepModel(db),
		ShopModel:      model.NewShopModel(db),
		TelegramBot:    bot,
		BarkClient:     barkClient,
		RateLimiter:    rateLimiter,
	}

	// Set the service context reference in the bot to avoid circular import
	if bot != nil {
		bot.SetServiceContext(ctx)
	}

	return ctx
}

// Implement ServiceContextInterface for telegram package
func (s *ServiceContext) GetDB() *sql.DB {
	return s.DB
}

func (s *ServiceContext) GetUserModel() *model.UserModel {
	return s.UserModel
}

func (s *ServiceContext) GetTaskModel() *model.TaskModel {
	return s.TaskModel
}

func (s *ServiceContext) GetCharacterModel() *model.CharacterModel {
	return s.CharacterModel
}
