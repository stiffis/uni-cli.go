package models

import (
	"time"

	"github.com/google/uuid"
)

// Course represents an academic course/class
type Course struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`        // "Calculus I"
	Code        string           `json:"code"`        // "MATH 101"
	Professor   string           `json:"professor"`   // "Dr. Smith"
	Location    string           `json:"location"`    // "Room 301"
	Semester    string           `json:"semester"`    // "Fall 2025"
	Credits     int              `json:"credits"`     // 3
	Color       string           `json:"color"`       // For calendar
	Description string           `json:"description"` // Course description
	Schedule    []CourseSchedule `json:"schedule"`    // Weekly schedule
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

// CourseSchedule represents a recurring class schedule
type CourseSchedule struct {
	ID        string    `json:"id"`
	CourseID  string    `json:"course_id"`
	DayOfWeek int       `json:"day_of_week"` // 1=Monday, 7=Sunday
	StartTime string    `json:"start_time"`  // "09:00"
	EndTime   string    `json:"end_time"`    // "10:30"
	CreatedAt time.Time `json:"created_at"`
}

// CourseNote represents notes taken for a course
type CourseNote struct {
	ID        string    `json:"id"`
	CourseID  string    `json:"course_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"` // Markdown content
	Date      time.Time `json:"date"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CourseAttendance tracks attendance for course sessions
type CourseAttendance struct {
	ID        string    `json:"id"`
	CourseID  string    `json:"course_id"`
	Date      time.Time `json:"date"`
	Status    string    `json:"status"` // "present", "absent", "late"
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
}

// NewCourse creates a new course with generated ID
func NewCourse(name string) *Course {
	return &Course{
		ID:        uuid.New().String(),
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// NewCourseSchedule creates a new course schedule
func NewCourseSchedule(courseID string, dayOfWeek int, startTime, endTime string) *CourseSchedule {
	return &CourseSchedule{
		ID:        uuid.New().String(),
		CourseID:  courseID,
		DayOfWeek: dayOfWeek,
		StartTime: startTime,
		EndTime:   endTime,
		CreatedAt: time.Now(),
	}
}

// NewCourseNote creates a new course note
func NewCourseNote(courseID, title, content string) *CourseNote {
	return &CourseNote{
		ID:        uuid.New().String(),
		CourseID:  courseID,
		Title:     title,
		Content:   content,
		Date:      time.Now(),
		Tags:      []string{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// NewCourseAttendance creates a new attendance record
func NewCourseAttendance(courseID string, date time.Time, status string) *CourseAttendance {
	return &CourseAttendance{
		ID:        uuid.New().String(),
		CourseID:  courseID,
		Date:      date,
		Status:    status,
		CreatedAt: time.Now(),
	}
}

func (cs *CourseSchedule) DayOfWeekString() string {
	days := []string{"", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	if cs.DayOfWeek >= 1 && cs.DayOfWeek <= 7 {
		return days[cs.DayOfWeek]
	}
	return ""
}

func (cs *CourseSchedule) DayOfWeekShort() string {
	days := []string{"", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	if cs.DayOfWeek >= 1 && cs.DayOfWeek <= 7 {
		return days[cs.DayOfWeek]
	}
	return ""
}

func (c *Course) GenerateEventsForWeek(weekStart time.Time) []*Event {
	var events []*Event

	weekStart = getStartOfWeek(weekStart)

	for _, schedule := range c.Schedule {
		event := c.generateEventFromSchedule(schedule, weekStart)
		if event != nil {
			events = append(events, event)
		}
	}

	return events
}

func (c *Course) GenerateEventsForMonth(year int, month time.Month) []*Event {
	var events []*Event

	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	lastDay := firstDay.AddDate(0, 1, -1)

	for d := firstDay; d.Before(lastDay) || d.Equal(lastDay); d = d.AddDate(0, 0, 7) {
		weekEvents := c.GenerateEventsForWeek(d)
		for _, event := range weekEvents {
			if event.StartDatetime.Month() == month {
				events = append(events, event)
			}
		}
	}

	return events
}

func (c *Course) GenerateEventsForDateRange(start, end time.Time) []*Event {
	var events []*Event

	for d := start; d.Before(end); d = d.AddDate(0, 0, 1) {
		for _, schedule := range c.Schedule {
			dayOfWeek := int(d.Weekday())
			if dayOfWeek == 0 {
				dayOfWeek = 7
			}

			if schedule.DayOfWeek == dayOfWeek {
				event := c.generateEventFromScheduleForDate(schedule, d)
				if event != nil {
					events = append(events, event)
				}
			}
		}
	}

	return events
}

func (c *Course) generateEventFromSchedule(schedule CourseSchedule, weekStart time.Time) *Event {
	targetDay := weekStart
	currentWeekday := int(weekStart.Weekday())
	if currentWeekday == 0 {
		currentWeekday = 7
	}

	daysToAdd := schedule.DayOfWeek - currentWeekday
	if daysToAdd < 0 {
		return nil
	}
	targetDay = targetDay.AddDate(0, 0, daysToAdd)

	return c.generateEventFromScheduleForDate(schedule, targetDay)
}

func (c *Course) generateEventFromScheduleForDate(schedule CourseSchedule, date time.Time) *Event {
	startTime, err := time.Parse("15:04", schedule.StartTime)
	if err != nil {
		return nil
	}
	endTime, err := time.Parse("15:04", schedule.EndTime)
	if err != nil {
		return nil
	}

	start := time.Date(date.Year(), date.Month(), date.Day(),
		startTime.Hour(), startTime.Minute(), 0, 0, time.Local)
	end := time.Date(date.Year(), date.Month(), date.Day(),
		endTime.Hour(), endTime.Minute(), 0, 0, time.Local)

	event := NewEvent(c.Name, start)
	event.Type = "class"
	event.EndDatetime = &end
	event.Description = c.Code
	if c.Location != "" {
		event.Description += "\n" + c.Location
	}
	if c.Professor != "" {
		event.Description += "\n" + c.Professor
	}

	return event
}

func getStartOfWeek(date time.Time) time.Time {
	weekday := int(date.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	daysToMonday := weekday - 1
	monday := date.AddDate(0, 0, -daysToMonday)
	return time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, time.Local)
}
