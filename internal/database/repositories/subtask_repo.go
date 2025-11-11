package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/stiffis/UniCLI/internal/models"
)

// SubtaskRepository handles subtask data operations.
type SubtaskRepository struct {
	*BaseRepository
}

// NewSubtaskRepository creates a new subtask repository.
func NewSubtaskRepository(db *sql.DB) *SubtaskRepository {
	return &SubtaskRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// FindByTaskID retrieves all subtasks for a given task ID.
func (r *SubtaskRepository) FindByTaskID(taskID string) ([]models.Subtask, error) {
	query := `
		SELECT id, task_id, title, is_completed, created_at
		FROM subtasks
		WHERE task_id = ?
		ORDER BY created_at ASC
	`
	rows, err := r.DB().Query(query, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to query subtasks: %w", err)
	}
	defer rows.Close()

	var subtasks []models.Subtask
	for rows.Next() {
		var subtask models.Subtask
		err := rows.Scan(
			&subtask.ID,
			&subtask.TaskID,
			&subtask.Title,
			&subtask.IsCompleted,
			&subtask.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subtask: %w", err)
		}
		subtasks = append(subtasks, subtask)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating subtasks: %w", err)
	}

	return subtasks, nil
}

// Create inserts a new subtask into the database.
func (r *SubtaskRepository) Create(subtask *models.Subtask) error {
	query := `
		INSERT INTO subtasks (task_id, title, is_completed, created_at)
		VALUES (?, ?, ?, ?)
	`
	subtask.CreatedAt = time.Now()

	result, err := r.DB().Exec(
		query,
		subtask.TaskID,
		subtask.Title,
		subtask.IsCompleted,
		subtask.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create subtask: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID for subtask: %w", err)
	}
	subtask.ID = int(id)

	return nil
}
