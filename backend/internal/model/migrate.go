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
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS character_stats (
			user_id INTEGER PRIMARY KEY,
			level INTEGER DEFAULT 1,
			exp INTEGER DEFAULT 0,
			strength REAL DEFAULT 5.0,
			intelligence REAL DEFAULT 5.0,
			vitality REAL DEFAULT 5.0,
			spirit REAL DEFAULT 5.0,
			hp INTEGER DEFAULT 100,
			max_hp INTEGER DEFAULT 100,
			gold INTEGER DEFAULT 0,
			title TEXT DEFAULT '新手🌱',
			last_activity_date TEXT DEFAULT (date('now')),
			energy INTEGER DEFAULT 100,
			max_energy INTEGER DEFAULT 100,
			FOREIGN KEY(user_id) REFERENCES users(id)
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
			reward_exp INTEGER DEFAULT 0,
			reward_gold INTEGER DEFAULT 0,
			reward_strength REAL DEFAULT 0,
			reward_intelligence REAL DEFAULT 0,
			reward_vitality REAL DEFAULT 0,
			reward_spirit REAL DEFAULT 0,
			penalty_exp INTEGER DEFAULT 0,
			penalty_gold INTEGER DEFAULT 0,
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

	tableNames := []string{"users", "character_stats", "tasks", "task_logs", "sleep_records", "shop_items", "inventory", "purchase_history"}

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

	// Apply incremental schema updates
	if err := applySchemaUpdates(db); err != nil {
		return fmt.Errorf("schema update error: %w", err)
	}

	return nil
}

func applySchemaUpdates(db *sql.DB) error {
	fmt.Println("🔧 Applying schema updates...")

	// Each update is idempotent — uses ADD COLUMN which fails silently if column exists
	updates := []struct {
		desc string
		sql  string
	}{
		{"character_stats.mental_power", "ALTER TABLE character_stats ADD COLUMN mental_power INTEGER DEFAULT 100"},
		{"character_stats.physical_power", "ALTER TABLE character_stats ADD COLUMN physical_power INTEGER DEFAULT 100"},
		{"character_stats.mental_sleep_aid", "ALTER TABLE character_stats ADD COLUMN mental_sleep_aid INTEGER DEFAULT 0"},
		{"character_stats.physical_sleep_aid", "ALTER TABLE character_stats ADD COLUMN physical_sleep_aid INTEGER DEFAULT 0"},
		{"character_stats.last_energy_reset", "ALTER TABLE character_stats ADD COLUMN last_energy_reset TEXT DEFAULT ''"},
		{"tasks.cost_mental", "ALTER TABLE tasks ADD COLUMN cost_mental INTEGER DEFAULT 0"},
		{"tasks.cost_physical", "ALTER TABLE tasks ADD COLUMN cost_physical INTEGER DEFAULT 0"},
		{"users.bark_key", "ALTER TABLE users ADD COLUMN bark_key TEXT DEFAULT ''"},
	}

	for _, u := range updates {
		if _, err := db.Exec(u.sql); err != nil {
			// SQLite returns "duplicate column name" if column already exists — skip
			fmt.Printf("  ⏭️  %s (already exists or skipped)\n", u.desc)
		} else {
			fmt.Printf("  ✅ Added %s\n", u.desc)
		}
	}

	fmt.Println("✅ Schema updates complete")
	return nil
}
