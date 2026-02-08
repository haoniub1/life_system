package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Database DatabaseConfig
	Auth     AuthConfig
	Telegram TelegramConfig
}

type DatabaseConfig struct {
	Path string
}

type AuthConfig struct {
	Secret string
	Expire int64
}

type TelegramConfig struct {
	BotToken string
	Enabled  bool
}
