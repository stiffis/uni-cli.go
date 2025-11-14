package main

import (
	"fmt"
	"os"
	"time"

	"github.com/stiffis/UniCLI/internal/config"
	"github.com/stiffis/UniCLI/internal/database"
	"github.com/stiffis/UniCLI/internal/models"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Initialize database
	db, err := database.New(cfg.DatabasePath)
	if err != nil {
		fmt.Printf("Error initializing database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Run migrations
	if err := db.Migrate(); err != nil {
		fmt.Printf("Error running migrations: %v\n", err)
		os.Exit(1)
	}

	// Create sample tasks
	tasks := []*models.Task{
		createTask("Study for Calculus exam", "Review chapters 5-7 and practice problems", models.TaskPriorityHigh, models.TaskStatusPending, 0, []string{"study", "math"}, []string{"Review chapter 5", "Review chapter 6", "Practice problems"}),
		createTask("Complete project proposal", "Write and submit the final project proposal for CS class", models.TaskPriorityUrgent, models.TaskStatusPending, 0, []string{"project", "cs"}, []string{"Write introduction", "Gather resources", "Write conclusion"}),
		createTask("Read Chapter 5", "Read chapter 5 of Operating Systems textbook", models.TaskPriorityMedium, models.TaskStatusPending, 3, []string{"reading"}, nil),
		createTask("Group meeting preparation", "Prepare slides for tomorrow's group meeting", models.TaskPriorityMedium, models.TaskStatusPending, 1, []string{"meeting", "project"}, []string{"Create agenda", "Prepare slides"}),
		createTask("Submit homework", "Submit physics homework before deadline", models.TaskPriorityHigh, models.TaskStatusPending, -1, []string{"homework", "physics"}, nil),
		createTask("Grocery shopping", "Buy ingredients for the week", models.TaskPriorityLow, models.TaskStatusPending, 2, []string{"personal"}, []string{"Buy milk", "Buy eggs", "Buy bread"}),
	}

	fmt.Println("Creating sample tasks...")
	for _, task := range tasks {
		if err := db.Tasks().Create(task); err != nil {
			fmt.Printf("Error creating task '%s': %v\n", task.Title, err)
		} else {
			fmt.Printf("✓ Created: %s\n", task.Title)
			for _, subtask := range task.Subtasks {
				subtask.TaskID = task.ID
				if err := db.Tasks().CreateSubtask(&subtask); err != nil {
					fmt.Printf("  Error creating subtask '%s': %v\n", subtask.Title, err)
				} else {
					fmt.Printf("  ✓ Subtask Created: %s\n", subtask.Title)
				}
			}
		}
	}

	// Create sample categories
	categories := createCategories(db)

	// Create sample events
	createEvents(db, categories)

	fmt.Println("\nSample tasks, categories and events created successfully!")
	fmt.Println("Run './unicli' to see them in the app.")
}

func createTask(title, description string, priority models.TaskPriority, status models.TaskStatus, dueDays int, tags []string, subtaskTitles []string) *models.Task {
	task := models.NewTask(title)
	task.Description = description
	task.Priority = priority
	task.Status = status
	task.Tags = tags

	if dueDays != 0 {
		dueDate := time.Now().AddDate(0, 0, dueDays)
		task.DueDate = &dueDate
	}

	for _, subtaskTitle := range subtaskTitles {
		task.Subtasks = append(task.Subtasks, models.Subtask{Title: subtaskTitle})
	}

	return task
}

func createCategories(db *database.DB) []models.Category {
	categories := []models.Category{
		*models.NewCategory("Personal", "#FFC0CB"), // Pink
		*models.NewCategory("Work", "#ADD8E6"),     // Light Blue
		*models.NewCategory("University", "#90EE90"), // Light Green
		*models.NewCategory("Health", "#FFD700"),    // Gold
		*models.NewCategory("Social", "#FFA07A"),    // Light Salmon
	}

	fmt.Println("\nCreating sample categories...")
	for _, category := range categories {
		// In a real app, you'd use db.Categories().Create(category)
		// For now, we'll just insert it directly
		query := "INSERT INTO categories (id, name, color) VALUES (?, ?, ?)"
		if _, err := db.Conn().Exec(query, category.ID, category.Name, category.Color); err != nil {
			fmt.Printf("Error creating category '%s': %v\n", category.Name, err)
		} else {
			fmt.Printf("✓ Created: %s\n", category.Name)
		}
	}
	return categories
}

func createEvents(db *database.DB, categories []models.Category) {
	events := []*models.Event{
		{
			Title:         "Doctor's Appointment",
			Description:   "Annual check-up",
			StartDatetime: time.Now().AddDate(0, 0, 3),
			CategoryID:    getCategoryIDByName("Health", categories),
		},
		{
			Title:         "Team Meeting",
			Description:   "Weekly sync-up",
			StartDatetime: time.Now().AddDate(0, 0, 1),
			CategoryID:    getCategoryIDByName("Work", categories),
		},
		{
			Title:         "Lunch with Mom",
			Description:   "At the new Italian place",
			StartDatetime: time.Now().AddDate(0, 0, 5),
			CategoryID:    getCategoryIDByName("Personal", categories),
		},
		{
			Title:         "Finals Study Group",
			Description:   "Library, 2nd floor",
			StartDatetime: time.Now().AddDate(0, 0, 2),
			CategoryID:    getCategoryIDByName("University", categories),
		},
	}

	fmt.Println("\nCreating sample events...")
	for _, event := range events {
		event.ID = models.NewEvent(event.Title, event.StartDatetime).ID
		event.CreatedAt = time.Now()
		if err := db.Events().Create(event); err != nil {
			fmt.Printf("Error creating event '%s': %v\n", event.Title, err)
		} else {
			fmt.Printf("✓ Created: %s\n", event.Title)
		}
	}
}

func getCategoryIDByName(name string, categories []models.Category) string {
	for _, category := range categories {
		if category.Name == name {
			return category.ID
		}
	}
	return ""
}
