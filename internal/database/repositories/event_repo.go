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

func (r *EventRepository) Create(event *models.Event) error {
	query := `
		INSERT INTO events (
			id, title, description, start_datetime, end_datetime, type, category_id,
			recurrence_rule, recurrence_end_date, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.DB().Exec(
		query,
		event.ID,
		event.Title,
		event.Description,
		event.StartDatetime,
		event.EndDatetime,
		event.Type,
		event.CategoryID,
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
		SELECT id, title, description, start_datetime, end_datetime, type, category_id,
			   recurrence_rule, recurrence_end_date, created_at
		FROM events
		WHERE id = ?
	`

	event := &models.Event{}
	var endDatetime, recurrenceEndDate sql.NullTime
	var recurrenceRule, categoryID sql.NullString

	err := r.DB().QueryRow(query, id).Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.StartDatetime,
		&endDatetime,
		&event.Type,
		&categoryID,
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
	if categoryID.Valid {
		event.CategoryID = categoryID.String
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
	allEvents, err := r.FindAll()
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
			if event.StartDatetime.Year() == year && event.StartDatetime.Month() == month {
				eventsForMonth = append(eventsForMonth, event)
			}
		}
	}

	return eventsForMonth, nil
}

// FindAll retrieves all events from the database
func (r *EventRepository) FindAll() ([]models.Event, error) {
	query := `
		SELECT id, title, description, start_datetime, end_datetime, type, category_id,
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
		var recurrenceRule, categoryID sql.NullString

		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.StartDatetime,
			&endDatetime,
			&event.Type,
			&categoryID,
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
		if categoryID.Valid {
			event.CategoryID = categoryID.String
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
	
	// If no end date, set a limit to 2 years from now to avoid infinite loops
	limitDate := time.Now().AddDate(2, 0, 0)
	if event.RecurrenceEndDate != nil {
		limitDate = *event.RecurrenceEndDate
	}


	for {
		if currentTime.After(limitDate) {
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

func (r *EventRepository) Update(event *models.Event) error {
	query := `
		UPDATE events
		SET title = ?, description = ?, start_datetime = ?, end_datetime = ?, type = ?, category_id = ?,
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
		event.CategoryID,
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

// GetEventsWithCoursesForMonth gets all events AND course classes for a specific month
func (r *EventRepository) GetEventsWithCoursesForMonth(year int, month time.Month, courseRepo *CourseRepository) ([]models.Event, error) {
	events, err := r.GetEventsByMonth(year, month)
	if err != nil {
		return nil, err
	}

	courses, err := courseRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// Generate class events from courses
	for _, course := range courses {
		classEvents := course.GenerateEventsForMonth(year, month)
		for _, classEvent := range classEvents {
			if course.Color != "" {
				// We'll need to create a category for this or use the color directly
				classEvent.CategoryID = "course_" + course.ID
			}
			events = append(events, *classEvent)
		}
	}

	return events, nil
}

// GetEventsWithCoursesForWeek gets all events AND course classes for a specific week
func (r *EventRepository) GetEventsWithCoursesForWeek(weekStart time.Time, courseRepo *CourseRepository) ([]models.Event, error) {
	weekEnd := weekStart.AddDate(0, 0, 7)

	allEvents, err := r.FindAll()
	if err != nil {
		return nil, err
	}

	var events []models.Event
	for _, event := range allEvents {
		if event.RecurrenceRule != "" && event.RecurrenceRule != "none" {
			// Generate occurrences for recurring events within the week
			occurrences := generateOccurrencesForRange(event, weekStart, weekEnd)
			events = append(events, occurrences...)
		} else {
			if (event.StartDatetime.Equal(weekStart) || event.StartDatetime.After(weekStart)) &&
				event.StartDatetime.Before(weekEnd) {
				events = append(events, event)
			}
		}
	}

	courses, err := courseRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// Generate class events from courses
	for _, course := range courses {
		classEvents := course.GenerateEventsForWeek(weekStart)
		for _, classEvent := range classEvents {
			if course.Color != "" {
				classEvent.CategoryID = "course_" + course.ID
			}
			events = append(events, *classEvent)
		}
	}

	return events, nil
}

// GetEventsWithCoursesForDay gets all events AND course classes for a specific day
func (r *EventRepository) GetEventsWithCoursesForDay(date time.Time, courseRepo *CourseRepository) ([]models.Event, error) {
	// Normalize date to start of day
	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
	dayEnd := dayStart.AddDate(0, 0, 1)

	allEvents, err := r.FindAll()
	if err != nil {
		return nil, err
	}

	var events []models.Event
	for _, event := range allEvents {
		if event.RecurrenceRule != "" && event.RecurrenceRule != "none" {
			// Generate occurrences for recurring events within the day
			occurrences := generateOccurrencesForRange(event, dayStart, dayEnd)
			events = append(events, occurrences...)
		} else {
			if (event.StartDatetime.Equal(dayStart) || event.StartDatetime.After(dayStart)) &&
				event.StartDatetime.Before(dayEnd) {
				events = append(events, event)
			}
		}
	}

	courses, err := courseRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// Generate class events from courses
	for _, course := range courses {
		classEvents := course.GenerateEventsForDateRange(dayStart, dayEnd)
		for _, classEvent := range classEvents {
			if course.Color != "" {
				classEvent.CategoryID = "course_" + course.ID
			}
			events = append(events, *classEvent)
		}
	}

	return events, nil
}

// generateOccurrencesForRange generates event occurrences within a specific date range
func generateOccurrencesForRange(event models.Event, start, end time.Time) []models.Event {
	var occurrences []models.Event

	currentTime := event.StartDatetime
	limitDate := time.Now().AddDate(2, 0, 0)
	if event.RecurrenceEndDate != nil {
		limitDate = *event.RecurrenceEndDate
	}

	for currentTime.Before(limitDate) && currentTime.Before(end) {
		if (currentTime.Equal(start) || currentTime.After(start)) && currentTime.Before(end) {
			occurrence := event
			occurrence.ID = uuid.New().String()
			occurrence.StartDatetime = currentTime
			
			if event.EndDatetime != nil {
				duration := event.EndDatetime.Sub(event.StartDatetime)
				newEnd := currentTime.Add(duration)
				occurrence.EndDatetime = &newEnd
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
			break
		}

		if currentTime.After(limitDate) {
			break
		}
	}

	return occurrences
}
