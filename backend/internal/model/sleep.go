package model

import (
	"database/sql"
	"time"
)

type SleepRecord struct {
	ID            int64
	UserID        int64
	SleepStart    time.Time
	SleepEnd      time.Time
	DurationHours float64
	Quality       string // poor, fair, good, excellent
	EnergyGained  int    // Energy restored from this sleep
	CreatedAt     time.Time
}

type SleepModel struct {
	db *sql.DB
}

func NewSleepModel(db *sql.DB) *SleepModel {
	return &SleepModel{db: db}
}

func (m *SleepModel) Create(record *SleepRecord) (int64, error) {
	result, err := m.db.Exec(`
		INSERT INTO sleep_records (user_id, sleep_start, sleep_end, duration_hours, quality, energy_gained, created_at)
		VALUES (?, ?, ?, ?, ?, ?, datetime('now'))
	`, record.UserID, record.SleepStart, record.SleepEnd, record.DurationHours, record.Quality, record.EnergyGained)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (m *SleepModel) FindByUserID(userID int64, limit int) ([]*SleepRecord, error) {
	query := `
		SELECT id, user_id, sleep_start, sleep_end, duration_hours, quality, energy_gained, created_at
		FROM sleep_records
		WHERE user_id = ?
		ORDER BY sleep_start DESC
	`

	if limit > 0 {
		query += ` LIMIT ?`
	}

	var rows *sql.Rows
	var err error

	if limit > 0 {
		rows, err = m.db.Query(query, userID, limit)
	} else {
		rows, err = m.db.Query(query, userID)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*SleepRecord
	for rows.Next() {
		var record SleepRecord
		err := rows.Scan(
			&record.ID, &record.UserID, &record.SleepStart, &record.SleepEnd,
			&record.DurationHours, &record.Quality, &record.EnergyGained, &record.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, &record)
	}

	return records, rows.Err()
}

func (m *SleepModel) FindByID(id int64) (*SleepRecord, error) {
	var record SleepRecord
	err := m.db.QueryRow(`
		SELECT id, user_id, sleep_start, sleep_end, duration_hours, quality, energy_gained, created_at
		FROM sleep_records
		WHERE id = ?
	`, id).Scan(
		&record.ID, &record.UserID, &record.SleepStart, &record.SleepEnd,
		&record.DurationHours, &record.Quality, &record.EnergyGained, &record.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &record, nil
}

// CalculateEnergyGain calculates energy gain based on sleep duration and quality
func CalculateEnergyGain(durationHours float64, quality string, maxEnergy int) int {
	// Base energy from sleep duration
	var baseEnergy float64

	switch {
	case durationHours >= 8:
		baseEnergy = float64(maxEnergy) // Full restore for 8+ hours
	case durationHours >= 6:
		baseEnergy = float64(maxEnergy) * 0.8 // 80% for 6-8 hours
	case durationHours >= 4:
		baseEnergy = float64(maxEnergy) * 0.5 // 50% for 4-6 hours
	default:
		baseEnergy = float64(maxEnergy) * 0.3 // 30% for < 4 hours
	}

	// Quality multiplier
	var qualityMultiplier float64
	switch quality {
	case "excellent":
		qualityMultiplier = 1.2
	case "good":
		qualityMultiplier = 1.0
	case "fair":
		qualityMultiplier = 0.8
	case "poor":
		qualityMultiplier = 0.6
	default:
		qualityMultiplier = 1.0
	}

	energyGained := int(baseEnergy * qualityMultiplier)

	// Cap at max energy
	if energyGained > maxEnergy {
		energyGained = maxEnergy
	}

	return energyGained
}
