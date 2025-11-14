package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/stiffis/UniCLI/internal/models"
)

// CourseRepository handles course database operations
type CourseRepository struct {
	db *sql.DB
}

// NewCourseRepository creates a new course repository
func NewCourseRepository(db *sql.DB) *CourseRepository {
	return &CourseRepository{db: db}
}

// Create creates a new course
func (r *CourseRepository) Create(course *models.Course) error {
	query := `
		INSERT INTO courses (id, name, code, professor, location, semester, credits, color, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query,
		course.ID,
		course.Name,
		course.Code,
		course.Professor,
		course.Location,
		course.Semester,
		course.Credits,
		course.Color,
		course.Description,
		course.CreatedAt,
		course.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create course: %w", err)
	}

	// Create schedules
	for i := range course.Schedule {
		if err := r.CreateSchedule(&course.Schedule[i]); err != nil {
			return err
		}
	}

	return nil
}

// Update updates an existing course
func (r *CourseRepository) Update(course *models.Course) error {
	course.UpdatedAt = time.Now()

	query := `
		UPDATE courses
		SET name = ?, code = ?, professor = ?, location = ?, semester = ?, 
		    credits = ?, color = ?, description = ?, updated_at = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(query,
		course.Name,
		course.Code,
		course.Professor,
		course.Location,
		course.Semester,
		course.Credits,
		course.Color,
		course.Description,
		course.UpdatedAt,
		course.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update course: %w", err)
	}

	// Delete existing schedules and recreate
	if err := r.DeleteSchedules(course.ID); err != nil {
		return err
	}

	for i := range course.Schedule {
		course.Schedule[i].CourseID = course.ID
		if err := r.CreateSchedule(&course.Schedule[i]); err != nil {
			return err
		}
	}

	return nil
}

// Delete deletes a course
func (r *CourseRepository) Delete(id string) error {
	query := "DELETE FROM courses WHERE id = ?"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete course: %w", err)
	}
	return nil
}

// GetByID retrieves a course by ID
func (r *CourseRepository) GetByID(id string) (*models.Course, error) {
	query := `
		SELECT id, name, code, professor, location, semester, credits, color, description, created_at, updated_at
		FROM courses
		WHERE id = ?
	`
	course := &models.Course{}
	err := r.db.QueryRow(query, id).Scan(
		&course.ID,
		&course.Name,
		&course.Code,
		&course.Professor,
		&course.Location,
		&course.Semester,
		&course.Credits,
		&course.Color,
		&course.Description,
		&course.CreatedAt,
		&course.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("course not found")
		}
		return nil, fmt.Errorf("failed to get course: %w", err)
	}

	// Load schedules
	schedules, err := r.GetSchedules(course.ID)
	if err != nil {
		return nil, err
	}
	course.Schedule = schedules

	return course, nil
}

// GetAll retrieves all courses
func (r *CourseRepository) GetAll() ([]models.Course, error) {
	query := `
		SELECT id, name, code, professor, location, semester, credits, color, description, created_at, updated_at
		FROM courses
		ORDER BY name ASC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get courses: %w", err)
	}
	defer rows.Close()

	var courses []models.Course
	for rows.Next() {
		var course models.Course
		err := rows.Scan(
			&course.ID,
			&course.Name,
			&course.Code,
			&course.Professor,
			&course.Location,
			&course.Semester,
			&course.Credits,
			&course.Color,
			&course.Description,
			&course.CreatedAt,
			&course.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan course: %w", err)
		}

		// Load schedules
		schedules, err := r.GetSchedules(course.ID)
		if err != nil {
			return nil, err
		}
		course.Schedule = schedules

		courses = append(courses, course)
	}

	return courses, nil
}

// GetBySemester retrieves courses by semester
func (r *CourseRepository) GetBySemester(semester string) ([]models.Course, error) {
	query := `
		SELECT id, name, code, professor, location, semester, credits, color, description, created_at, updated_at
		FROM courses
		WHERE semester = ?
		ORDER BY name ASC
	`
	rows, err := r.db.Query(query, semester)
	if err != nil {
		return nil, fmt.Errorf("failed to get courses by semester: %w", err)
	}
	defer rows.Close()

	var courses []models.Course
	for rows.Next() {
		var course models.Course
		err := rows.Scan(
			&course.ID,
			&course.Name,
			&course.Code,
			&course.Professor,
			&course.Location,
			&course.Semester,
			&course.Credits,
			&course.Color,
			&course.Description,
			&course.CreatedAt,
			&course.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan course: %w", err)
		}

		// Load schedules
		schedules, err := r.GetSchedules(course.ID)
		if err != nil {
			return nil, err
		}
		course.Schedule = schedules

		courses = append(courses, course)
	}

	return courses, nil
}

// CreateSchedule creates a course schedule
func (r *CourseRepository) CreateSchedule(schedule *models.CourseSchedule) error {
	query := `
		INSERT INTO course_schedules (id, course_id, day_of_week, start_time, end_time, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query,
		schedule.ID,
		schedule.CourseID,
		schedule.DayOfWeek,
		schedule.StartTime,
		schedule.EndTime,
		schedule.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create course schedule: %w", err)
	}
	return nil
}

