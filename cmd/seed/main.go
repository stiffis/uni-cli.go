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
		createTask("Study for Calculus exam", "Review chapters 5-7 and practice problems", models.TaskPriorityHigh, models.TaskStatusPending, 0, []string{"study", "math"}),
		createTask("Complete project proposal", "Write and submit the final project proposal for CS class", models.TaskPriorityUrgent, models.TaskStatusPending, 0, []string{"project", "cs"}),
		createTask("Read Chapter 5", "Read chapter 5 of Operating Systems textbook", models.TaskPriorityMedium, models.TaskStatusPending, 3, []string{"reading"}),
		createTask("Group meeting preparation", "Prepare slides for tomorrow's group meeting", models.TaskPriorityMedium, models.TaskStatusPending, 1, []string{"meeting", "project"}),
		createTask("Submit homework", "Submit physics homework before deadline", models.TaskPriorityHigh, models.TaskStatusPending, -1, []string{"homework", "physics"}),
		createTask("Grocery shopping", "Buy ingredients for the week", models.TaskPriorityLow, models.TaskStatusPending, 2, []string{"personal"}),
	}

	fmt.Println("Creating sample tasks...")
	for _, task := range tasks {
		if err := db.Tasks().Create(task); err != nil {
			fmt.Printf("Error creating task '%s': %v\n", task.Title, err)
		} else {
			fmt.Printf("✓ Created: %s\n", task.Title)
		}
	}

	fmt.Println("\nSample tasks created successfully!")
	fmt.Println("Run './unicli' to see them in the app.")
}

func createTask(title, description string, priority models.TaskPriority, status models.TaskStatus, dueDays int, tags []string) *models.Task {
	task := models.NewTask(title)
	task.Description = description
	task.Priority = priority
	task.Status = status
	task.Tags = tags

	if dueDays != 0 {
		dueDate := time.Now().AddDate(0, 0, dueDays)
		task.DueDate = &dueDate
	}

	return task
}
