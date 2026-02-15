package logic

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"life-system-backend/internal/model"
	"life-system-backend/internal/realm"
	"life-system-backend/internal/svc"
	"life-system-backend/internal/types"
)

// Difficulty template: maps difficulty stars (0-5) to preset values
type difficultyPreset struct {
	Fatigue      int
	SpiritStones int
	AttrBonus    float64
}

var difficultyTable = map[int]difficultyPreset{
	0: {Fatigue: 1, SpiritStones: 10, AttrBonus: 0},
	1: {Fatigue: 5, SpiritStones: 50, AttrBonus: 0.1},
	2: {Fatigue: 10, SpiritStones: 120, AttrBonus: 0.2},
	3: {Fatigue: 20, SpiritStones: 300, AttrBonus: 0.4},
	4: {Fatigue: 40, SpiritStones: 800, AttrBonus: 0.7},
	5: {Fatigue: 90, SpiritStones: 2500, AttrBonus: 1.0},
}

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
	if req.Title == "" {
		return nil, fmt.Errorf("title is required")
	}

	var deadline sql.NullTime
	if req.Deadline != "" {
		t, err := time.Parse(time.RFC3339, req.Deadline)
		if err != nil {
			return nil, fmt.Errorf("invalid deadline format")
		}
		deadline = sql.NullTime{Time: t, Valid: true}
	}

	taskType := req.Type
	if taskType == "" {
		taskType = "once"
	}

	// For challenge tasks, default penalty to reward values
	penaltyExp := req.PenaltyExp
	penaltySpiritStones := req.PenaltySpiritStones
	if taskType == "challenge" {
		if penaltyExp == 0 {
			penaltyExp = req.RewardExp
		}
		if penaltySpiritStones == 0 {
			penaltySpiritStones = req.RewardSpiritStones
		}
	}

	fatigueCost := req.FatigueCost
	if fatigueCost == 0 {
		fatigueCost = 10
	}

	task := &model.Task{
		UserID:              userID,
		Title:               req.Title,
		Description:         req.Description,
		Category:            req.Category,
		Type:                taskType,
		Status:              "active",
		Deadline:            deadline,
		PrimaryAttribute:    req.PrimaryAttribute,
		Difficulty:          req.Difficulty,
		RewardExp:           req.RewardExp,
		RewardSpiritStones:  req.RewardSpiritStones,
		RewardPhysique:      req.RewardPhysique,
		RewardWillpower:     req.RewardWillpower,
		RewardIntelligence:  req.RewardIntelligence,
		RewardPerception:    req.RewardPerception,
		RewardCharisma:      req.RewardCharisma,
		RewardAgility:       req.RewardAgility,
		FatigueCost:         fatigueCost,
		PenaltyExp:          penaltyExp,
		PenaltySpiritStones: penaltySpiritStones,
		DailyLimit:          req.DailyLimit,
		TotalLimit:          req.TotalLimit,
		RemindBefore:        req.RemindBefore,
		RemindInterval:      req.RemindInterval,
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
	task, err := l.svcCtx.TaskModel.FindByID(taskID)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, fmt.Errorf("task not found")
	}
	if task.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}

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
	if req.PrimaryAttribute != nil {
		task.PrimaryAttribute = *req.PrimaryAttribute
	}
	if req.Difficulty != nil {
		task.Difficulty = *req.Difficulty
	}
	if req.RewardExp != nil {
		task.RewardExp = *req.RewardExp
	}
	if req.RewardSpiritStones != nil {
		task.RewardSpiritStones = *req.RewardSpiritStones
	}
	if req.RewardPhysique != nil {
		task.RewardPhysique = *req.RewardPhysique
	}
	if req.RewardWillpower != nil {
		task.RewardWillpower = *req.RewardWillpower
	}
	if req.RewardIntelligence != nil {
		task.RewardIntelligence = *req.RewardIntelligence
	}
	if req.RewardPerception != nil {
		task.RewardPerception = *req.RewardPerception
	}
	if req.RewardCharisma != nil {
		task.RewardCharisma = *req.RewardCharisma
	}
	if req.RewardAgility != nil {
		task.RewardAgility = *req.RewardAgility
	}
	if req.PenaltyExp != nil {
		task.PenaltyExp = *req.PenaltyExp
	}
	if req.PenaltySpiritStones != nil {
		task.PenaltySpiritStones = *req.PenaltySpiritStones
	}
	if req.FatigueCost != nil {
		task.FatigueCost = *req.FatigueCost
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
	task, err := l.svcCtx.TaskModel.FindByID(taskID)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, fmt.Errorf("task not found")
	}
	if task.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}
	if task.Status != "active" {
		return nil, fmt.Errorf("task is not active")
	}

	today := time.Now().Format("2006-01-02")

	// Handle repeatable tasks: check limits
	if task.Type == "repeatable" {
		if task.TotalLimit > 0 && task.CompletedCount >= task.TotalLimit {
			return nil, fmt.Errorf("已达到总完成次数上限")
		}
		if task.LastCompletedDate != today {
			task.TodayCompletionCount = 0
		}
		if task.DailyLimit > 0 && task.TodayCompletionCount >= task.DailyLimit {
			return nil, fmt.Errorf("已达到今日完成次数上限（%d/%d）", task.TodayCompletionCount, task.DailyLimit)
		}
		task.CompletedCount++
		task.TodayCompletionCount++
		task.LastCompletedDate = today
		if task.TotalLimit > 0 && task.CompletedCount >= task.TotalLimit {
			task.Status = "completed"
		}
	} else {
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

	// Consume fatigue (allow overdraft)
	stats.Fatigue += task.FatigueCost
	if stats.Fatigue > stats.FatigueCap {
		overdraft := float64(stats.Fatigue - stats.FatigueCap)
		stats.OverdraftPenalty += overdraft
	}

	// Add spirit stones
	stats.SpiritStones += task.RewardSpiritStones

	// Update last activity date
	stats.LastActivityDate = today

	// Get attributes for reward processing
	attrs, err := l.svcCtx.CharacterModel.FindAttributesByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Build map for quick lookup
	attrMap := make(map[string]*model.CharacterAttribute)
	for _, a := range attrs {
		attrMap[a.AttrKey] = a
	}

	// Apply attribute rewards using realm.ProcessAttrGain
	attrRewards := map[string]float64{
		"physique":     task.RewardPhysique,
		"willpower":    task.RewardWillpower,
		"intelligence": task.RewardIntelligence,
		"perception":   task.RewardPerception,
		"charisma":     task.RewardCharisma,
		"agility":      task.RewardAgility,
	}

	for key, gain := range attrRewards {
		if gain <= 0 {
			continue
		}
		attr, ok := attrMap[key]
		if !ok {
			continue
		}

		result := realm.ProcessAttrGain(attr.Value, gain, attr.Realm, attr.RealmExp, attr.IsBottleneck, attr.AccumulationPool)
		attr.Value = result.NewValue
		attr.AccumulationPool = result.NewAccPool
		attr.RealmExp = result.NewRealmExp
		attr.IsBottleneck = result.NewIsBottleneck

		if err := l.svcCtx.CharacterModel.UpdateAttribute(attr); err != nil {
			return nil, err
		}
	}

	// Reload attributes after updates
	attrs, err = l.svcCtx.CharacterModel.FindAttributesByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Update task
	if err := l.svcCtx.TaskModel.Update(task); err != nil {
		return nil, err
	}

	// Update character stats
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

	charLogic := NewCharacterLogic(l.svcCtx)
	charResp := charLogic.statsToResp(stats, attrs)

	message := fmt.Sprintf("✅ 任务「%s」已完成！获得 %d灵石", task.Title, task.RewardSpiritStones)

	return &CompleteTaskResult{
		Task:      l.taskToResp(task),
		Character: *charResp,
		Message:   message,
	}, nil
}

