package model

import (
	"database/sql"
)

type CharacterStats struct {
	UserID           int64
	Level            int
	Exp              int
	Strength         float64
	Intelligence     float64
	Vitality         float64
	Spirit           float64
	HP               int
	MaxHP            int
	Gold             int
	Title            string
	LastActivityDate string // YYYY-MM-DD format for tracking daily activity
	Energy           int    // Current energy (kept for compatibility)
	MaxEnergy        int    // Maximum energy (kept for compatibility)
	MentalPower      int    // Current mental power
	PhysicalPower    int    // Current physical power
	MentalSleepAid   int    // Accumulated mental sleep aid
	PhysicalSleepAid int    // Accumulated physical sleep aid
	LastEnergyReset  string // Last daily energy reset date (YYYY-MM-DD)
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
		SELECT user_id, level, exp, strength, intelligence, vitality, spirit, hp, max_hp, gold, title, last_activity_date, energy, max_energy,
		       mental_power, physical_power, mental_sleep_aid, physical_sleep_aid, last_energy_reset
		FROM character_stats WHERE user_id = ?
	`, userID).Scan(
		&stats.UserID, &stats.Level, &stats.Exp, &stats.Strength, &stats.Intelligence,
		&stats.Vitality, &stats.Spirit, &stats.HP, &stats.MaxHP, &stats.Gold, &stats.Title,
		&stats.LastActivityDate, &stats.Energy, &stats.MaxEnergy,
		&stats.MentalPower, &stats.PhysicalPower, &stats.MentalSleepAid, &stats.PhysicalSleepAid, &stats.LastEnergyReset,
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
	_, err := m.db.Exec(`
		INSERT INTO character_stats (user_id, level, exp, strength, intelligence, vitality, spirit, hp, max_hp, gold, title, last_activity_date, energy, max_energy,
		                             mental_power, physical_power, mental_sleep_aid, physical_sleep_aid, last_energy_reset)
		VALUES (?, 1, 0, 5.0, 5.0, 5.0, 5.0, 100, 100, 0, '新手🌱', date('now'), 100, 100, 100, 100, 0, 0, '')
	`, userID)

	return err
}

func (m *CharacterModel) Update(stats *CharacterStats) error {
	_, err := m.db.Exec(`
		UPDATE character_stats
		SET level = ?, exp = ?, strength = ?, intelligence = ?, vitality = ?, spirit = ?, hp = ?, max_hp = ?, gold = ?, title = ?, last_activity_date = ?, energy = ?, max_energy = ?,
		    mental_power = ?, physical_power = ?, mental_sleep_aid = ?, physical_sleep_aid = ?, last_energy_reset = ?
		WHERE user_id = ?
	`, stats.Level, stats.Exp, stats.Strength, stats.Intelligence, stats.Vitality, stats.Spirit,
		stats.HP, stats.MaxHP, stats.Gold, stats.Title, stats.LastActivityDate, stats.Energy, stats.MaxEnergy,
		stats.MentalPower, stats.PhysicalPower, stats.MentalSleepAid, stats.PhysicalSleepAid, stats.LastEnergyReset, stats.UserID)

	return err
}

// FindInactiveCharacters finds all characters that haven't been active for the specified number of days
func (m *CharacterModel) FindInactiveCharacters(daysThreshold int) ([]*CharacterStats, error) {
	rows, err := m.db.Query(`
		SELECT user_id, level, exp, strength, intelligence, vitality, spirit, hp, max_hp, gold, title, last_activity_date, energy, max_energy,
		       mental_power, physical_power, mental_sleep_aid, physical_sleep_aid, last_energy_reset
		FROM character_stats
		WHERE last_activity_date < date('now', '-' || ? || ' days')
		  AND (strength > 5.0 OR intelligence > 5.0 OR vitality > 5.0 OR spirit > 5.0)
	`, daysThreshold)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var characters []*CharacterStats
	for rows.Next() {
		var stats CharacterStats
		err := rows.Scan(
			&stats.UserID, &stats.Level, &stats.Exp, &stats.Strength, &stats.Intelligence,
			&stats.Vitality, &stats.Spirit, &stats.HP, &stats.MaxHP, &stats.Gold, &stats.Title,
			&stats.LastActivityDate, &stats.Energy, &stats.MaxEnergy,
			&stats.MentalPower, &stats.PhysicalPower, &stats.MentalSleepAid, &stats.PhysicalSleepAid, &stats.LastEnergyReset,
		)
		if err != nil {
			return nil, err
		}
		characters = append(characters, &stats)
	}

	return characters, rows.Err()
}
