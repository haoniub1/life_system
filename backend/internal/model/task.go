package model

import (
	"database/sql"
	"fmt"
	"time"
)

type Task struct {
	ID                   int64
	UserID               int64
	Title                string
	Description          string
	Category             string
	Type                 string // once, repeatable, challenge
	Status               string // active, completed, failed, deleted
	Deadline             sql.NullTime
	RewardExp            int
	RewardGold           int
	RewardStrength       float64
	RewardIntelligence   float64
	RewardVitality       float64
	RewardSpirit         float64
	PenaltyExp           int
	PenaltyGold          int
	DailyLimit           int
	TotalLimit           int
	CompletedCount       int
	TodayCompletionCount int    // Today's completion count for repeatable tasks
	LastCompletedDate    string // Last completion date (YYYY-MM-DD) for daily reset
	RemindBefore         int    // minutes
	RemindInterval       int    // minutes
	LastRemindedAt       sql.NullTime
	CreatedAt            time.Time
	UpdatedAt            time.Time
	CostMental           int // Mental power cost when completing
	CostPhysical         int // Physical power cost when completing
}

type TaskLog struct {
	ID        int64
	TaskID    int64
	UserID    int64
	Action    string // complete, fail, delete
	Source    string // web, telegram
	CreatedAt time.Time
}

type TaskModel struct {
	db *sql.DB
}

func NewTaskModel(db *sql.DB) *TaskModel {
	return &TaskModel{db: db}
}

