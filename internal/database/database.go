package database

import (
	"database/sql"
	"fmt"

	"github.com/stiffis/UniCLI/internal/database/repositories"
	_ "modernc.org/sqlite"
)

// DB wraps the database connection
type DB struct {
	conn         *sql.DB
	taskRepo     *repositories.TaskRepository
	eventRepo    *repositories.EventRepository
	categoryRepo *repositories.CategoryRepository
	courseRepo   *repositories.CourseRepository
}

// New creates a new database connection
func New(path string) (*DB, error) {
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable foreign keys
	if _, err := conn.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	db := &DB{conn: conn}

	// Initialize repositories
	db.taskRepo = repositories.NewTaskRepository(conn)
	db.eventRepo = repositories.NewEventRepository(conn)
	db.categoryRepo = repositories.NewCategoryRepository(conn)
	db.courseRepo = repositories.NewCourseRepository(conn)

	return db, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// Conn returns the underlying database connection
func (db *DB) Conn() *sql.DB {
	return db.conn
}

// Tasks returns the task repository
func (db *DB) Tasks() *repositories.TaskRepository {
	return db.taskRepo
}

// Events returns the event repository
func (db *DB) Events() *repositories.EventRepository {
	return db.eventRepo
}

// Categories returns the category repository
func (db *DB) Categories() *repositories.CategoryRepository {
	return db.categoryRepo
}

// Courses returns the course repository
func (db *DB) Courses() *repositories.CourseRepository {
	return db.courseRepo
}

// Migrate runs database migrations
func (db *DB) Migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS tasks (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT,
		status TEXT NOT NULL DEFAULT 'pending',
		priority TEXT NOT NULL DEFAULT 'medium',
		category TEXT,
		due_date DATETIME,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		completed_at DATETIME
	);

	CREATE TABLE IF NOT EXISTS tags (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL
	);

	CREATE TABLE IF NOT EXISTS task_tags (
		task_id TEXT NOT NULL,
		tag_id INTEGER NOT NULL,
		FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
		FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE,
		PRIMARY KEY (task_id, tag_id)
	);

	CREATE TABLE IF NOT EXISTS subtasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		task_id TEXT NOT NULL,
		title TEXT NOT NULL,
		is_completed BOOLEAN NOT NULL DEFAULT 0,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS courses (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		code TEXT,
		professor TEXT,
		location TEXT,
		semester TEXT,
		credits INTEGER DEFAULT 0,
		color TEXT,
		description TEXT,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS course_schedules (
		id TEXT PRIMARY KEY,
		course_id TEXT NOT NULL,
		day_of_week INTEGER NOT NULL,
		start_time TEXT NOT NULL,
		end_time TEXT NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS course_notes (
		id TEXT PRIMARY KEY,
		course_id TEXT NOT NULL,
		title TEXT NOT NULL,
		content TEXT,
		date DATETIME,
		tags TEXT,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS course_attendance (
		id TEXT PRIMARY KEY,
		course_id TEXT NOT NULL,
		date DATETIME NOT NULL,
		status TEXT NOT NULL,
		notes TEXT,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS grades (
		id TEXT PRIMARY KEY,
		course_id TEXT NOT NULL,
		name TEXT NOT NULL,
		score REAL NOT NULL,
		max_score REAL NOT NULL,
		weight REAL DEFAULT 1.0,
		date DATE,
		type TEXT,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS events (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT,
		start_datetime DATETIME NOT NULL,
		end_datetime DATETIME,
		type TEXT,
		category_id TEXT,
		recurrence_rule TEXT,
		recurrence_end_date DATETIME,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
	);

	CREATE TABLE IF NOT EXISTS categories (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		color TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS notes (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		content TEXT,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
	CREATE INDEX IF NOT EXISTS idx_tasks_due_date ON tasks(due_date);
	CREATE INDEX IF NOT EXISTS idx_subtasks_task_id ON subtasks(task_id);
	CREATE INDEX IF NOT EXISTS idx_events_start ON events(start_datetime);
	CREATE INDEX IF NOT EXISTS idx_course_schedules_course_id ON course_schedules(course_id);
	CREATE INDEX IF NOT EXISTS idx_course_notes_course_id ON course_notes(course_id);
	CREATE INDEX IF NOT EXISTS idx_course_attendance_course_id ON course_attendance(course_id);
	`

	if _, err := db.conn.Exec(schema); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	if err := db.addColumnIfNotExists("events", "recurrence_rule", "TEXT"); err != nil {
		return err
	}
	if err := db.addColumnIfNotExists("events", "recurrence_end_date", "DATETIME"); err != nil {
		return err
	}
	if err := db.addColumnIfNotExists("events", "category_id", "TEXT"); err != nil {
		return err
	}

	return nil
}

func (db *DB) addColumnIfNotExists(tableName, columnName, columnType string) error {
	rows, err := db.conn.Query(fmt.Sprintf("PRAGMA table_info(%s)", tableName))
	if err != nil {
		return fmt.Errorf("failed to query table info: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var name string
		var type_ string
		var notnull int
		var dflt_value any
		var pk int
		if err := rows.Scan(&cid, &name, &type_, &notnull, &dflt_value, &pk); err != nil {
			return fmt.Errorf("failed to scan table info: %w", err)
		}
		if name == columnName {
			// Column already exists
			return nil
		}
	}

	query := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", tableName, columnName, columnType)
	if _, err := db.conn.Exec(query); err != nil {
		return fmt.Errorf("failed to add column %s to table %s: %w", columnName, tableName, err)
	}

	return nil
}
