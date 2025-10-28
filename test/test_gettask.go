package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/stiffis/UniCLI/internal/models"
)

// Simulate the GetTask method from TaskForm
func simulateGetTask(taskID string, titleValue string, descValue string, priority models.TaskPriority) *models.Task {
	var task *models.Task
	if taskID != "" {
		// Editing existing task
		task = &models.Task{
			ID: taskID,
		}
		fmt.Printf("  Editing existing task: %s\n", taskID)
	} else {
		// Creating new task
		task = models.NewTask("") // Title will be set below
		fmt.Printf("  Creating NEW task with ID: %s\n", task.ID)
	}

	task.Title = titleValue
	task.Description = descValue
	task.Priority = priority
	if taskID != "" {
		task.Status = models.TaskStatusInProgress // Just for testing
	} else {
		task.Status = models.TaskStatusPending // Always create as pending
	}
	task.UpdatedAt = time.Now()

	return task
}

func main() {
	fmt.Println("🧪 Testing TaskForm.GetTask() logic...")
	fmt.Println()

	// Test 1: Creating a new task (taskID is empty)
	fmt.Println("Test 1: Create new task")
	newTask := simulateGetTask("", "My New Task", "This is a test", models.TaskPriorityHigh)
	fmt.Printf("  Result:\n")
	fmt.Printf("    ID: %s\n", newTask.ID)
	fmt.Printf("    Title: %s\n", newTask.Title)
	fmt.Printf("    Description: %s\n", newTask.Description)
	fmt.Printf("    Status: %s\n", newTask.Status)
	fmt.Printf("    Priority: %s\n", newTask.Priority)
	fmt.Printf("    ID is empty: %v\n", newTask.ID == "")
	fmt.Printf("    Title is empty: %v\n", strings.TrimSpace(newTask.Title) == "")
	fmt.Println()

	// Test 2: Editing existing task
	fmt.Println("Test 2: Edit existing task")
	existingTask := simulateGetTask("existing-123", "Updated Title", "Updated desc", models.TaskPriorityLow)
	fmt.Printf("  Result:\n")
	fmt.Printf("    ID: %s\n", existingTask.ID)
	fmt.Printf("    Title: %s\n", existingTask.Title)
	fmt.Printf("    Status: %s\n", existingTask.Status)
	fmt.Println()

	// Test 3: Creating task with empty title (should still have ID)
	fmt.Println("Test 3: Create task with empty title")
	emptyTask := simulateGetTask("", "", "", models.TaskPriorityMedium)
	fmt.Printf("  Result:\n")
	fmt.Printf("    ID: %s\n", emptyTask.ID)
	fmt.Printf("    Title: '%s'\n", emptyTask.Title)
	fmt.Printf("    ID is empty: %v\n", emptyTask.ID == "")
	fmt.Println()

	fmt.Println("✅ Test complete")
}
