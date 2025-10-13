package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/stiffis/UniCLI/internal/models"
)

// TaskRepository handles task data operations
type TaskRepository struct {
	*BaseRepository
}

// NewTaskRepository creates a new task repository
func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create inserts a new task into the database
func (r *TaskRepository) Create(task *models.Task) error {
	query := `
		INSERT INTO tasks (
			id, title, description, status, priority, category,
			due_date, created_at, updated_at, completed_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.DB().Exec(
		query,
		task.ID,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.Category,
		task.DueDate,
		task.CreatedAt,
		task.UpdatedAt,
		task.CompletedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	// Insert tags if any
	if len(task.Tags) > 0 {
		if err := r.updateTags(task.ID, task.Tags); err != nil {
			return fmt.Errorf("failed to create task tags: %w", err)
		}
	}

	return nil
}

// FindByID retrieves a task by its ID
func (r *TaskRepository) FindByID(id string) (*models.Task, error) {
	query := `
		SELECT id, title, description, status, priority, category,
			   due_date, created_at, updated_at, completed_at
		FROM tasks
		WHERE id = ?
	`

	task := &models.Task{}
	var dueDate, completedAt sql.NullTime

	err := r.DB().QueryRow(query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.Category,
		&dueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
		&completedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found: %s", id)
		}
		return nil, fmt.Errorf("failed to find task: %w", err)
	}

	// Handle nullable fields
	if dueDate.Valid {
		task.DueDate = &dueDate.Time
	}
	if completedAt.Valid {
		task.CompletedAt = &completedAt.Time
	}

	// Load tags
	tags, err := r.loadTags(task.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to load tags: %w", err)
	}
	task.Tags = tags

	return task, nil
}

// FindAll retrieves all tasks
func (r *TaskRepository) FindAll() ([]models.Task, error) {
	query := `
		SELECT id, title, description, status, priority, category,
			   due_date, created_at, updated_at, completed_at
		FROM tasks
		ORDER BY created_at DESC
	`

	rows, err := r.DB().Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	return r.scanTasks(rows)
}

// FindByStatus retrieves tasks by status
func (r *TaskRepository) FindByStatus(status models.TaskStatus) ([]models.Task, error) {
	query := `
		SELECT id, title, description, status, priority, category,
			   due_date, created_at, updated_at, completed_at
		FROM tasks
		WHERE status = ?
		ORDER BY created_at DESC
	`

	rows, err := r.DB().Query(query, status)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks by status: %w", err)
	}
	defer rows.Close()

	return r.scanTasks(rows)
}

// FindDueToday retrieves tasks due today
func (r *TaskRepository) FindDueToday() ([]models.Task, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	query := `
		SELECT id, title, description, status, priority, category,
			   due_date, created_at, updated_at, completed_at
		FROM tasks
		WHERE due_date >= ? AND due_date < ?
		ORDER BY due_date ASC
	`

	rows, err := r.DB().Query(query, startOfDay, endOfDay)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks due today: %w", err)
	}
	defer rows.Close()

	return r.scanTasks(rows)
}

// FindUpcoming retrieves tasks due in the next 7 days (excluding today)
func (r *TaskRepository) FindUpcoming() ([]models.Task, error) {
	now := time.Now()
	tomorrow := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Add(24 * time.Hour)
	nextWeek := tomorrow.Add(7 * 24 * time.Hour)

	query := `
		SELECT id, title, description, status, priority, category,
			   due_date, created_at, updated_at, completed_at
		FROM tasks
		WHERE due_date >= ? AND due_date < ? AND status != ?
		ORDER BY due_date ASC
	`

	rows, err := r.DB().Query(query, tomorrow, nextWeek, models.TaskStatusCompleted)
	if err != nil {
		return nil, fmt.Errorf("failed to query upcoming tasks: %w", err)
	}
	defer rows.Close()

	return r.scanTasks(rows)
}

// FindOverdue retrieves overdue tasks
func (r *TaskRepository) FindOverdue() ([]models.Task, error) {
	now := time.Now()

	query := `
		SELECT id, title, description, status, priority, category,
			   due_date, created_at, updated_at, completed_at
		FROM tasks
		WHERE due_date < ? AND status != ?
		ORDER BY due_date ASC
	`

	rows, err := r.DB().Query(query, now, models.TaskStatusCompleted)
	if err != nil {
		return nil, fmt.Errorf("failed to query overdue tasks: %w", err)
	}
	defer rows.Close()

	return r.scanTasks(rows)
}

// Update updates an existing task
func (r *TaskRepository) Update(task *models.Task) error {
	task.UpdatedAt = time.Now()

	// Set CompletedAt if status is completed and it's not already set
	if task.Status == models.TaskStatusCompleted && task.CompletedAt == nil {
		now := time.Now()
		task.CompletedAt = &now
	}

	// Clear CompletedAt if status is not completed
	if task.Status != models.TaskStatusCompleted {
		task.CompletedAt = nil
	}

	query := `
		UPDATE tasks
		SET title = ?, description = ?, status = ?, priority = ?,
		    category = ?, due_date = ?, updated_at = ?, completed_at = ?
		WHERE id = ?
	`

	result, err := r.DB().Exec(
		query,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.Category,
		task.DueDate,
		task.UpdatedAt,
		task.CompletedAt,
		task.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("task not found: %s", task.ID)
	}

	// Update tags
	if err := r.updateTags(task.ID, task.Tags); err != nil {
		return fmt.Errorf("failed to update task tags: %w", err)
	}

	return nil
}

// Delete removes a task from the database
func (r *TaskRepository) Delete(id string) error {
	query := `DELETE FROM tasks WHERE id = ?`

	result, err := r.DB().Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("task not found: %s", id)
	}

	return nil
}

// ToggleComplete toggles the completion status of a task
func (r *TaskRepository) ToggleComplete(id string) error {
	task, err := r.FindByID(id)
	if err != nil {
		return err
	}

	if task.Status == models.TaskStatusCompleted {
		// Mark as pending
		task.Status = models.TaskStatusPending
		task.CompletedAt = nil
	} else {
		// Mark as completed
		task.Status = models.TaskStatusCompleted
		now := time.Now()
		task.CompletedAt = &now
	}

	return r.Update(task)
}

// scanTasks scans multiple tasks from query rows
func (r *TaskRepository) scanTasks(rows *sql.Rows) ([]models.Task, error) {
	var tasks []models.Task

	for rows.Next() {
		var task models.Task
		var dueDate, completedAt sql.NullTime

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.Category,
			&dueDate,
			&task.CreatedAt,
			&task.UpdatedAt,
			&completedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}

		// Handle nullable fields
		if dueDate.Valid {
			task.DueDate = &dueDate.Time
		}
		if completedAt.Valid {
			task.CompletedAt = &completedAt.Time
		}

		// Load tags
		tags, err := r.loadTags(task.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to load tags for task %s: %w", task.ID, err)
		}
		task.Tags = tags

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tasks: %w", err)
	}

	return tasks, nil
}

// loadTags loads tags for a task
func (r *TaskRepository) loadTags(taskID string) ([]string, error) {
	query := `
		SELECT t.name
		FROM tags t
		JOIN task_tags tt ON t.id = tt.tag_id
		WHERE tt.task_id = ?
		ORDER BY t.name
	`

	rows, err := r.DB().Query(query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, rows.Err()
}

// updateTags updates tags for a task
func (r *TaskRepository) updateTags(taskID string, tags []string) error {
	// Start transaction
	tx, err := r.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete existing task_tags
	if _, err := tx.Exec("DELETE FROM task_tags WHERE task_id = ?", taskID); err != nil {
		return err
	}

	// Insert new tags
	for _, tag := range tags {
		// Get or create tag
		var tagID int64
		err := tx.QueryRow("SELECT id FROM tags WHERE name = ?", tag).Scan(&tagID)
		if err == sql.ErrNoRows {
			// Create new tag
			result, err := tx.Exec("INSERT INTO tags (name) VALUES (?)", tag)
			if err != nil {
				return err
			}
			tagID, err = result.LastInsertId()
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}

		// Link task with tag
		if _, err := tx.Exec("INSERT INTO task_tags (task_id, tag_id) VALUES (?, ?)", taskID, tagID); err != nil {
			return err
		}
	}

	return tx.Commit()
}
