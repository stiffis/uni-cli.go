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
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	db, err := database.New(cfg.DatabasePath)
	if err != nil {
		fmt.Printf("Error initializing database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Migrate(); err != nil {
		fmt.Printf("Error running migrations: %v\n", err)
		os.Exit(1)
	}

	tasks := []*models.Task{
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

	categories := createCategories(db)
	createEvents(db, categories)
	createCourses(db, categories)

	fmt.Println("\nSample tasks, categories, events and courses created successfully!")
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
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	
	events := []*models.Event{
		{
			Title:         "Morning Standup",
			Description:   "Daily team sync - Discuss progress and blockers",
			StartDatetime: today.Add(9 * time.Hour),
			EndDatetime:   timePtr(today.Add(9*time.Hour + 30*time.Minute)),
			Type:          "event",
			CategoryID:    getCategoryIDByName("Work", categories),
		},
		{
			Title:         "Deep Work Session",
			Description:   "Focus time - Work on backend API authentication module",
			StartDatetime: today.Add(10 * time.Hour),
			EndDatetime:   timePtr(today.Add(13 * time.Hour)),
			Type:          "event",
			CategoryID:    getCategoryIDByName("Projects", categories),
		},
		{
			Title:         "Lunch Break",
			Description:   "Team lunch at downtown cafe",
			StartDatetime: today.Add(13 * time.Hour),
			EndDatetime:   timePtr(today.Add(14 * time.Hour)),
			Type:          "event",
			CategoryID:    getCategoryIDByName("Personal", categories),
		},
		{
			Title:         "Code Review",
			Description:   "Review pull requests from team members",
			StartDatetime: today.Add(14*time.Hour + 30*time.Minute),
			EndDatetime:   timePtr(today.Add(15*time.Hour + 30*time.Minute)),
			Type:          "event",
			CategoryID:    getCategoryIDByName("Work", categories),
		},
		{
			Title:         "Gym Workout",
			Description:   "Leg day - Squats, lunges, and cardio",
			StartDatetime: today.Add(18 * time.Hour),
			EndDatetime:   timePtr(today.Add(19*time.Hour + 30*time.Minute)),
			Type:          "event",
			CategoryID:    getCategoryIDByName("Health", categories),
		},
		{
			Title:         "Calculus Lecture",
			Description:   "Chapter 7: Integration techniques - Room 301",
			StartDatetime: today.AddDate(0, 0, 1).Add(8 * time.Hour),
			EndDatetime:   timePtr(today.AddDate(0, 0, 1).Add(10 * time.Hour)),
			Type:          "event",
			CategoryID:    getCategoryIDByName("University", categories),
		},
		{
			Title:         "Study Group",
			Description:   "Physics study session - Library 2nd floor",
			StartDatetime: today.AddDate(0, 0, 1).Add(14 * time.Hour),
			EndDatetime:   timePtr(today.AddDate(0, 0, 1).Add(16 * time.Hour)),
			Type:          "event",
			CategoryID:    getCategoryIDByName("University", categories),
		},
		{
			Title:         "Client Presentation",
			Description:   "Q4 project demo - Prepare slides and demo environment",
			StartDatetime: today.AddDate(0, 0, 1).Add(16*time.Hour + 30*time.Minute),
			EndDatetime:   timePtr(today.AddDate(0, 0, 1).Add(17*time.Hour + 30*time.Minute)),
			Type:          "event",
			CategoryID:    getCategoryIDByName("Work", categories),
		},
		{
			Title:         "Doctor's Appointment",
			Description:   "Annual check-up with Dr. Smith - Don't forget insurance card",
			StartDatetime: today.AddDate(0, 0, 2).Add(10*time.Hour + 30*time.Minute),
			EndDatetime:   timePtr(today.AddDate(0, 0, 2).Add(11*time.Hour + 30*time.Minute)),
			Type:          "event",
			CategoryID:    getCategoryIDByName("Health", categories),
		},
		{
			Title:         "Job Interview",
			Description:   "Software Engineer position at TechCorp - Virtual interview via Zoom",
			StartDatetime: today.AddDate(0, 0, 2).Add(15 * time.Hour),
			EndDatetime:   timePtr(today.AddDate(0, 0, 2).Add(16*time.Hour + 30*time.Minute)),
			Type:          "event",
			CategoryID:    getCategoryIDByName("Career", categories),
		},
		{
			Title:         "Chemistry Lab",
			Description:   "Experiment: Acid-base titration - Lab coat required",
			StartDatetime: today.AddDate(0, 0, 3).Add(13 * time.Hour),
			EndDatetime:   timePtr(today.AddDate(0, 0, 3).Add(16 * time.Hour)),
			Type:          "event",
			CategoryID:    getCategoryIDByName("University", categories),
		},
		{
			Title:         "Dinner with Friends",
			Description:   "Celebrating Sarah's promotion at Italian restaurant",
			StartDatetime: today.AddDate(0, 0, 3).Add(19*time.Hour + 30*time.Minute),
			EndDatetime:   timePtr(today.AddDate(0, 0, 3).Add(22 * time.Hour)),
			Type:          "event",
			CategoryID:    getCategoryIDByName("Social", categories),
		},
		{
			Title:         "Dentist Appointment",
			Description:   "Regular cleaning and checkup - Downtown Dental Clinic",
			StartDatetime: today.AddDate(0, 0, 4).Add(9 * time.Hour),
			EndDatetime:   timePtr(today.AddDate(0, 0, 4).Add(10 * time.Hour)),
			Type:          "event",
			CategoryID:    getCategoryIDByName("Health", categories),
		},
		{
			Title:         "Team Building Event",
			Description:   "Escape room activity with work team",
			StartDatetime: today.AddDate(0, 0, 4).Add(14 * time.Hour),
			EndDatetime:   timePtr(today.AddDate(0, 0, 4).Add(17 * time.Hour)),
			Type:          "event",
			CategoryID:    getCategoryIDByName("Work", categories),
		},
		{
			Title:         "Birthday Party",
			Description:   "John's 25th birthday party - Bring gift!",
			StartDatetime: today.AddDate(0, 0, 5).Add(18 * time.Hour),
			EndDatetime:   timePtr(today.AddDate(0, 0, 5).Add(23 * time.Hour)),
			Type:          "event",
			CategoryID:    getCategoryIDByName("Social", categories),
		},
		{
			Title:         "Project Deadline",
			Description:   "Final submission for CS project - GitHub repository due",
			StartDatetime: today.AddDate(0, 0, 6).Add(23*time.Hour + 59*time.Minute),
			EndDatetime:   timePtr(today.AddDate(0, 0, 6).Add(23*time.Hour + 59*time.Minute)),
			Type:          "event",
			CategoryID:    getCategoryIDByName("Projects", categories),
		},
		{
			Title:         "Guitar Lesson",
			Description:   "Weekly lesson with instructor - Practice sheet music from last week",
			StartDatetime: today.AddDate(0, 0, 7).Add(16 * time.Hour),
			EndDatetime:   timePtr(today.AddDate(0, 0, 7).Add(17 * time.Hour)),
			Type:          "event",
			CategoryID:    getCategoryIDByName("Hobbies", categories),
		},
		{
			Title:         "Budget Review",
			Description:   "Monthly budget review and planning for next month",
			StartDatetime: today.AddDate(0, 0, 7).Add(20 * time.Hour),
			EndDatetime:   timePtr(today.AddDate(0, 0, 7).Add(21 * time.Hour)),
			Type:          "event",
			CategoryID:    getCategoryIDByName("Finance", categories),
		},
		{
			Title:         "Monday Morning Planning",
			Description:   "Weekly planning session - Set goals for the week",
			StartDatetime: today.AddDate(0, 0, 7).Add(8 * time.Hour),
			EndDatetime:   timePtr(today.AddDate(0, 0, 7).Add(9 * time.Hour)),
			Type:          "event",
			RecurrenceRule: "weekly",
			CategoryID:    getCategoryIDByName("Personal", categories),
		},
	}

	fmt.Println("\nCreating sample events...")
	for _, event := range events {
		event.ID = models.NewEvent(event.Title, event.StartDatetime).ID
		event.CreatedAt = time.Now()
		if err := db.Events().Create(event); err != nil {
			fmt.Printf("Error creating event '%s': %v\n", event.Title, err)
		} else {
			fmt.Printf("✓ Created: %s (at %s)\n", event.Title, event.StartDatetime.Format("Mon Jan 02 15:04"))
		}
	}
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func getCategoryIDByName(name string, categories []models.Category) string {
	for _, category := range categories {
		if category.Name == name {
			return category.ID
		}
	}
	return ""
}

func createCourses(db *database.DB, categories []models.Category) {
	now := time.Now()
	year := now.Year()
	
	courses := []struct {
		name        string
		code        string
		professor   string
		location    string
		semester    string
		credits     int
		description string
		color       string
		schedules   []struct {
			day       int    // 1=Mon, 2=Tue, 3=Wed, 4=Thu, 5=Fri
			startTime string
			endTime   string
		}
	}{
		{
			name:        "Data Structures & Algorithms",
			code:        "CS-201",
			professor:   "Dr. Sarah Chen",
			location:    "Building A, Room 301",
			semester:    fmt.Sprintf("Fall %d", year),
			credits:     4,
			description: "Study of fundamental data structures and algorithms with analysis of their time and space complexity.",
			color:       "#C34043",
			schedules: []struct {
				day       int
				startTime string
				endTime   string
			}{
				{1, "09:00", "10:30"},
				{3, "09:00", "10:30"},
				{5, "09:00", "10:30"},
			},
		},
		{
			name:        "Operating Systems",
			code:        "CS-301",
			professor:   "Prof. Michael Rodriguez",
			location:    "Building B, Room 205",
			semester:    fmt.Sprintf("Fall %d", year),
			credits:     4,
			description: "Comprehensive study of operating system concepts including process management, memory management, and file systems.",
			color:       "#76946A",
			schedules: []struct {
				day       int
				startTime string
				endTime   string
			}{
				{2, "14:00", "15:30"},
				{4, "14:00", "15:30"},
			},
		},
		{
			name:        "Database Systems",
			code:        "CS-350",
			professor:   "Dr. Emily Zhang",
			location:    "Tech Center, Lab 4",
			semester:    fmt.Sprintf("Fall %d", year),
			credits:     3,
			description: "Design and implementation of database systems including SQL, normalization, and transaction processing.",
			color:       "#C8C093",
			schedules: []struct {
				day       int
				startTime string
				endTime   string
			}{
				{1, "13:00", "14:30"},
				{3, "13:00", "14:30"},
			},
		},
		{
			name:        "Software Engineering",
			code:        "CS-320",
			professor:   "Prof. James Anderson",
			location:    "Engineering Building, Room 102",
			semester:    fmt.Sprintf("Fall %d", year),
			credits:     3,
			description: "Principles and practices of software development including design patterns, testing, and agile methodologies.",
			color:       "#957FB8",
			schedules: []struct {
				day       int
				startTime string
				endTime   string
			}{
				{2, "10:00", "11:30"},
				{4, "10:00", "11:30"},
			},
		},
		{
			name:        "Computer Networks",
			code:        "CS-410",
			professor:   "Dr. Lisa Thompson",
			location:    "Building C, Room 450",
			semester:    fmt.Sprintf("Fall %d", year),
			credits:     3,
			description: "Study of computer network architectures, protocols, and applications including TCP/IP and network security.",
			color:       "#7AA89F",
			schedules: []struct {
				day       int
				startTime string
				endTime   string
			}{
				{1, "15:00", "16:30"},
				{3, "15:00", "16:30"},
			},
		},
	}

	fmt.Println("\nCreating sample courses...")
	for _, c := range courses {
		course := models.NewCourse(c.name)
		course.Code = c.code
		course.Professor = c.professor
		course.Location = c.location
		course.Semester = c.semester
		course.Credits = c.credits
		course.Description = c.description
		course.Color = c.color

		for _, sched := range c.schedules {
			schedule := models.NewCourseSchedule(course.ID, sched.day, sched.startTime, sched.endTime)
			course.Schedule = append(course.Schedule, *schedule)
		}

		if err := db.Courses().Create(course); err != nil {
			fmt.Printf("Error creating course '%s': %v\n", course.Name, err)
		} else {
			fmt.Printf("✓ Created: %s (%s)\n", course.Name, course.Code)
			for _, sched := range course.Schedule {
				dayName := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}[sched.DayOfWeek%7]
				fmt.Printf("  → %s %s-%s\n", dayName, sched.StartTime, sched.EndTime)
			}
		}
	}
}
