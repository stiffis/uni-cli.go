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

// DayOfWeekString returns the day name
func (cs *CourseSchedule) DayOfWeekString() string {
	days := []string{"", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	if cs.DayOfWeek >= 1 && cs.DayOfWeek <= 7 {
		return days[cs.DayOfWeek]
	}
	return ""
}

// DayOfWeekShort returns the short day name
func (cs *CourseSchedule) DayOfWeekShort() string {
	days := []string{"", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	if cs.DayOfWeek >= 1 && cs.DayOfWeek <= 7 {
		return days[cs.DayOfWeek]
	}
	return ""
}
