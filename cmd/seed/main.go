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
		createTask("Study for Calculus exam", "Review chapters 5-7 and practice problems", models.TaskPriorityHigh, models.TaskStatusPending, 1, []string{"study", "math"}, []string{"Review chapter 5", "Review chapter 6", "Practice problems"}),
		createTask("Complete project proposal", "Write and submit the final project proposal for CS class", models.TaskPriorityUrgent, models.TaskStatusPending, 2, []string{"project", "cs"}, []string{"Write introduction", "Gather resources", "Write conclusion"}),
		createTask("Read Chapter 5", "Read chapter 5 of Operating Systems textbook", models.TaskPriorityMedium, models.TaskStatusPending, 3, []string{"reading"}, nil),
		createTask("Group meeting preparation", "Prepare slides for tomorrow's group meeting", models.TaskPriorityMedium, models.TaskStatusPending, 1, []string{"meeting", "project"}, []string{"Create agenda", "Prepare slides"}),
		createTask("Submit physics homework", "Submit physics homework before deadline", models.TaskPriorityHigh, models.TaskStatusPending, 4, []string{"homework", "physics"}, nil),
		createTask("Grocery shopping", "Buy ingredients for the week", models.TaskPriorityLow, models.TaskStatusPending, 2, []string{"personal"}, []string{"Buy milk", "Buy eggs", "Buy bread"}),
		createTask("Plan weekend trip", "Research and plan a weekend trip", models.TaskPriorityLow, models.TaskStatusPending, 5, []string{"personal", "travel"}, []string{"Choose destination", "Book accommodation", "Plan activities"}),
		createTask("Finish lab report", "Complete and submit the chemistry lab report", models.TaskPriorityHigh, models.TaskStatusPending, 3, []string{"lab", "chemistry"}, []string{"Analyze data", "Write report", "Proofread"}),
		createTask("Update resume", "Update resume with recent experience", models.TaskPriorityMedium, models.TaskStatusPending, 6, []string{"career"}, nil),
		createTask("Clean the apartment", "Clean the kitchen, bathroom, and living room", models.TaskPriorityLow, models.TaskStatusPending, 5, []string{"personal", "chores"}, []string{"Clean kitchen", "Clean bathroom", "Clean living room"}),
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
	// Using Kanagawa syntax highlighting colors (like in Neovim)
	categories := []models.Category{
		*models.NewCategory("Personal", "#D27E99"),      // sakuraPink (constants, keywords)
		*models.NewCategory("Work", "#7E9CD8"),          // springBlue (functions, methods)
		*models.NewCategory("University", "#76946A"),    // autumnGreen (success, diff add)
		*models.NewCategory("Health", "#7AA89F"),        // waveAqua2 (cyan bright)
		*models.NewCategory("Social", "#C34043"),        // autumnRed (errors, diff delete)
		*models.NewCategory("Finance", "#C8C093"),       // oldWhite (foreground alt)
		*models.NewCategory("Hobbies", "#957FB8"),       // oniViolet (keywords, macros)
		*models.NewCategory("Travel", "#7FB4CA"),        // lightBlue (identifiers)
		*models.NewCategory("Career", "#FF9E3B"),        // roninYellow (warnings, important)
		*models.NewCategory("Projects", "#957FB8"),      // oniViolet (keywords, macros)
	}

	fmt.Println("\nCreating sample categories...")
	for _, category := range categories {
		// In a real app, you'd use db.Categories().Create(category)
		// For now, we'll just insert it directly
		query := "INSERT OR IGNORE INTO categories (id, name, color) VALUES (?, ?, ?)"
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
		{
			Title:         "Dentist Appointment",
			Description:   "Regular cleaning",
			StartDatetime: time.Now().AddDate(0, 0, 4),
			CategoryID:    getCategoryIDByName("Health", categories),
		},
		{
			Title:         "Project Deadline",
			Description:   "Submit the final version",
			StartDatetime: time.Now().AddDate(0, 0, 6),
			CategoryID:    getCategoryIDByName("Projects", categories),
		},
		{
			Title:         "Birthday Party",
			Description:   "John's birthday party",
			StartDatetime: time.Now().AddDate(0, 0, 5),
			CategoryID:    getCategoryIDByName("Social", categories),
		},
		{
			Title:         "Job Interview",
			Description:   "Interview for the software engineer position",
			StartDatetime: time.Now().AddDate(0, 0, 3),
			CategoryID:    getCategoryIDByName("Career", categories),
		},
		{
			Title:         "Pay bills",
			Description:   "Pay electricity and internet bills",
			StartDatetime: time.Now().AddDate(0, 0, 1),
			CategoryID:    getCategoryIDByName("Finance", categories),
		},
		{
			Title:         "Go to the gym",
			Description:   "Workout session",
			StartDatetime: time.Now().AddDate(0, 0, 2),
			CategoryID:    getCategoryIDByName("Health", categories),
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
