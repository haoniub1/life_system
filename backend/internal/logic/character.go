package logic

import (
	"context"
	"fmt"
	"math"
	"time"

	"life-system-backend/internal/model"
	"life-system-backend/internal/svc"
	"life-system-backend/internal/types"
)

type CharacterLogic struct {
	svcCtx *svc.ServiceContext
}

func NewCharacterLogic(svcCtx *svc.ServiceContext) *CharacterLogic {
	return &CharacterLogic{
		svcCtx: svcCtx,
	}
}

func (l *CharacterLogic) GetCharacter(ctx context.Context, userID int64) (*types.CharacterResp, error) {
	stats, err := l.svcCtx.CharacterModel.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if stats == nil {
		return nil, fmt.Errorf("character not found")
	}

	// Lazily trigger daily energy reset
	if l.CheckAndResetDailyEnergy(stats) {
		if err := l.svcCtx.CharacterModel.Update(stats); err != nil {
			return nil, err
		}
	}

	return l.statsToResp(stats), nil
}

func (l *CharacterLogic) UpdateCharacter(ctx context.Context, userID int64, req *types.UpdateCharacterReq) (*types.CharacterResp, error) {
	stats, err := l.svcCtx.CharacterModel.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if stats == nil {
		return nil, fmt.Errorf("character not found")
	}

	// Update fields if provided
	if req.Strength != nil {
		stats.Strength = *req.Strength
	}
	if req.Intelligence != nil {
		stats.Intelligence = *req.Intelligence
	}
	if req.Vitality != nil {
		stats.Vitality = *req.Vitality
	}
	if req.Spirit != nil {
		stats.Spirit = *req.Spirit
	}

	// Recalculate derived stats
	stats.MaxHP = 100 + int(stats.Strength*2) + int(stats.Vitality*3)
	if stats.HP > stats.MaxHP {
		stats.HP = stats.MaxHP
	}

	// Update in database
	if err := l.svcCtx.CharacterModel.Update(stats); err != nil {
		return nil, err
	}

	return l.statsToResp(stats), nil
}

func (l *CharacterLogic) statsToResp(stats *model.CharacterStats) *types.CharacterResp {
	return &types.CharacterResp{
		UserID:           stats.UserID,
		Level:            stats.Level,
		Exp:              stats.Exp,
		Strength:         stats.Strength,
		Intelligence:     stats.Intelligence,
		Vitality:         stats.Vitality,
		Spirit:           stats.Spirit,
		HP:               stats.HP,
		MaxHP:            stats.MaxHP,
		Gold:             stats.Gold,
		Title:            stats.Title,
		LastActivityDate: stats.LastActivityDate,
		Energy:           stats.Energy,
		MaxEnergy:        stats.MaxEnergy,
		MentalPower:      stats.MentalPower,
		PhysicalPower:    stats.PhysicalPower,
		MentalSleepAid:   stats.MentalSleepAid,
		PhysicalSleepAid: stats.PhysicalSleepAid,
	}
}

// CheckAndResetDailyEnergy resets energy based on sleep aid values when a new day starts.
// Returns true if a reset was performed (caller should persist).
func (l *CharacterLogic) CheckAndResetDailyEnergy(stats *model.CharacterStats) bool {
	today := time.Now().Format("2006-01-02")
	if stats.LastEnergyReset == today {
		return false
	}

	// Calculate sleep quality from accumulated sleep aid
	sleepQuality := 60.0 + float64(stats.MentalSleepAid)*0.2 + float64(stats.PhysicalSleepAid)*0.2
	sleepQuality = math.Min(100, sleepQuality)

	sq := int(sleepQuality)
	stats.MentalPower = sq
	stats.PhysicalPower = sq

	// Also sync legacy energy field
	stats.Energy = sq
	stats.MaxEnergy = 100

	// Reset sleep aid accumulators
	stats.MentalSleepAid = 0
	stats.PhysicalSleepAid = 0

	stats.LastEnergyReset = today

	fmt.Printf("🔄 Daily energy reset for user %d: sleepQuality=%d, mentalPower=%d, physicalPower=%d\n",
		stats.UserID, sq, stats.MentalPower, stats.PhysicalPower)

	return true
}

// Determine title based on level
func DetermineTitleByLevel(level int) string {
	switch {
	case level >= 50:
		return "生命大师🏆"
	case level >= 40:
		return "自律宗师👑"
	case level >= 30:
		return "优化专家⭐"
	case level >= 20:
		return "进化者🚀"
	case level >= 15:
		return "修行者🧘"
	case level >= 10:
		return "探索者🔍"
	case level >= 5:
		return "学徒📚"
	default:
		return "新手🌱"
	}
}

// Calculate experience needed for a level
func ExpForLevel(level int) int {
	// ExpForLevel(n) = 100 * 1.5^(n-1)
	exp := 100.0
	for i := 1; i < level; i++ {
		exp *= 1.5
	}
	return int(exp)
}

// Check if character should level up and do it
func CheckAndApplyLevelUp(stats *model.CharacterStats) bool {
	nextLevelExp := ExpForLevel(stats.Level + 1)
	if stats.Exp >= nextLevelExp {
		stats.Level++
		stats.Title = DetermineTitleByLevel(stats.Level)
		// Don't reset exp, it carries over
		return true
	}
	return false
}
