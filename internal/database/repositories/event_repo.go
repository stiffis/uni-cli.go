package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/stiffis/UniCLI/internal/models"
)

// EventRepository handles event data operations
type EventRepository struct {
	*BaseRepository
}

// NewEventRepository creates a new event repository
func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create inserts a new event into the database
func (r *EventRepository) Create(event *models.Event) error {
	query := `
		INSERT INTO events (
			id, title, description, start_datetime, end_datetime, type, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.DB().Exec(
		query,
		event.ID,
		event.Title,
		event.Description,
		event.StartDatetime,
		event.EndDatetime,
		event.Type,
		event.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}

	return nil
}

// FindByID retrieves an event by its ID
func (r *EventRepository) FindByID(id string) (*models.Event, error) {
	query := `
		SELECT id, title, description, start_datetime, end_datetime, type, created_at
		FROM events
		WHERE id = ?
	`

	event := &models.Event{}
	var endDatetime sql.NullTime

	err := r.DB().QueryRow(query, id).Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.StartDatetime,
		&endDatetime,
		&event.Type,
		&event.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("event not found: %s", id)
		}
		return nil, fmt.Errorf("failed to find event: %w", err)
	}

	if endDatetime.Valid {
		event.EndDatetime = &endDatetime.Time
	}

	return event, nil
}

// GetEventsByMonth retrieves all events for a given month and year
func (r *EventRepository) GetEventsByMonth(year int, month time.Month) ([]models.Event, error) {
	startDate := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second) // Last second of the month

	query := `
		SELECT id, title, description, start_datetime, end_datetime, type, created_at
		FROM events
		WHERE start_datetime >= ? AND start_datetime <= ?
		ORDER BY start_datetime ASC
	`

	rows, err := r.DB().Query(query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query events for month %d/%d: %w", month, year, err)
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		var endDatetime sql.NullTime

		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.StartDatetime,
			&endDatetime,
			&event.Type,
			&event.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}

		if endDatetime.Valid {
			event.EndDatetime = &endDatetime.Time
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating events: %w", err)
	}

	return events, nil
}

// Update updates an existing event
func (r *EventRepository) Update(event *models.Event) error {
	query := `
		UPDATE events
		SET title = ?, description = ?, start_datetime = ?, end_datetime = ?, type = ?
		WHERE id = ?
	`

	result, err := r.DB().Exec(
		query,
		event.Title,
		event.Description,
		event.StartDatetime,
		event.EndDatetime,
		event.Type,
		event.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("event not found: %s", event.ID)
	}

	return nil
}

// Delete removes an event from the database
func (r *EventRepository) Delete(id string) error {
	query := `DELETE FROM events WHERE id = ?`

	result, err := r.DB().Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("event not found: %s", id)
	}

	return nil
}