// GetSchedules retrieves all schedules for a course
func (r *CourseRepository) GetSchedules(courseID string) ([]models.CourseSchedule, error) {
	query := `
		SELECT id, course_id, day_of_week, start_time, end_time, created_at
		FROM course_schedules
		WHERE course_id = ?
		ORDER BY day_of_week, start_time
	`
	rows, err := r.db.Query(query, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get course schedules: %w", err)
	}
	defer rows.Close()

	var schedules []models.CourseSchedule
	for rows.Next() {
		var schedule models.CourseSchedule
		err := rows.Scan(
			&schedule.ID,
			&schedule.CourseID,
			&schedule.DayOfWeek,
			&schedule.StartTime,
			&schedule.EndTime,
			&schedule.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan course schedule: %w", err)
		}
		schedules = append(schedules, schedule)
	}

	return schedules, nil
}

// DeleteSchedules deletes all schedules for a course
func (r *CourseRepository) DeleteSchedules(courseID string) error {
	query := "DELETE FROM course_schedules WHERE course_id = ?"
	_, err := r.db.Exec(query, courseID)
	if err != nil {
		return fmt.Errorf("failed to delete course schedules: %w", err)
	}
	return nil
}

// CreateNote creates a course note
func (r *CourseRepository) CreateNote(note *models.CourseNote) error {
	tagsJSON, err := json.Marshal(note.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	query := `
		INSERT INTO course_notes (id, course_id, title, content, date, tags, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = r.db.Exec(query,
		note.ID,
		note.CourseID,
		note.Title,
		note.Content,
		note.Date,
		string(tagsJSON),
		note.CreatedAt,
		note.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create course note: %w", err)
	}
	return nil
}

// GetNotes retrieves all notes for a course
func (r *CourseRepository) GetNotes(courseID string) ([]models.CourseNote, error) {
	query := `
		SELECT id, course_id, title, content, date, tags, created_at, updated_at
		FROM course_notes
		WHERE course_id = ?
		ORDER BY date DESC
	`
	rows, err := r.db.Query(query, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get course notes: %w", err)
	}
	defer rows.Close()

	var notes []models.CourseNote
	for rows.Next() {
		var note models.CourseNote
		var tagsJSON string
		err := rows.Scan(
			&note.ID,
			&note.CourseID,
			&note.Title,
			&note.Content,
			&note.Date,
			&tagsJSON,
			&note.CreatedAt,
			&note.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan course note: %w", err)
		}

		if tagsJSON != "" {
			if err := json.Unmarshal([]byte(tagsJSON), &note.Tags); err != nil {
				return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
			}
		}

		notes = append(notes, note)
	}

	return notes, nil
}

// CreateAttendance creates an attendance record
func (r *CourseRepository) CreateAttendance(attendance *models.CourseAttendance) error {
	query := `
		INSERT INTO course_attendance (id, course_id, date, status, notes, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query,
		attendance.ID,
		attendance.CourseID,
		attendance.Date,
		attendance.Status,
		attendance.Notes,
		attendance.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create attendance: %w", err)
	}
	return nil
}

// GetAttendance retrieves all attendance records for a course
func (r *CourseRepository) GetAttendance(courseID string) ([]models.CourseAttendance, error) {
	query := `
		SELECT id, course_id, date, status, notes, created_at
		FROM course_attendance
		WHERE course_id = ?
		ORDER BY date DESC
	`
	rows, err := r.db.Query(query, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get attendance: %w", err)
	}
	defer rows.Close()

	var records []models.CourseAttendance
	for rows.Next() {
		var record models.CourseAttendance
		err := rows.Scan(
			&record.ID,
			&record.CourseID,
			&record.Date,
			&record.Status,
			&record.Notes,
			&record.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan attendance: %w", err)
		}
		records = append(records, record)
	}

	return records, nil
}
