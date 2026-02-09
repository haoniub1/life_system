package logic

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"life-system-backend/internal/model"
	"life-system-backend/internal/svc"
	"life-system-backend/internal/types"
)

type TaskLogic struct {
	svcCtx *svc.ServiceContext
}

func NewTaskLogic(svcCtx *svc.ServiceContext) *TaskLogic {
	return &TaskLogic{
		svcCtx: svcCtx,
	}
}

func (l *TaskLogic) ListTasks(ctx context.Context, userID int64, taskType, status string) (*types.TaskListResp, error) {
	tasks, err := l.svcCtx.TaskModel.FindByUserID(userID, taskType, status)
	if err != nil {
		return nil, err
	}

	resp := &types.TaskListResp{
		Tasks: make([]types.TaskResp, 0),
	}

	for _, task := range tasks {
		resp.Tasks = append(resp.Tasks, l.taskToResp(task))
	}

	return resp, nil
}

func (l *TaskLogic) CreateTask(ctx context.Context, userID int64, req *types.CreateTaskReq) (*types.TaskResp, error) {
	// Validate input
	if req.Title == "" {
		return nil, fmt.Errorf("title is required")
	}

	// Parse deadline if provided
	var deadline sql.NullTime
	if req.Deadline != "" {
		t, err := time.Parse(time.RFC3339, req.Deadline)
		if err != nil {
			return nil, fmt.Errorf("invalid deadline format")
		}
		deadline = sql.NullTime{Time: t, Valid: true}
	}

	// Set default type
	taskType := req.Type
	if taskType == "" {
		taskType = "once"
	}

	// 对于挑战任务，如果没设置惩罚值，默认使用奖励值
	penaltyExp := req.PenaltyExp
	penaltyGold := req.PenaltyGold
	if taskType == "challenge" {
		if penaltyExp == 0 {
			penaltyExp = req.RewardExp
		}
		if penaltyGold == 0 {
			penaltyGold = req.RewardGold
		}
	}

	task := &model.Task{
		UserID:             userID,
		Title:              req.Title,
		Description:        req.Description,
		Category:           req.Category,
		Type:               taskType,
		Status:             "active",
		Deadline:           deadline,
		RewardExp:          req.RewardExp,
		RewardGold:         req.RewardGold,
		RewardStrength:     req.RewardStrength,
		RewardIntelligence: req.RewardIntelligence,
		RewardVitality:     req.RewardVitality,
		RewardSpirit:       req.RewardSpirit,
		PenaltyExp:         penaltyExp,
		PenaltyGold:        penaltyGold,
		DailyLimit:         req.DailyLimit,
		TotalLimit:         req.TotalLimit,
		RemindBefore:       req.RemindBefore,
		RemindInterval:     req.RemindInterval,
		CostMental:         req.CostMental,
		CostPhysical:       req.CostPhysical,
	}

	taskID, err := l.svcCtx.TaskModel.Create(task)
	if err != nil {
		return nil, err
	}

	task.ID = taskID
	resp := l.taskToResp(task)
	return &resp, nil
}

