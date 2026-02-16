package logic

import (
	"context"
	"fmt"
	"time"

	"life-system-backend/internal/model"
	"life-system-backend/internal/realm"
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

	attrs, err := l.svcCtx.CharacterModel.FindAttributesByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Lazily trigger daily fatigue reset
	if l.CheckAndResetDailyFatigue(stats) {
		if err := l.svcCtx.CharacterModel.Update(stats); err != nil {
			return nil, err
		}
	}

	return l.statsToResp(stats, attrs), nil
}

// CheckAndResetDailyFatigue resets fatigue when a new day starts.
// Returns true if a reset was performed (caller should persist).
func (l *CharacterLogic) CheckAndResetDailyFatigue(stats *model.CharacterStats) bool {
	today := time.Now().Format("2006-01-02")
	if stats.LastFatigueReset == today {
		return false
	}

	// Reset fatigue to 0
	stats.Fatigue = 0

	// Reset fatigue cap (no penalty)
	stats.FatigueCap = realm.FatigueCapForLevel(stats.FatigueLevel)
	stats.OverdraftPenalty = 0

	stats.LastFatigueReset = today

	fmt.Printf("🔄 Daily fatigue reset for user %d: fatigueCap=%d\n",
		stats.UserID, stats.FatigueCap)

	return true
}

func (l *CharacterLogic) statsToResp(stats *model.CharacterStats, attrs []*model.CharacterAttribute) *types.CharacterResp {
	resp := &types.CharacterResp{
		UserID:           stats.UserID,
		SpiritStones:     stats.SpiritStones,
		Fatigue:          stats.Fatigue,
		FatigueCap:       stats.FatigueCap,
		FatigueLevel:     stats.FatigueLevel,
		OverdraftPenalty: stats.OverdraftPenalty,
		Title:            stats.Title,
		LastActivityDate: stats.LastActivityDate,
		Attributes:       make([]types.AttributeResp, 0, len(attrs)),
	}

	today := time.Now().Format("2006-01-02")
	
	for _, attr := range attrs {
		display, ok := realm.AttrDisplay[attr.AttrKey]
		if !ok {
			continue
		}

		// Reset today_gain if not today
		todayGain := attr.TodayGain
		if attr.LastGainDate != today {
			todayGain = 0
		}

		cap := realm.AttrCap(attr.Realm)
		minVal := realm.AttrMin(attr.Realm)
		rangeVal := cap - minVal
		var progressPercent float64
		if rangeVal > 0 {
			progressPercent = (attr.Value - minVal) / rangeVal * 100
			if progressPercent > 100 {
				progressPercent = 100
			}
			if progressPercent < 0 {
				progressPercent = 0
			}
		}

		resp.Attributes = append(resp.Attributes, types.AttributeResp{
			AttrKey:          attr.AttrKey,
			DisplayName:      display.Name,
			Emoji:            display.Emoji,
			Value:            attr.Value,
			TodayGain:        todayGain,
			Realm:            attr.Realm,
			RealmName:        realm.GetRealmName(attr.Realm),
			SubRealm:         attr.SubRealm,
			SubRealmName:     realm.GetSubRealmName(attr.SubRealm),
			RealmExp:         attr.RealmExp,
			IsBottleneck:     attr.IsBottleneck,
			AccumulationPool: attr.AccumulationPool,
			AttrCap:          cap,
			ProgressPercent:  progressPercent,
			Color:            display.Color,
		})
	}

	return resp
}