func (l *TaskLogic) FailTask(ctx context.Context, taskID int64, reason string) error {
	task, err := l.svcCtx.TaskModel.FindByID(taskID)
	if err != nil {
		return err
	}
	if task == nil {
		return fmt.Errorf("task not found")
	}
	if task.Type != "challenge" {
		return fmt.Errorf("only challenge tasks can fail")
	}
	if task.Status != "active" {
		return nil
	}

	// Get character stats
	stats, err := l.svcCtx.CharacterModel.FindByUserID(task.UserID)
	if err != nil {
		return err
	}
	if stats == nil {
		return fmt.Errorf("character not found")
	}

	// Deduct spirit stones (min 0)
	stats.SpiritStones -= task.PenaltySpiritStones
	if stats.SpiritStones < 0 {
		stats.SpiritStones = 0
	}

	// Get attributes for penalty processing
	attrs, err := l.svcCtx.CharacterModel.FindAttributesByUserID(task.UserID)
	if err != nil {
		return err
	}

	// Attribute penalty: reduce by penaltyExp/10 but not below realm base value
	attrPenalty := float64(task.PenaltyExp) / 10.0
	if attrPenalty > 0 {
		for _, attr := range attrs {
			if attr.AttrKey == "luck" {
				continue
			}
			minVal := realm.AttrMin(attr.Realm)
			attr.Value -= attrPenalty
			if attr.Value < minVal {
				attr.Value = minVal
			}
			if err := l.svcCtx.CharacterModel.UpdateAttribute(attr); err != nil {
				return err
			}
		}
	}

	// Update task status
	task.Status = "failed"

	if err := l.svcCtx.TaskModel.Update(task); err != nil {
		return err
	}
	if err := l.svcCtx.CharacterModel.Update(stats); err != nil {
		return err
	}

	log := &model.TaskLog{
		TaskID: taskID,
		UserID: task.UserID,
		Action: "fail",
		Source: "system",
	}
	if err := l.svcCtx.TaskModel.CreateLog(log); err != nil {
		return err
	}

	fmt.Printf("❌ Task #%d failed: %s. Penalties applied: -%d spiritStones\n",
		taskID, reason, task.PenaltySpiritStones)

	return nil
}