func (l *TaskLogic) UpdateTask(ctx context.Context, userID int64, taskID int64, req *types.UpdateTaskReq) (*types.TaskResp, error) {
	// Get existing task
	task, err := l.svcCtx.TaskModel.FindByID(taskID)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, fmt.Errorf("task not found")
	}

	// Check ownership
	if task.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}

	// Update fields if provided
	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.Category != nil {
		task.Category = *req.Category
	}
	if req.Type != nil {
		task.Type = *req.Type
	}
	if req.Deadline != nil {
		if *req.Deadline == "" {
			task.Deadline = sql.NullTime{}
		} else {
			t, err := time.Parse(time.RFC3339, *req.Deadline)
			if err != nil {
				return nil, fmt.Errorf("invalid deadline format")
			}
			task.Deadline = sql.NullTime{Time: t, Valid: true}
		}
	}
	if req.RewardExp != nil {
		task.RewardExp = *req.RewardExp
	}
	if req.RewardGold != nil {
		task.RewardGold = *req.RewardGold
	}
	if req.RewardStrength != nil {
		task.RewardStrength = *req.RewardStrength
	}
	if req.RewardIntelligence != nil {
		task.RewardIntelligence = *req.RewardIntelligence
	}
	if req.RewardVitality != nil {
		task.RewardVitality = *req.RewardVitality
	}
	if req.RewardSpirit != nil {
		task.RewardSpirit = *req.RewardSpirit
	}
	if req.PenaltyExp != nil {
		task.PenaltyExp = *req.PenaltyExp
	}
	if req.PenaltyGold != nil {
		task.PenaltyGold = *req.PenaltyGold
	}
	if req.DailyLimit != nil {
		task.DailyLimit = *req.DailyLimit
	}
	if req.TotalLimit != nil {
		task.TotalLimit = *req.TotalLimit
	}
	if req.RemindBefore != nil {
		task.RemindBefore = *req.RemindBefore
	}
	if req.RemindInterval != nil {
		task.RemindInterval = *req.RemindInterval
	}
	if req.CostMental != nil {
		task.CostMental = *req.CostMental
	}
	if req.CostPhysical != nil {
		task.CostPhysical = *req.CostPhysical
	}

	// Update in database
	if err := l.svcCtx.TaskModel.Update(task); err != nil {
		return nil, err
	}

	updatedResp := l.taskToResp(task)
	return &updatedResp, nil
}

type CompleteTaskResult struct {
	Task      types.TaskResp      `json:"task"`
	Character types.CharacterResp `json:"character"`
	Message   string              `json:"message"`
}

func (l *TaskLogic) CompleteTask(ctx context.Context, userID int64, taskID int64, source string) (*CompleteTaskResult, error) {
	// Get task
	task, err := l.svcCtx.TaskModel.FindByID(taskID)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, fmt.Errorf("task not found")
	}

	// Check ownership
	if task.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}

	// Check if task is already completed/failed
	if task.Status != "active" {
		return nil, fmt.Errorf("task is not active")
	}

	// Get today's date
	today := time.Now().Format("2006-01-02")

	// Handle repeatable tasks: check limits
	if task.Type == "repeatable" {
		// Check total limit
		if task.TotalLimit > 0 && task.CompletedCount >= task.TotalLimit {
			return nil, fmt.Errorf("已达到总完成次数上限")
		}

		// Reset today's count if it's a new day
		if task.LastCompletedDate != today {
			task.TodayCompletionCount = 0
		}

		// Check daily limit
		if task.DailyLimit > 0 && task.TodayCompletionCount >= task.DailyLimit {
			return nil, fmt.Errorf("已达到今日完成次数上限（%d/%d）", task.TodayCompletionCount, task.DailyLimit)
		}

		// Increment counters
		task.CompletedCount++
		task.TodayCompletionCount++
		task.LastCompletedDate = today

		// Mark as completed if total limit reached
		if task.TotalLimit > 0 && task.CompletedCount >= task.TotalLimit {
			task.Status = "completed"
		}
	} else {
		// For once/challenge tasks, mark as completed
		task.Status = "completed"
	}

	// Get character stats
	stats, err := l.svcCtx.CharacterModel.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if stats == nil {
		return nil, fmt.Errorf("character not found")
	}

	// Consume mental/physical power (allow going negative — overdraft)
	stats.MentalPower -= task.CostMental
	stats.PhysicalPower -= task.CostPhysical

	// Accumulate sleep aid: += cost * 0.75
	stats.MentalSleepAid += int(float64(task.CostMental) * 0.75)
	stats.PhysicalSleepAid += int(float64(task.CostPhysical) * 0.75)

	// Sync legacy energy field (average of mental + physical, clamped to 0)
	avgPower := (stats.MentalPower + stats.PhysicalPower) / 2
	if avgPower < 0 {
		avgPower = 0
	}
	stats.Energy = avgPower

	// Apply rewards
	stats.Exp += task.RewardExp
	stats.Gold += task.RewardGold
	stats.Strength += task.RewardStrength
	stats.Intelligence += task.RewardIntelligence
	stats.Vitality += task.RewardVitality
	stats.Spirit += task.RewardSpirit

	// Update last activity date
	stats.LastActivityDate = time.Now().Format("2006-01-02")

	// Handle level up
	CheckAndApplyLevelUp(stats)

	// Recalculate HP
	stats.MaxHP = 100 + int(stats.Strength*2) + int(stats.Vitality*3)
	if stats.HP > stats.MaxHP {
		stats.HP = stats.MaxHP
	}

	// Update task
	if err := l.svcCtx.TaskModel.Update(task); err != nil {
		return nil, err
	}

	// Update character
	if err := l.svcCtx.CharacterModel.Update(stats); err != nil {
		return nil, err
	}

	// Create task log
	log := &model.TaskLog{
		TaskID: taskID,
		UserID: userID,
		Action: "complete",
		Source: source,
	}
	if err := l.svcCtx.TaskModel.CreateLog(log); err != nil {
		return nil, err
	}

	// Build response
	charLogic := NewCharacterLogic(l.svcCtx)
	charResp := charLogic.statsToResp(stats)

	message := fmt.Sprintf("✅ 任务「%s」已完成！获得 %d经验 %d金币", task.Title, task.RewardExp, task.RewardGold)

	return &CompleteTaskResult{
		Task:      l.taskToResp(task),
		Character: *charResp,
		Message:   message,
	}, nil
}

