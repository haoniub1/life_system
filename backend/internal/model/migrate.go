package model

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite" // SQLite driver
)

func NewDB(path string) (*sql.DB, error) {
	fmt.Printf("📂 Opening database at: %s\n", path)

	// Ensure parent directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	fmt.Println("🔌 Testing database connection...")
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("✅ Database connection successful")
	return db, nil
}

func Migrate(db *sql.DB) error {
	fmt.Println("🔧 Starting database migration...")

	statements := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			display_name TEXT DEFAULT '',
			avatar TEXT DEFAULT '',
			tg_chat_id INTEGER DEFAULT 0,
			tg_username TEXT DEFAULT '',
			tg_bind_code TEXT DEFAULT '',
			tg_bind_expire DATETIME,
			bark_key TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS character_stats (
			user_id INTEGER PRIMARY KEY,
			spirit_stones INTEGER DEFAULT 0,
			fatigue INTEGER DEFAULT 0,
			fatigue_cap INTEGER DEFAULT 100,
			fatigue_level INTEGER DEFAULT 0,
			overdraft_penalty REAL DEFAULT 0,
			title TEXT DEFAULT '凡人',
			last_activity_date TEXT DEFAULT (date('now')),
			last_fatigue_reset TEXT DEFAULT '',
			FOREIGN KEY(user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS character_attributes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			attr_key TEXT NOT NULL,
			value REAL DEFAULT 100,
			realm INTEGER DEFAULT 0,
			sub_realm INTEGER DEFAULT 0,
			realm_exp INTEGER DEFAULT 0,
			is_bottleneck INTEGER DEFAULT 0,
			accumulation_pool REAL DEFAULT 0,
			FOREIGN KEY(user_id) REFERENCES users(id),
			UNIQUE(user_id, attr_key)
		)`,
		`CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			description TEXT DEFAULT '',
			category TEXT DEFAULT 'general',
			type TEXT NOT NULL DEFAULT 'once',
			status TEXT NOT NULL DEFAULT 'active',
			deadline DATETIME,
			primary_attribute TEXT DEFAULT '',
			difficulty INTEGER DEFAULT 1,
			reward_exp INTEGER DEFAULT 0,
			reward_spirit_stones INTEGER DEFAULT 0,
			reward_physique REAL DEFAULT 0,
			reward_willpower REAL DEFAULT 0,
			reward_intelligence REAL DEFAULT 0,
			reward_perception REAL DEFAULT 0,
			reward_charisma REAL DEFAULT 0,
			reward_agility REAL DEFAULT 0,
			fatigue_cost INTEGER DEFAULT 10,
			penalty_exp INTEGER DEFAULT 0,
			penalty_spirit_stones INTEGER DEFAULT 0,
			daily_limit INTEGER DEFAULT 0,
			total_limit INTEGER DEFAULT 0,
			completed_count INTEGER DEFAULT 0,
			today_completion_count INTEGER DEFAULT 0,
			last_completed_date TEXT DEFAULT '',
			remind_before INTEGER DEFAULT 0,
			remind_interval INTEGER DEFAULT 0,
			last_reminded_at DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS task_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			task_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			action TEXT NOT NULL,
			source TEXT DEFAULT 'web',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS sleep_records (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			sleep_start DATETIME NOT NULL,
			sleep_end DATETIME NOT NULL,
			duration_hours REAL NOT NULL,
			quality TEXT DEFAULT 'good',
			energy_gained INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS shop_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			description TEXT DEFAULT '',
			price INTEGER NOT NULL,
			item_type TEXT NOT NULL,
			effect TEXT NOT NULL,
			effect_value INTEGER DEFAULT 0,
			icon TEXT DEFAULT '',
			image TEXT DEFAULT '',
			stock INTEGER DEFAULT -1,
			sell_price INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS inventory (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			item_id INTEGER NOT NULL,
			quantity INTEGER DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(user_id) REFERENCES users(id),
			FOREIGN KEY(item_id) REFERENCES shop_items(id),
			UNIQUE(user_id, item_id)
		)`,
		`CREATE TABLE IF NOT EXISTS purchase_history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			item_id INTEGER NOT NULL,
			item_name TEXT NOT NULL,
			quantity INTEGER NOT NULL,
			total_price INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(user_id) REFERENCES users(id)
		)`,
	}

	// Add sell_price column to shop_items if not exists (migration)
	migrations := []string{
		`ALTER TABLE shop_items ADD COLUMN sell_price INTEGER DEFAULT 0`,
		`ALTER TABLE tasks ADD COLUMN sort_order INTEGER DEFAULT 0`,
		`ALTER TABLE character_attributes ADD COLUMN today_gain REAL DEFAULT 0`,
		`ALTER TABLE character_attributes ADD COLUMN last_gain_date TEXT DEFAULT ''`,
	}
	for _, m := range migrations {
		db.Exec(m) // Ignore errors (column may already exist)
	}

	tableNames := []string{"users", "character_stats", "character_attributes", "tasks", "task_logs", "sleep_records", "shop_items", "inventory", "purchase_history"}

	for i, stmt := range statements {
		fmt.Printf("  Creating table '%s'...\n", tableNames[i])
		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("migration error on table '%s': %w", tableNames[i], err)
		}
		fmt.Printf("  ✅ Table '%s' created\n", tableNames[i])
	}

	// Verify tables were created
	fmt.Println("🔍 Verifying tables...")
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' ORDER BY name")
	if err != nil {
		return fmt.Errorf("failed to verify tables: %w", err)
	}
	defer rows.Close()

	tableCount := 0
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("failed to scan table name: %w", err)
		}
		fmt.Printf("  Found table: %s\n", name)
		tableCount++
	}

	if tableCount == 0 {
		return fmt.Errorf("⚠️  WARNING: No tables found in database after migration!")
	}

	fmt.Printf("✅ Migration completed successfully! Created %d tables\n", tableCount)
	return nil
}