func (l *TaskLogic) DeleteTask(ctx context.Context, userID int64, taskID int64, source string) error {
	task, err := l.svcCtx.TaskModel.FindByID(taskID)
	if err != nil {
		return err
	}
	if task == nil {
		return fmt.Errorf("task not found")
	}
	if task.UserID != userID {
		return fmt.Errorf("unauthorized")
	}
	if task.Status != "active" {
		return fmt.Errorf("只能删除进行中的任务")
	}

	if err := l.svcCtx.TaskModel.Delete(taskID); err != nil {
		return err
	}

	if source == "" {
		source = "web"
	}

	log := &model.TaskLog{
		TaskID: taskID,
		UserID: userID,
		Action: "delete",
		Source: source,
	}
	return l.svcCtx.TaskModel.CreateLog(log)
}

type QuickTaskResult struct {
	Task      types.TaskResp      `json:"task"`
	Character types.CharacterResp `json:"character,omitempty"`
	Message   string              `json:"message"`
	Completed bool                `json:"completed"` // true if auto-completed (once), false if just created (repeatable/challenge)
}

func (l *TaskLogic) QuickComplete(ctx context.Context, userID int64, req *types.QuickTaskReq) (*QuickTaskResult, error) {
	// Validate difficulty
	preset, ok := difficultyTable[req.Difficulty]
	if !ok {
		return nil, fmt.Errorf("invalid difficulty: %d (must be 0-5)", req.Difficulty)
	}

	// Validate categories (can be either attribute keys or tags)
	validAttrs := map[string]bool{
		"physique": true, "willpower": true, "intelligence": true,
		"perception": true, "charisma": true, "agility": true,
	}
	
	// Attribute-to-default-tag mapping
	attrToTag := map[string]string{
		"physique":     "运动",
		"willpower":    "专注",
		"intelligence": "学习",
		"perception":   "观察",
		"charisma":     "社交",
		"agility":      "灵活",
	}
	
	// Deduplicate categories and collect both attrs and tags
	seen := make(map[string]bool)
	categoryTags := []string{}
	
	for _, cat := range req.Categories {
		// If it's an attribute key, use its default tag
		if validAttrs[cat] {
			tag := attrToTag[cat]
			if !seen[tag] {
				categoryTags = append(categoryTags, tag)
				seen[tag] = true
			}
			seen[cat] = true
		} else {
			// It's a custom tag, use directly
			if !seen[cat] {
				categoryTags = append(categoryTags, cat)
				seen[cat] = true
			}
		}
	}

	// Validate task type
	taskType := req.Type
	if taskType == "" {
		taskType = "once"
	}
	if taskType != "once" && taskType != "repeatable" && taskType != "challenge" {
		return nil, fmt.Errorf("invalid type: %s (must be once, repeatable, or challenge)", taskType)
	}

	// Challenge requires deadline
	if taskType == "challenge" && req.Deadline == "" {
		return nil, fmt.Errorf("challenge tasks require a deadline")
	}

	// Build title
	title := req.Title
	if title == "" {
		title = fmt.Sprintf("快速任务 (★%d)", req.Difficulty)
	}

	// Build attribute rewards from categories
	var rewardPhysique, rewardWillpower, rewardIntelligence float64
	var rewardPerception, rewardCharisma, rewardAgility float64
	for cat := range seen {
		switch cat {
		case "physique":
			rewardPhysique = preset.AttrBonus
		case "willpower":
			rewardWillpower = preset.AttrBonus
		case "intelligence":
			rewardIntelligence = preset.AttrBonus
		case "perception":
			rewardPerception = preset.AttrBonus
		case "charisma":
			rewardCharisma = preset.AttrBonus
		case "agility":
			rewardAgility = preset.AttrBonus
		}
	}

	source := req.Source
	if source == "" {
		source = "api"
	}

	// Create the task
	categoryStr := ""
	if len(categoryTags) > 0 {
		categoryStr = categoryTags[0]
		for i := 1; i < len(categoryTags); i++ {
			categoryStr += "," + categoryTags[i]
		}
	}
	
	createReq := &types.CreateTaskReq{
		Title:              title,
		Category:           categoryStr,
		Type:               taskType,
		Difficulty:         req.Difficulty,
		FatigueCost:        preset.Fatigue,
		RewardSpiritStones: preset.SpiritStones,
		RewardPhysique:     rewardPhysique,
		RewardWillpower:    rewardWillpower,
		RewardIntelligence: rewardIntelligence,
		RewardPerception:   rewardPerception,
		RewardCharisma:     rewardCharisma,
		RewardAgility:      rewardAgility,
		DailyLimit:         req.DailyLimit,
		TotalLimit:         req.TotalLimit,
		Deadline:           req.Deadline,
	}

	taskResp, err := l.CreateTask(ctx, userID, createReq)
	if err != nil {
		return nil, fmt.Errorf("create task failed: %w", err)
	}

	// For "once" type: auto-complete immediately
	if taskType == "once" {
		result, err := l.CompleteTask(ctx, userID, taskResp.ID, source)
		if err != nil {
			return nil, fmt.Errorf("complete task failed: %w", err)
		}
		return &QuickTaskResult{
			Task:      result.Task,
			Character: result.Character,
			Message:   result.Message,
			Completed: true,
		}, nil
	}

	// For repeatable/challenge: just create, return task info
	message := fmt.Sprintf("✅ 任务「%s」已创建", title)
	if taskType == "repeatable" {
		message += "（可重复）"
	} else {
		message += "（挑战）"
	}

	return &QuickTaskResult{
		Task:      *taskResp,
		Message:   message,
		Completed: false,
	}, nil
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
		PrimaryAttribute:     task.PrimaryAttribute,
		Difficulty:           task.Difficulty,
		RewardExp:            task.RewardExp,
		RewardSpiritStones:   task.RewardSpiritStones,
		RewardPhysique:       task.RewardPhysique,
		RewardWillpower:      task.RewardWillpower,
		RewardIntelligence:   task.RewardIntelligence,
		RewardPerception:     task.RewardPerception,
		RewardCharisma:       task.RewardCharisma,
		RewardAgility:        task.RewardAgility,
		FatigueCost:          task.FatigueCost,
		PenaltyExp:           task.PenaltyExp,
		PenaltySpiritStones:  task.PenaltySpiritStones,
		DailyLimit:           task.DailyLimit,
		TotalLimit:           task.TotalLimit,
		CompletedCount:       task.CompletedCount,
		TodayCompletionCount: task.TodayCompletionCount,
		LastCompletedDate:    task.LastCompletedDate,
		RemindBefore:         task.RemindBefore,
		RemindInterval:       task.RemindInterval,
		LastRemindedAt:       lastRemindedAt,
		SortOrder:            task.SortOrder,
		CreatedAt:            task.CreatedAt.Format(time.RFC3339),
		UpdatedAt:            task.UpdatedAt.Format(time.RFC3339),
	}
}

// TelegramTaskCompleter is an adapter for telegram bot to complete tasks
type TelegramTaskCompleter struct {
	svcCtx *svc.ServiceContext
}

func NewTelegramTaskCompleter(svcCtx *svc.ServiceContext) *TelegramTaskCompleter {
	return &TelegramTaskCompleter{svcCtx: svcCtx}
}

func (t *TelegramTaskCompleter) CompleteTask(userID int64, taskID int64) (expGained int, spiritStonesGained int, realmTitle string, spiritStones int, err error) {
	logic := NewTaskLogic(t.svcCtx)
	result, err := logic.CompleteTask(context.Background(), userID, taskID, "telegram")
	if err != nil {
		return 0, 0, "", 0, err
	}

	task, _ := t.svcCtx.TaskModel.FindByID(taskID)
	if task != nil {
		return task.RewardExp, task.RewardSpiritStones, result.Character.Title, result.Character.SpiritStones, nil
	}

	return 0, 0, result.Character.Title, result.Character.SpiritStones, nil
}

func (t *TelegramTaskCompleter) DeleteTask(userID int64, taskID int64) error {
	logic := NewTaskLogic(t.svcCtx)
	return logic.DeleteTask(context.Background(), userID, taskID, "telegram")
}
