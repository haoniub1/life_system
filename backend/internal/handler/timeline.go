package handler

import (
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/zeromicro/go-zero/rest/httpx"
	"life-system-backend/internal/middleware"
	"life-system-backend/internal/svc"
	"life-system-backend/internal/types"
)

func GetTimelineHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetUserID(r.Context())
		if err != nil {
			httpx.OkJson(w, types.CommonResp{Code: 401, Message: "unauthorized"})
			return
		}

		var events []types.TimelineEvent

		// 1. Fetch task logs with task details
		taskLogs, err := svcCtx.DB.Query(`
			SELECT tl.id, tl.action, tl.source, tl.created_at, t.title, t.reward_exp, t.reward_spirit_stones
			FROM task_logs tl
			LEFT JOIN tasks t ON tl.task_id = t.id
			WHERE tl.user_id = ?
			ORDER BY tl.created_at DESC
			LIMIT 100
		`, userID)
		if err == nil {
			defer taskLogs.Close()
			for taskLogs.Next() {
				var id int64
				var action, source, createdAt, title string
				var rewardExp, rewardSpiritStones int
				if err := taskLogs.Scan(&id, &action, &source, &createdAt, &title, &rewardExp, &rewardSpiritStones); err != nil {
					continue
				}

				eventType := "task_complete"
				eventTitle := fmt.Sprintf("完成任务：%s", title)
				desc := fmt.Sprintf("通过 %s 完成", source)

				switch action {
				case "fail":
					eventType = "task_fail"
					eventTitle = fmt.Sprintf("任务失败：%s", title)
					desc = "任务超时未完成"
				case "delete":
					eventType = "task_delete"
					eventTitle = fmt.Sprintf("删除任务：%s", title)
					desc = fmt.Sprintf("通过 %s 删除", source)
				}

				event := types.TimelineEvent{
					ID:          fmt.Sprintf("task_%d", id),
					Type:        eventType,
					Title:       eventTitle,
					Description: desc,
					Timestamp:   createdAt,
				}

				if action == "complete" && (rewardExp > 0 || rewardSpiritStones > 0) {
					event.Rewards = &types.TimelineRewards{
						Exp:          rewardExp,
						SpiritStones: rewardSpiritStones,
					}
				}

				events = append(events, event)
			}
		}

		// 2. Fetch sleep records
		sleepRows, err := svcCtx.DB.Query(`
			SELECT id, duration_hours, quality, energy_gained, created_at
			FROM sleep_records
			WHERE user_id = ?
			ORDER BY created_at DESC
			LIMIT 50
		`, userID)
		if err == nil {
			defer sleepRows.Close()
			for sleepRows.Next() {
				var id int64
				var durationHours float64
				var quality, createdAt string
				var energyGained int
				if err := sleepRows.Scan(&id, &durationHours, &quality, &energyGained, &createdAt); err != nil {
					continue
				}

				qualityMap := map[string]string{
					"poor": "较差", "fair": "一般", "good": "良好", "excellent": "优秀",
				}
				qualityStr := qualityMap[quality]
				if qualityStr == "" {
					qualityStr = quality
				}

				event := types.TimelineEvent{
					ID:          fmt.Sprintf("sleep_%d", id),
					Type:        "sleep",
					Title:       "记录睡眠",
					Description: fmt.Sprintf("睡眠时长: %.1f 小时 | 质量: %s", durationHours, qualityStr),
					Timestamp:   createdAt,
				}

				events = append(events, event)
			}
		}

		// 3. Fetch purchase history
		purchaseRows, err := svcCtx.DB.Query(`
			SELECT id, item_name, quantity, total_price, created_at
			FROM purchase_history
			WHERE user_id = ?
			ORDER BY created_at DESC
			LIMIT 50
		`, userID)
		if err == nil {
			defer purchaseRows.Close()
			for purchaseRows.Next() {
				var id int64
				var itemName, createdAt string
				var quantity, totalPrice int
				if err := purchaseRows.Scan(&id, &itemName, &quantity, &totalPrice, &createdAt); err != nil {
					continue
				}

				events = append(events, types.TimelineEvent{
					ID:          fmt.Sprintf("purchase_%d", id),
					Type:        "purchase",
					Title:       fmt.Sprintf("购买：%s", itemName),
					Description: fmt.Sprintf("购买了 %d 个，花费 %d 灵石", quantity, totalPrice),
					Timestamp:   createdAt,
					Rewards:     &types.TimelineRewards{SpiritStones: -totalPrice},
				})
			}
		}

		// Sort all events by timestamp descending
		parseTime := func(s string) time.Time {
			formats := []string{
				"2006-01-02 15:04:05",
				"2006-01-02T15:04:05Z",
				time.RFC3339,
			}
			for _, f := range formats {
				if t, err := time.Parse(f, s); err == nil {
					return t
				}
			}
			return time.Time{}
		}
		sort.Slice(events, func(i, j int) bool {
			return parseTime(events[i].Timestamp).After(parseTime(events[j].Timestamp))
		})

		// Limit to 100 events total
		if len(events) > 100 {
			events = events[:100]
		}

		// Get stats
		var tasksCompleted int
		svcCtx.DB.QueryRow(`SELECT COUNT(*) FROM task_logs WHERE user_id = ? AND action = 'complete'`, userID).Scan(&tasksCompleted)

		var sleepCount int
		svcCtx.DB.QueryRow(`SELECT COUNT(*) FROM sleep_records WHERE user_id = ?`, userID).Scan(&sleepCount)

		// Get character stats for totals
		character, _ := svcCtx.CharacterModel.FindByUserID(userID)
		totalSpiritStones := 0
		if character != nil {
			totalSpiritStones = character.SpiritStones
		}

		httpx.OkJson(w, types.CommonResp{
			Code:    0,
			Message: "success",
			Data: types.TimelineResp{
				Events:            events,
				TasksCompleted:    tasksCompleted,
				TotalExp:          0,
				TotalSpiritStones: totalSpiritStones,
				SleepRecords:      sleepCount,
			},
		})
	}
}
