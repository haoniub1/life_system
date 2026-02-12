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
	PrimaryAttribute     string
	Difficulty           int
	RewardExp            int
	RewardSpiritStones   int
	RewardPhysique       float64
	RewardWillpower      float64
	RewardIntelligence   float64
	RewardPerception     float64
	RewardCharisma       float64
	RewardAgility        float64
	FatigueCost          int
	PenaltyExp           int
	PenaltySpiritStones  int
	DailyLimit           int
	TotalLimit           int
	CompletedCount       int
	TodayCompletionCount int
	LastCompletedDate    string
	RemindBefore         int // minutes
	RemindInterval       int // minutes
	LastRemindedAt       sql.NullTime
	SortOrder            int
	CreatedAt            time.Time
	UpdatedAt            time.Time
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

// taskColumns is the shared column list for task queries.
const taskColumns = `id, user_id, title, description, category, type, status, deadline,
       primary_attribute, difficulty,
       reward_exp, reward_spirit_stones,
       reward_physique, reward_willpower, reward_intelligence,
       reward_perception, reward_charisma, reward_agility,
       fatigue_cost, penalty_exp, penalty_spirit_stones,
       daily_limit, total_limit, completed_count,
       today_completion_count, last_completed_date,
       remind_before, remind_interval, last_reminded_at,
       COALESCE(sort_order, 0) as sort_order, created_at, updated_at`

// taskColumnsAliased is the same column list prefixed with "t." for use in JOIN queries.
const taskColumnsAliased = `t.id, t.user_id, t.title, t.description, t.category, t.type, t.status, t.deadline,
       t.primary_attribute, t.difficulty,
       t.reward_exp, t.reward_spirit_stones,
       t.reward_physique, t.reward_willpower, t.reward_intelligence,
       t.reward_perception, t.reward_charisma, t.reward_agility,
       t.fatigue_cost, t.penalty_exp, t.penalty_spirit_stones,
       t.daily_limit, t.total_limit, t.completed_count,
       t.today_completion_count, t.last_completed_date,
       t.remind_before, t.remind_interval, t.last_reminded_at,
       COALESCE(t.sort_order, 0) as sort_order, t.created_at, t.updated_at`

func scanTask(scanner interface{ Scan(...interface{}) error }) (*Task, error) {
	var task Task
	err := scanner.Scan(
		&task.ID, &task.UserID, &task.Title, &task.Description, &task.Category, &task.Type,
		&task.Status, &task.Deadline,
		&task.PrimaryAttribute, &task.Difficulty,
		&task.RewardExp, &task.RewardSpiritStones,
		&task.RewardPhysique, &task.RewardWillpower, &task.RewardIntelligence,
		&task.RewardPerception, &task.RewardCharisma, &task.RewardAgility,
		&task.FatigueCost, &task.PenaltyExp, &task.PenaltySpiritStones,
		&task.DailyLimit, &task.TotalLimit, &task.CompletedCount,
		&task.TodayCompletionCount, &task.LastCompletedDate,
		&task.RemindBefore, &task.RemindInterval, &task.LastRemindedAt,
		&task.SortOrder, &task.CreatedAt, &task.UpdatedAt,
	)
	return &task, err
}

func (m *TaskModel) FindByUserID(userID int64, taskType, status string) ([]*Task, error) {
	query := `SELECT ` + taskColumns + ` FROM tasks WHERE user_id = ?`
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

	query += ` ORDER BY COALESCE(sort_order, 0) ASC, created_at DESC`

	rows, err := m.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		task, err := scanTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, rows.Err()
}

