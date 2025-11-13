package database

import (
	"database/sql"
	"fmt"

	"github.com/stiffis/UniCLI/internal/database/repositories"
	_ "modernc.org/sqlite"
)

// DB wraps the database connection
type DB struct {
	conn      *sql.DB
	taskRepo  *repositories.TaskRepository
	eventRepo *repositories.EventRepository
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

	CREATE TABLE IF NOT EXISTS classes (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		professor TEXT,
		room TEXT,
		color TEXT,
		semester TEXT,
		credits INTEGER,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS schedules (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		class_id TEXT NOT NULL,
		day_of_week INTEGER NOT NULL,
		start_time TEXT NOT NULL,
		end_time TEXT NOT NULL,
		FOREIGN KEY (class_id) REFERENCES classes(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS grades (
		id TEXT PRIMARY KEY,
		class_id TEXT NOT NULL,
		name TEXT NOT NULL,
		score REAL NOT NULL,
		max_score REAL NOT NULL,
		weight REAL DEFAULT 1.0,
		date DATE,
		type TEXT,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (class_id) REFERENCES classes(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS events (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT,
		start_datetime DATETIME NOT NULL,
		end_datetime DATETIME,
		type TEXT,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
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
	`

	if _, err := db.conn.Exec(schema); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
