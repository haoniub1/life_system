package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"life-system-backend/internal/config"
	"life-system-backend/internal/handler"
	"life-system-backend/internal/logic"
	"life-system-backend/internal/middleware"
	"life-system-backend/internal/model"
	"life-system-backend/internal/svc"
	"life-system-backend/pkg/scheduler"
	"life-system-backend/pkg/telegram"
)

var configFile = flag.String("f", "etc/config.yaml", "the config file")

func main() {
	flag.Parse()

	// Load config
	var cfg config.Config
	conf.MustLoad(*configFile, &cfg)

	log.Printf("Starting Life System Backend on %s:%d", cfg.Host, cfg.Port)
	log.Printf("📋 Config loaded - Database.Path: '%s'", cfg.Database.Path)
	log.Printf("📋 Config loaded - Auth.Secret length: %d", len(cfg.Auth.Secret))
	log.Printf("📋 Config loaded - Telegram.Enabled: %v", cfg.Telegram.Enabled)

	// Initialize database
	db, err := model.NewDB(cfg.Database.Path)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := model.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Database initialized successfully")

	// Initialize Telegram bot if enabled
	var bot *telegram.Bot
	if cfg.Telegram.Enabled {
		var err error
		bot, err = telegram.NewBot(cfg.Telegram.BotToken, db)
		if err != nil {
			log.Fatalf("Failed to initialize Telegram bot: %v", err)
		}
		bot.Start()
		log.Println("Telegram bot started")
	}

	// Create service context
	svcCtx := svc.NewServiceContext(cfg, db, bot)

	// Set telegram task completer to avoid circular dependency
	if bot != nil {
		taskCompleter := logic.NewTelegramTaskCompleter(svcCtx)
		bot.SetTaskCompleter(taskCompleter)
	}

	// Initialize scheduler (always runs for daily reset and challenge task expiry)
	sched := scheduler.NewScheduler(bot, svcCtx, 0)
	sched.Start()
	log.Println("Task scheduler started (reminders, daily reset, challenge expiry)")

	// Create REST server
	server := rest.MustNewServer(cfg.RestConf)
	defer server.Stop()

	// Register CORS middleware globally
	server.Use(middleware.CORSMiddleware())

	// Register routes
	handler.RegisterRoutes(server, svcCtx)

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down gracefully...")

		// Stop scheduler
		sched.Stop()
		log.Println("Task scheduler stopped")

		// Stop bot
		if bot != nil {
			bot.Stop()
			log.Println("Telegram bot stopped")
		}

		// Stop server
		server.Stop()
		log.Println("Server stopped")

		os.Exit(0)
	}()

	// Start server
	log.Printf("REST server listening on %s:%d", cfg.Host, cfg.Port)
	server.Start()
}
