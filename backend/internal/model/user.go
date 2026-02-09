package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID           int64
	Username     string
	PasswordHash string
	DisplayName  string
	Avatar       string
	TgChatID     int64
	TgUsername   string
	TgBindCode   string
	TgBindExpire sql.NullTime
	BarkKey      string // Bark push notification device key
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserModel struct {
	db *sql.DB
}

func NewUserModel(db *sql.DB) *UserModel {
	return &UserModel{db: db}
}

func (m *UserModel) FindByUsername(username string) (*User, error) {
	var user User
	err := m.db.QueryRow(`
		SELECT id, username, password_hash, display_name, avatar, tg_chat_id, tg_username, tg_bind_code, tg_bind_expire, bark_key, created_at, updated_at
		FROM users WHERE username = ?
	`, username).Scan(
		&user.ID, &user.Username, &user.PasswordHash, &user.DisplayName, &user.Avatar,
		&user.TgChatID, &user.TgUsername, &user.TgBindCode, &user.TgBindExpire,
		&user.BarkKey, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (m *UserModel) FindByID(id int64) (*User, error) {
	var user User
	err := m.db.QueryRow(`
		SELECT id, username, password_hash, display_name, avatar, tg_chat_id, tg_username, tg_bind_code, tg_bind_expire, bark_key, created_at, updated_at
		FROM users WHERE id = ?
	`, id).Scan(
		&user.ID, &user.Username, &user.PasswordHash, &user.DisplayName, &user.Avatar,
		&user.TgChatID, &user.TgUsername, &user.TgBindCode, &user.TgBindExpire,
		&user.BarkKey, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (m *UserModel) FindByTgChatID(chatID int64) (*User, error) {
	var user User
	err := m.db.QueryRow(`
		SELECT id, username, password_hash, display_name, avatar, tg_chat_id, tg_username, tg_bind_code, tg_bind_expire, bark_key, created_at, updated_at
		FROM users WHERE tg_chat_id = ?
	`, chatID).Scan(
		&user.ID, &user.Username, &user.PasswordHash, &user.DisplayName, &user.Avatar,
		&user.TgChatID, &user.TgUsername, &user.TgBindCode, &user.TgBindExpire,
		&user.BarkKey, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (m *UserModel) Create(username, passwordHash string) (int64, error) {
	result, err := m.db.Exec(`
		INSERT INTO users (username, password_hash, display_name, created_at, updated_at)
		VALUES (?, ?, ?, datetime('now'), datetime('now'))
	`, username, passwordHash, username)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (m *UserModel) UpdateTgBinding(userID int64, chatID int64, username string) error {
	_, err := m.db.Exec(`
		UPDATE users SET tg_chat_id = ?, tg_username = ?, updated_at = datetime('now')
		WHERE id = ?
	`, chatID, username, userID)

	return err
}

func (m *UserModel) SetBindCode(userID int64, code string, expire time.Time) error {
	// Format expire time as UTC string to match SQLite's datetime('now') format
	expireStr := expire.UTC().Format("2006-01-02 15:04:05")
	_, err := m.db.Exec(`
		UPDATE users SET tg_bind_code = ?, tg_bind_expire = ?, updated_at = datetime('now')
		WHERE id = ?
	`, code, expireStr, userID)

	return err
}

func (m *UserModel) FindByBindCode(code string) (*User, error) {
	var user User
	err := m.db.QueryRow(`
		SELECT id, username, password_hash, display_name, avatar, tg_chat_id, tg_username, tg_bind_code, tg_bind_expire, bark_key, created_at, updated_at
		FROM users WHERE tg_bind_code = ? AND tg_bind_expire > datetime('now')
	`, code).Scan(
		&user.ID, &user.Username, &user.PasswordHash, &user.DisplayName, &user.Avatar,
		&user.TgChatID, &user.TgUsername, &user.TgBindCode, &user.TgBindExpire,
		&user.BarkKey, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (m *UserModel) ClearBindCode(userID int64) error {
	_, err := m.db.Exec(`
		UPDATE users SET tg_bind_code = '', tg_bind_expire = NULL, updated_at = datetime('now')
		WHERE id = ?
	`, userID)

	return err
}

func (m *UserModel) UpdateProfile(userID int64, displayName, avatar string) error {
	_, err := m.db.Exec(`
		UPDATE users SET display_name = ?, avatar = ?, updated_at = datetime('now')
		WHERE id = ?
	`, displayName, avatar, userID)

	return err
}

func (m *UserModel) UpdatePassword(userID int64, newPasswordHash string) error {
	_, err := m.db.Exec(`
		UPDATE users SET password_hash = ?, updated_at = datetime('now')
		WHERE id = ?
	`, newPasswordHash, userID)

	return err
}

func (m *UserModel) UpdateBarkKey(userID int64, barkKey string) error {
	_, err := m.db.Exec(`
		UPDATE users SET bark_key = ?, updated_at = datetime('now')
		WHERE id = ?
	`, barkKey, userID)

	return err
}

func (m *UserModel) GetBarkKey(userID int64) (string, error) {
	var barkKey string
	err := m.db.QueryRow(`SELECT bark_key FROM users WHERE id = ?`, userID).Scan(&barkKey)
	if err != nil {
		return "", err
	}
	return barkKey, nil
}