func (m *TaskModel) FindByID(id int64) (*Task, error) {
	row := m.db.QueryRow(`SELECT `+taskColumns+` FROM tasks WHERE id = ?`, id)
	task, err := scanTask(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return task, nil
}

func (m *TaskModel) Create(task *Task) (int64, error) {
	result, err := m.db.Exec(`
		INSERT INTO tasks (user_id, title, description, category, type, status, deadline,
		                   primary_attribute, difficulty,
		                   reward_exp, reward_spirit_stones,
		                   reward_physique, reward_willpower, reward_intelligence,
		                   reward_perception, reward_charisma, reward_agility,
		                   fatigue_cost, penalty_exp, penalty_spirit_stones,
		                   daily_limit, total_limit, completed_count,
		                   today_completion_count, last_completed_date,
		                   remind_before, remind_interval, sort_order, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?,
		        ?, ?,
		        ?, ?,
		        ?, ?, ?,
		        ?, ?, ?,
		        ?, ?, ?,
		        ?, ?, ?,
		        ?, ?,
		        ?, ?, ?, datetime('now'), datetime('now'))
	`,
		task.UserID, task.Title, task.Description, task.Category, task.Type, task.Status, task.Deadline,
		task.PrimaryAttribute, task.Difficulty,
		task.RewardExp, task.RewardSpiritStones,
		task.RewardPhysique, task.RewardWillpower, task.RewardIntelligence,
		task.RewardPerception, task.RewardCharisma, task.RewardAgility,
		task.FatigueCost, task.PenaltyExp, task.PenaltySpiritStones,
		task.DailyLimit, task.TotalLimit, task.CompletedCount,
		task.TodayCompletionCount, task.LastCompletedDate,
		task.RemindBefore, task.RemindInterval, task.SortOrder,
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
		    primary_attribute = ?, difficulty = ?,
		    reward_exp = ?, reward_spirit_stones = ?,
		    reward_physique = ?, reward_willpower = ?, reward_intelligence = ?,
		    reward_perception = ?, reward_charisma = ?, reward_agility = ?,
		    fatigue_cost = ?, penalty_exp = ?, penalty_spirit_stones = ?,
		    daily_limit = ?, total_limit = ?, completed_count = ?,
		    today_completion_count = ?, last_completed_date = ?,
		    remind_before = ?, remind_interval = ?, last_reminded_at = ?,
		    sort_order = ?, updated_at = datetime('now')
		WHERE id = ?
	`,
		task.Title, task.Description, task.Category, task.Type, task.Status, task.Deadline,
		task.PrimaryAttribute, task.Difficulty,
		task.RewardExp, task.RewardSpiritStones,
		task.RewardPhysique, task.RewardWillpower, task.RewardIntelligence,
		task.RewardPerception, task.RewardCharisma, task.RewardAgility,
		task.FatigueCost, task.PenaltyExp, task.PenaltySpiritStones,
		task.DailyLimit, task.TotalLimit, task.CompletedCount,
		task.TodayCompletionCount, task.LastCompletedDate,
		task.RemindBefore, task.RemindInterval, task.LastRemindedAt,
		task.SortOrder, task.ID,
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
		SELECT `+taskColumnsAliased+`,
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
		var tgChatID int64
		var username string
		var task Task

		err := rows.Scan(
			&task.ID, &task.UserID, &task.Title, &task.Description, &task.Category, &task.Type,
			&task.Status, &task.Deadline,
			&task.PrimaryAttribute, &task.Difficulty,
			&task.RewardExp, &task.RewardSpiritStones,
			&task.RewardPhysique, &task.RewardWillpower, &task.RewardIntelligence,
			&task.RewardPerception, &task.RewardCharisma, &task.RewardAgility,
			&task.FatigueCost, &task.PenaltyExp, &task.PenaltySpiritStones,
			&task.DailyLimit, &task.TotalLimit, &task.CompletedCount,
			&task.TodayCompletionCount, &task.LastCompletedDate,
			&task.RemindBefore, &task.RemindInterval, &task.LastRemindedAt,
			&task.SortOrder, &task.CreatedAt, &task.UpdatedAt,
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

	affected, _ := result.RowsAffected()
	fmt.Printf("✅ Reset %d tasks\n", affected)

	return nil
}

// ReorderTasks sets sort_order for a list of task IDs belonging to a user.
func (m *TaskModel) ReorderTasks(userID int64, taskIDs []int64) error {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`UPDATE tasks SET sort_order = ?, updated_at = datetime('now') WHERE id = ? AND user_id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i, id := range taskIDs {
		result, err := stmt.Exec(i+1, id, userID)
		if err != nil {
			return err
		}
		affected, _ := result.RowsAffected()
		if affected == 0 {
			return fmt.Errorf("task %d not found or unauthorized", id)
		}
	}

	return tx.Commit()
}

// FindExpiredChallengeTasks finds all active challenge tasks that have passed their deadline
func (m *TaskModel) FindExpiredChallengeTasks() ([]*Task, error) {
	rows, err := m.db.Query(`
		SELECT `+taskColumns+`
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
		task, err := scanTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, rows.Err()
}