func (l *TaskLogic) FailTask(ctx context.Context, taskID int64, reason string) error {
	// Get task
	task, err := l.svcCtx.TaskModel.FindByID(taskID)
	if err != nil {
		return err
	}
	if task == nil {
		return fmt.Errorf("task not found")
	}

	// Only challenge tasks can fail
	if task.Type != "challenge" {
		return fmt.Errorf("only challenge tasks can fail")
	}

	// Already failed or completed
	if task.Status != "active" {
		return nil // Already processed
	}

	// Get character stats
	stats, err := l.svcCtx.CharacterModel.FindByUserID(task.UserID)
	if err != nil {
		return err
	}
	if stats == nil {
		return fmt.Errorf("character not found")
	}

	// Apply penalties (ensure attributes don't go below minimum values)
	stats.Exp -= task.PenaltyExp
	if stats.Exp < 0 {
		stats.Exp = 0
	}

	stats.Gold -= task.PenaltyGold
	if stats.Gold < 0 {
		stats.Gold = 0
	}

	const minAttribute = 5.0
	stats.Strength = maxFloat(minAttribute, stats.Strength-float64(task.PenaltyExp)/10.0)
	stats.Intelligence = maxFloat(minAttribute, stats.Intelligence-float64(task.PenaltyExp)/10.0)
	stats.Vitality = maxFloat(minAttribute, stats.Vitality-float64(task.PenaltyExp)/10.0)
	stats.Spirit = maxFloat(minAttribute, stats.Spirit-float64(task.PenaltyExp)/10.0)

	// Update task status
	task.Status = "failed"

	// Update task
	if err := l.svcCtx.TaskModel.Update(task); err != nil {
		return err
	}

	// Update character
	if err := l.svcCtx.CharacterModel.Update(stats); err != nil {
		return err
	}

	// Create task log
	log := &model.TaskLog{
		TaskID: taskID,
		UserID: task.UserID,
		Action: "fail",
		Source: "system",
	}
	if err := l.svcCtx.TaskModel.CreateLog(log); err != nil {
		return err
	}

	fmt.Printf("❌ Task #%d failed: %s. Penalties applied: -%d exp, -%d gold\n",
		taskID, reason, task.PenaltyExp, task.PenaltyGold)

	return nil
}

