package model

import (
	"database/sql"
)

type CharacterStats struct {
	UserID           int64
	SpiritStones     int
	Fatigue          int
	FatigueCap       int
	FatigueLevel     int
	OverdraftPenalty float64
	Title            string
	LastActivityDate string
	LastFatigueReset string
}

type CharacterAttribute struct {
	ID               int64
	UserID           int64
	AttrKey          string
	Value            float64
	Realm            int
	SubRealm         int
	RealmExp         int
	IsBottleneck     bool
	AccumulationPool float64
}

type CharacterModel struct {
	db *sql.DB
}

func NewCharacterModel(db *sql.DB) *CharacterModel {
	return &CharacterModel{db: db}
}

func (m *CharacterModel) FindByUserID(userID int64) (*CharacterStats, error) {
	var stats CharacterStats
	err := m.db.QueryRow(`
		SELECT user_id, spirit_stones, fatigue, fatigue_cap, fatigue_level,
		       overdraft_penalty, title, last_activity_date, last_fatigue_reset
		FROM character_stats WHERE user_id = ?
	`, userID).Scan(
		&stats.UserID, &stats.SpiritStones, &stats.Fatigue, &stats.FatigueCap,
		&stats.FatigueLevel, &stats.OverdraftPenalty, &stats.Title,
		&stats.LastActivityDate, &stats.LastFatigueReset,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &stats, nil
}

func (m *CharacterModel) Create(userID int64) error {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert character_stats row
	_, err = tx.Exec(`
		INSERT INTO character_stats (user_id, spirit_stones, fatigue, fatigue_cap, fatigue_level,
		                             overdraft_penalty, title, last_activity_date, last_fatigue_reset)
		VALUES (?, 0, 0, 100, 0, 0, '凡人', date('now'), '')
	`, userID)
	if err != nil {
		return err
	}

	// Insert 7 attribute rows (6 cultivation + luck), all starting at value=100, realm=0
	attrKeys := []string{"physique", "willpower", "intelligence", "perception", "charisma", "agility", "luck"}
	for _, key := range attrKeys {
		_, err = tx.Exec(`
			INSERT INTO character_attributes (user_id, attr_key, value, realm, sub_realm, realm_exp, is_bottleneck, accumulation_pool)
			VALUES (?, ?, 100, 0, 0, 0, 0, 0)
		`, userID, key)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (m *CharacterModel) Update(stats *CharacterStats) error {
	_, err := m.db.Exec(`
		UPDATE character_stats
		SET spirit_stones = ?, fatigue = ?, fatigue_cap = ?, fatigue_level = ?,
		    overdraft_penalty = ?, title = ?, last_activity_date = ?, last_fatigue_reset = ?
		WHERE user_id = ?
	`, stats.SpiritStones, stats.Fatigue, stats.FatigueCap, stats.FatigueLevel,
		stats.OverdraftPenalty, stats.Title, stats.LastActivityDate, stats.LastFatigueReset,
		stats.UserID)

	return err
}

func (m *CharacterModel) FindAttributesByUserID(userID int64) ([]*CharacterAttribute, error) {
	rows, err := m.db.Query(`
		SELECT id, user_id, attr_key, value, realm, sub_realm, realm_exp, is_bottleneck, accumulation_pool
		FROM character_attributes WHERE user_id = ?
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attrs []*CharacterAttribute
	for rows.Next() {
		var attr CharacterAttribute
		err := rows.Scan(
			&attr.ID, &attr.UserID, &attr.AttrKey, &attr.Value,
			&attr.Realm, &attr.SubRealm, &attr.RealmExp,
			&attr.IsBottleneck, &attr.AccumulationPool,
		)
		if err != nil {
			return nil, err
		}
		attrs = append(attrs, &attr)
	}

	return attrs, rows.Err()
}

func (m *CharacterModel) UpdateAttribute(attr *CharacterAttribute) error {
	_, err := m.db.Exec(`
		UPDATE character_attributes
		SET value = ?, realm = ?, sub_realm = ?, realm_exp = ?,
		    is_bottleneck = ?, accumulation_pool = ?
		WHERE user_id = ? AND attr_key = ?
	`, attr.Value, attr.Realm, attr.SubRealm, attr.RealmExp,
		attr.IsBottleneck, attr.AccumulationPool,
		attr.UserID, attr.AttrKey)

	return err
}

// FindInactiveCharacters finds all characters that haven't been active for the specified number of days
func (m *CharacterModel) FindInactiveCharacters(daysThreshold int) ([]*CharacterStats, error) {
	rows, err := m.db.Query(`
		SELECT user_id, spirit_stones, fatigue, fatigue_cap, fatigue_level,
		       overdraft_penalty, title, last_activity_date, last_fatigue_reset
		FROM character_stats
		WHERE last_activity_date < date('now', '-' || ? || ' days')
	`, daysThreshold)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var characters []*CharacterStats
	for rows.Next() {
		var stats CharacterStats
		err := rows.Scan(
			&stats.UserID, &stats.SpiritStones, &stats.Fatigue, &stats.FatigueCap,
			&stats.FatigueLevel, &stats.OverdraftPenalty, &stats.Title,
			&stats.LastActivityDate, &stats.LastFatigueReset,
		)
		if err != nil {
			return nil, err
		}
		characters = append(characters, &stats)
	}

	return characters, rows.Err()
}
