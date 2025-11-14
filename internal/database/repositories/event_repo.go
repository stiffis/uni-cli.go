package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
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
			id, title, description, start_datetime, end_datetime, type, 
			recurrence_rule, recurrence_end_date, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.DB().Exec(
		query,
		event.ID,
		event.Title,
		event.Description,
		event.StartDatetime,
		event.EndDatetime,
		event.Type,
		event.RecurrenceRule,
		event.RecurrenceEndDate,
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
		SELECT id, title, description, start_datetime, end_datetime, type, 
			   recurrence_rule, recurrence_end_date, created_at
		FROM events
		WHERE id = ?
	`

	event := &models.Event{}
	var endDatetime, recurrenceEndDate sql.NullTime
	var recurrenceRule sql.NullString

	err := r.DB().QueryRow(query, id).Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.StartDatetime,
		&endDatetime,
		&event.Type,
		&recurrenceRule,
		&recurrenceEndDate,
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
	if recurrenceRule.Valid {
		event.RecurrenceRule = recurrenceRule.String
	}
	if recurrenceEndDate.Valid {
		event.RecurrenceEndDate = &recurrenceEndDate.Time
	}

	return event, nil
}

// GetEventsByMonth retrieves all events for a given month and year
func (r *EventRepository) GetEventsByMonth(year int, month time.Month) ([]models.Event, error) {
	// 1. Fetch all events from the database (this could be optimized)
	allEvents, err := r.findAllEvents()
	if err != nil {
		return nil, err
	}

	// 2. Generate occurrences for the given month
	var eventsForMonth []models.Event
	for _, event := range allEvents {
		if event.RecurrenceRule != "" && event.RecurrenceRule != "none" {
			// Generate occurrences for recurring events
			occurrences := generateOccurrences(event, year, month)
			eventsForMonth = append(eventsForMonth, occurrences...)
		} else {
			// Add non-recurring events that fall within the month
			if event.StartDatetime.Year() == year && event.StartDatetime.Month() == month {
				eventsForMonth = append(eventsForMonth, event)
			}
		}
	}

	return eventsForMonth, nil
}

// findAllEvents retrieves all events from the database
func (r *EventRepository) findAllEvents() ([]models.Event, error) {
	query := `
		SELECT id, title, description, start_datetime, end_datetime, type, 
			   recurrence_rule, recurrence_end_date, created_at
		FROM events
		ORDER BY start_datetime ASC
	`

	rows, err := r.DB().Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query all events: %w", err)
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		var endDatetime, recurrenceEndDate sql.NullTime
		var recurrenceRule sql.NullString

		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.StartDatetime,
			&endDatetime,
			&event.Type,
			&recurrenceRule,
			&recurrenceEndDate,
			&event.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}

		if endDatetime.Valid {
			event.EndDatetime = &endDatetime.Time
		}
		if recurrenceRule.Valid {
			event.RecurrenceRule = recurrenceRule.String
		}
		if recurrenceEndDate.Valid {
			event.RecurrenceEndDate = &recurrenceEndDate.Time
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating events: %w", err)
	}

	return events, nil
}

// generateOccurrences generates event occurrences for a given month
func generateOccurrences(event models.Event, year int, month time.Month) []models.Event {
	var occurrences []models.Event

	// Loop from the start of the event until the recurrence end date or the end of the month
	currentTime := event.StartDatetime
	for {
		if event.RecurrenceEndDate != nil && currentTime.After(*event.RecurrenceEndDate) {
			break
		}

		if currentTime.Year() > year || (currentTime.Year() == year && currentTime.Month() > month) {
			break
		}

		if currentTime.Year() == year && currentTime.Month() == month {
			occurrence := event
			occurrence.ID = uuid.New().String() // Give each occurrence a unique ID
			occurrence.StartDatetime = currentTime
			if event.EndDatetime != nil {
				duration := event.EndDatetime.Sub(event.StartDatetime)
				newEndDate := currentTime.Add(duration)
				occurrence.EndDatetime = &newEndDate
			}
			occurrences = append(occurrences, occurrence)
		}

		switch event.RecurrenceRule {
		case "daily":
			currentTime = currentTime.AddDate(0, 0, 1)
		case "weekly":
			currentTime = currentTime.AddDate(0, 0, 7)
		case "monthly":
			currentTime = currentTime.AddDate(0, 1, 0)
		default:
			return occurrences // No valid recurrence rule
		}
	}

	return occurrences
}

// Update updates an existing event
func (r *EventRepository) Update(event *models.Event) error {
	query := `
		UPDATE events
		SET title = ?, description = ?, start_datetime = ?, end_datetime = ?, type = ?,
			recurrence_rule = ?, recurrence_end_date = ?
		WHERE id = ?
	`

	result, err := r.DB().Exec(
		query,
		event.Title,
		event.Description,
		event.StartDatetime,
		event.EndDatetime,
		event.Type,
		event.RecurrenceRule,
		event.RecurrenceEndDate,
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