func (l *TaskLogic) DeleteTask(ctx context.Context, userID int64, taskID int64, source string) error {
	// Get task
	task, err := l.svcCtx.TaskModel.FindByID(taskID)
	if err != nil {
		return err
	}
	if task == nil {
		return fmt.Errorf("task not found")
	}

	// Check ownership
	if task.UserID != userID {
		return fmt.Errorf("unauthorized")
	}

	// Only active tasks can be deleted
	if task.Status != "active" {
		return fmt.Errorf("只能删除进行中的任务")
	}

	// Soft delete
	if err := l.svcCtx.TaskModel.Delete(taskID); err != nil {
		return err
	}

	if source == "" {
		source = "web"
	}

	// Create task log
	log := &model.TaskLog{
		TaskID: taskID,
		UserID: userID,
		Action: "delete",
		Source: source,
	}
	return l.svcCtx.TaskModel.CreateLog(log)
}

func (l *TaskLogic) taskToResp(task *model.Task) types.TaskResp {
	var deadline, lastRemindedAt *string

	if task.Deadline.Valid {
		deadlineStr := task.Deadline.Time.Format(time.RFC3339)
		deadline = &deadlineStr
	}

	if task.LastRemindedAt.Valid {
		lastRemindedStr := task.LastRemindedAt.Time.Format(time.RFC3339)
		lastRemindedAt = &lastRemindedStr
	}

	return types.TaskResp{
		ID:                   task.ID,
		UserID:               task.UserID,
		Title:                task.Title,
		Description:          task.Description,
		Category:             task.Category,
		Type:                 task.Type,
		Status:               task.Status,
		Deadline:             deadline,
		RewardExp:            task.RewardExp,
		RewardGold:           task.RewardGold,
		RewardStrength:       task.RewardStrength,
		RewardIntelligence:   task.RewardIntelligence,
		RewardVitality:       task.RewardVitality,
		RewardSpirit:         task.RewardSpirit,
		PenaltyExp:           task.PenaltyExp,
		PenaltyGold:          task.PenaltyGold,
		DailyLimit:           task.DailyLimit,
		TotalLimit:           task.TotalLimit,
		CompletedCount:       task.CompletedCount,
		TodayCompletionCount: task.TodayCompletionCount,
		LastCompletedDate:    task.LastCompletedDate,
		RemindBefore:         task.RemindBefore,
		RemindInterval:       task.RemindInterval,
		LastRemindedAt:       lastRemindedAt,
		CreatedAt:            task.CreatedAt.Format(time.RFC3339),
		UpdatedAt:            task.UpdatedAt.Format(time.RFC3339),
		CostMental:           task.CostMental,
		CostPhysical:         task.CostPhysical,
	}
}

// TelegramTaskCompleter is an adapter for telegram bot to complete tasks
// This avoids circular dependency between telegram and logic packages
type TelegramTaskCompleter struct {
	svcCtx *svc.ServiceContext
}

func NewTelegramTaskCompleter(svcCtx *svc.ServiceContext) *TelegramTaskCompleter {
	return &TelegramTaskCompleter{svcCtx: svcCtx}
}

func (t *TelegramTaskCompleter) CompleteTask(userID int64, taskID int64) (expGained int, goldGained int, newLevel int, newExp int, err error) {
	logic := NewTaskLogic(t.svcCtx)
	result, err := logic.CompleteTask(context.Background(), userID, taskID, "telegram")
	if err != nil {
		return 0, 0, 0, 0, err
	}

	// Get task to return reward amounts
	task, _ := t.svcCtx.TaskModel.FindByID(taskID)
	if task != nil {
		return task.RewardExp, task.RewardGold, result.Character.Level, result.Character.Exp, nil
	}

	return 0, 0, result.Character.Level, result.Character.Exp, nil
}

func (t *TelegramTaskCompleter) DeleteTask(userID int64, taskID int64) error {
	logic := NewTaskLogic(t.svcCtx)
	return logic.DeleteTask(context.Background(), userID, taskID, "telegram")
}

// Helper function for max of two floats
func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
