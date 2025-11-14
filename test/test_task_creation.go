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
	fmt.Println("🧪 Testing task creation...")
	fmt.Println()

	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("❌ Error loading config: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ Config loaded: %s\n", cfg.DatabasePath)

	db, err := database.New(cfg.DatabasePath)
	if err != nil {
		fmt.Printf("❌ Error initializing database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()
	fmt.Println("✓ Database connected")

	if err := db.Migrate(); err != nil {
		fmt.Printf("❌ Error running migrations: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Migrations completed")
	fmt.Println()

	// Count existing tasks
	existingTasks, err := db.Tasks().FindAll()
	if err != nil {
		fmt.Printf("❌ Error counting tasks: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("📊 Existing tasks: %d\n", len(existingTasks))
	fmt.Println()

	fmt.Println("🔨 Creating test task...")
	testTask := models.NewTask("Test Task from Script")
	testTask.Description = "This is a test task to verify creation works"
	testTask.Priority = models.TaskPriorityMedium
	testTask.Status = models.TaskStatusPending
	testTask.Tags = []string{"test", "debug"}
	testTask.UpdatedAt = time.Now()

	// Try to create the task
	err = db.Tasks().Create(testTask)
	if err != nil {
		fmt.Printf("❌ Error creating task: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ Task created with ID: %s\n", testTask.ID)
	fmt.Println()

	// Verify it was saved
	savedTask, err := db.Tasks().FindByID(testTask.ID)
	if err != nil {
		fmt.Printf("❌ Error retrieving task: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Task retrieved from database:")
	fmt.Printf("  - ID: %s\n", savedTask.ID)
	fmt.Printf("  - Title: %s\n", savedTask.Title)
	fmt.Printf("  - Description: %s\n", savedTask.Description)
	fmt.Printf("  - Status: %s\n", savedTask.Status)
	fmt.Printf("  - Priority: %s\n", savedTask.Priority)
	fmt.Printf("  - Tags: %v\n", savedTask.Tags)
	fmt.Printf("  - Created: %s\n", savedTask.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Println()

	// Count tasks again
	allTasks, err := db.Tasks().FindAll()
	if err != nil {
		fmt.Printf("❌ Error counting tasks: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("📊 Total tasks now: %d\n", len(allTasks))
	fmt.Println()
	fmt.Println("✅ All tests passed! Task creation is working correctly.")
}