func (m *TaskModel) FindByUserID(userID int64, taskType, status string) ([]*Task, error) {
	query := `
		SELECT id, user_id, title, description, category, type, status, deadline,
		       reward_exp, reward_gold, reward_strength, reward_intelligence, reward_vitality, reward_spirit,
		       penalty_exp, penalty_gold, daily_limit, total_limit, completed_count,
		       today_completion_count, last_completed_date,
		       remind_before, remind_interval, last_reminded_at, created_at, updated_at,
		       cost_mental, cost_physical
		FROM tasks WHERE user_id = ?
	`
	args := []interface{}{userID}

	if taskType != "" {
		query += ` AND type = ?`
		args = append(args, taskType)
	}

	if status != "" {
		query += ` AND status = ?`
		args = append(args, status)
	} else {
		query += ` AND status != 'deleted'`
	}

	query += ` ORDER BY created_at DESC`

	rows, err := m.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var task Task
		err := rows.Scan(
			&task.ID, &task.UserID, &task.Title, &task.Description, &task.Category, &task.Type,
			&task.Status, &task.Deadline, &task.RewardExp, &task.RewardGold,
			&task.RewardStrength, &task.RewardIntelligence, &task.RewardVitality, &task.RewardSpirit,
			&task.PenaltyExp, &task.PenaltyGold, &task.DailyLimit, &task.TotalLimit, &task.CompletedCount,
			&task.TodayCompletionCount, &task.LastCompletedDate,
			&task.RemindBefore, &task.RemindInterval, &task.LastRemindedAt, &task.CreatedAt, &task.UpdatedAt,
			&task.CostMental, &task.CostPhysical,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	return tasks, rows.Err()
}

func (m *TaskModel) FindByID(id int64) (*Task, error) {
	var task Task
	err := m.db.QueryRow(`
		SELECT id, user_id, title, description, category, type, status, deadline,
		       reward_exp, reward_gold, reward_strength, reward_intelligence, reward_vitality, reward_spirit,
		       penalty_exp, penalty_gold, daily_limit, total_limit, completed_count,
		       today_completion_count, last_completed_date,
		       remind_before, remind_interval, last_reminded_at, created_at, updated_at,
		       cost_mental, cost_physical
		FROM tasks WHERE id = ?
	`, id).Scan(
		&task.ID, &task.UserID, &task.Title, &task.Description, &task.Category, &task.Type,
		&task.Status, &task.Deadline, &task.RewardExp, &task.RewardGold,
		&task.RewardStrength, &task.RewardIntelligence, &task.RewardVitality, &task.RewardSpirit,
		&task.PenaltyExp, &task.PenaltyGold, &task.DailyLimit, &task.TotalLimit, &task.CompletedCount,
		&task.TodayCompletionCount, &task.LastCompletedDate,
		&task.RemindBefore, &task.RemindInterval, &task.LastRemindedAt, &task.CreatedAt, &task.UpdatedAt,
		&task.CostMental, &task.CostPhysical,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &task, nil
}

func (m *TaskModel) Create(task *Task) (int64, error) {
	result, err := m.db.Exec(`
		INSERT INTO tasks (user_id, title, description, category, type, status, deadline,
		                   reward_exp, reward_gold, reward_strength, reward_intelligence, reward_vitality, reward_spirit,
		                   penalty_exp, penalty_gold, daily_limit, total_limit, completed_count,
		                   today_completion_count, last_completed_date,
		                   remind_before, remind_interval, cost_mental, cost_physical, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))
	`,
		task.UserID, task.Title, task.Description, task.Category, task.Type, task.Status, task.Deadline,
		task.RewardExp, task.RewardGold, task.RewardStrength, task.RewardIntelligence, task.RewardVitality, task.RewardSpirit,
		task.PenaltyExp, task.PenaltyGold, task.DailyLimit, task.TotalLimit, task.CompletedCount,
		task.TodayCompletionCount, task.LastCompletedDate,
		task.RemindBefore, task.RemindInterval, task.CostMental, task.CostPhysical,
	)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (m *TaskModel) Update(task *Task) error {
	_, err := m.db.Exec(`
		UPDATE tasks
		SET title = ?, description = ?, category = ?, type = ?, status = ?, deadline = ?,
		    reward_exp = ?, reward_gold = ?, reward_strength = ?, reward_intelligence = ?, reward_vitality = ?, reward_spirit = ?,
		    penalty_exp = ?, penalty_gold = ?, daily_limit = ?, total_limit = ?, completed_count = ?,
		    today_completion_count = ?, last_completed_date = ?,
		    remind_before = ?, remind_interval = ?, last_reminded_at = ?, cost_mental = ?, cost_physical = ?, updated_at = datetime('now')
		WHERE id = ?
	`,
		task.Title, task.Description, task.Category, task.Type, task.Status, task.Deadline,
		task.RewardExp, task.RewardGold, task.RewardStrength, task.RewardIntelligence, task.RewardVitality, task.RewardSpirit,
		task.PenaltyExp, task.PenaltyGold, task.DailyLimit, task.TotalLimit, task.CompletedCount,
		task.TodayCompletionCount, task.LastCompletedDate,
		task.RemindBefore, task.RemindInterval, task.LastRemindedAt, task.CostMental, task.CostPhysical, task.ID,
	)

	return err
}

func (m *TaskModel) UpdateStatus(id int64, status string) error {
	_, err := m.db.Exec(`
		UPDATE tasks SET status = ?, updated_at = datetime('now') WHERE id = ?
	`, status, id)

	return err
}

func (m *TaskModel) Delete(id int64) error {
	_, err := m.db.Exec(`
		UPDATE tasks SET status = 'deleted', updated_at = datetime('now') WHERE id = ?
	`, id)

	return err
}

type TaskWithUser struct {
	Task     *Task
	TgChatID int64
	Username string
}

func (m *TaskModel) FindTasksNeedingReminder() ([]*TaskWithUser, error) {
	rows, err := m.db.Query(`
		SELECT t.id, t.user_id, t.title, t.description, t.category, t.type, t.status, t.deadline,
		       t.reward_exp, t.reward_gold, t.reward_strength, t.reward_intelligence, t.reward_vitality, t.reward_spirit,
		       t.penalty_exp, t.penalty_gold, t.daily_limit, t.total_limit, t.completed_count,
		       t.today_completion_count, t.last_completed_date,
		       t.remind_before, t.remind_interval, t.last_reminded_at, t.created_at, t.updated_at,
		       t.cost_mental, t.cost_physical,
		       u.tg_chat_id, u.username
		FROM tasks t
		JOIN users u ON t.user_id = u.id
		WHERE t.status = 'active' AND t.deadline IS NOT NULL AND t.remind_before > 0 AND u.tg_chat_id > 0
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*TaskWithUser
	for rows.Next() {
		var task Task
		var tgChatID int64
		var username string

		err := rows.Scan(
			&task.ID, &task.UserID, &task.Title, &task.Description, &task.Category, &task.Type,
			&task.Status, &task.Deadline, &task.RewardExp, &task.RewardGold,
			&task.RewardStrength, &task.RewardIntelligence, &task.RewardVitality, &task.RewardSpirit,
			&task.PenaltyExp, &task.PenaltyGold, &task.DailyLimit, &task.TotalLimit, &task.CompletedCount,
			&task.TodayCompletionCount, &task.LastCompletedDate,
			&task.RemindBefore, &task.RemindInterval, &task.LastRemindedAt, &task.CreatedAt, &task.UpdatedAt,
			&task.CostMental, &task.CostPhysical,
			&tgChatID, &username,
		)
		if err != nil {
			return nil, err
		}

		results = append(results, &TaskWithUser{
			Task:     &task,
			TgChatID: tgChatID,
			Username: username,
		})
	}

	return results, rows.Err()
}

func (m *TaskModel) UpdateLastReminded(id int64, t time.Time) error {
	_, err := m.db.Exec(`
		UPDATE tasks SET last_reminded_at = ?, updated_at = datetime('now') WHERE id = ?
	`, t, id)

	return err
}

func (m *TaskModel) CreateLog(log *TaskLog) error {
	_, err := m.db.Exec(`
		INSERT INTO task_logs (task_id, user_id, action, source, created_at)
		VALUES (?, ?, ?, ?, datetime('now'))
	`, log.TaskID, log.UserID, log.Action, log.Source)

	return err
}

// ResetDailyCompletionCounts resets today_completion_count for all repeatable tasks
// that haven't been completed today (last_completed_date != today)
func (m *TaskModel) ResetDailyCompletionCounts(today string) error {
	fmt.Printf("🔄 Resetting daily completion counts for date: %s\n", today)

	result, err := m.db.Exec(`
		UPDATE tasks
		SET today_completion_count = 0, updated_at = datetime('now')
		WHERE type = 'repeatable'
		  AND status = 'active'
		  AND (last_completed_date IS NULL OR last_completed_date = '' OR last_completed_date != ?)
	`, today)

	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	fmt.Printf("✅ Reset %d tasks\n", rows)

	return nil
}

// FindExpiredChallengeTasks finds all active challenge tasks that have passed their deadline
func (m *TaskModel) FindExpiredChallengeTasks() ([]*Task, error) {
	rows, err := m.db.Query(`
		SELECT id, user_id, title, description, category, type, status, deadline,
		       reward_exp, reward_gold, reward_strength, reward_intelligence, reward_vitality, reward_spirit,
		       penalty_exp, penalty_gold, daily_limit, total_limit, completed_count,
		       today_completion_count, last_completed_date,
		       remind_before, remind_interval, last_reminded_at, created_at, updated_at,
		       cost_mental, cost_physical
		FROM tasks
		WHERE type = 'challenge'
		  AND status = 'active'
		  AND deadline IS NOT NULL
		  AND deadline < datetime('now')
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var task Task
		err := rows.Scan(
			&task.ID, &task.UserID, &task.Title, &task.Description, &task.Category, &task.Type,
			&task.Status, &task.Deadline, &task.RewardExp, &task.RewardGold,
			&task.RewardStrength, &task.RewardIntelligence, &task.RewardVitality, &task.RewardSpirit,
			&task.PenaltyExp, &task.PenaltyGold, &task.DailyLimit, &task.TotalLimit, &task.CompletedCount,
			&task.TodayCompletionCount, &task.LastCompletedDate,
			&task.RemindBefore, &task.RemindInterval, &task.LastRemindedAt, &task.CreatedAt, &task.UpdatedAt,
			&task.CostMental, &task.CostPhysical,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	return tasks, rows.Err()
}
