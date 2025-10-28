# Quick Start Guide - UniCLI

## Overview

UniCLI is a modern Terminal User Interface (TUI) application for student organization and productivity, built with Go and Bubble Tea.

## Installation

### Requirements
- Go 1.21 or higher
- A terminal with true color support (recommended)

### Quick Install
```bash
# Clone the repository
git clone https://github.com/yourusername/UniCLI.git
cd UniCLI

# Build & Run
./run.sh
```

### Manual Build
```bash
# Build
go build -o unicli ./cmd/unicli

# Run
./unicli
```

## Current Features

### ✅ Task Management (Kanban Board)
- Create, edit, and delete tasks
- 3-column board: To Do → In Progress → Done
- Priority levels: Urgent, High, Medium, Low
- Due dates with overdue indicators
- Tags support
- Full persistence to SQLite database

### ✅ User Interface
- Beautiful multi-panel interface
- Sidebar navigation
- Command mode (`:s` sidebar, `:q` quit, `:h` help)
- Nord/Zen color theme
- Keyboard-driven navigation

## Keyboard Shortcuts

### Global Navigation
| Key | Action |
|-----|--------|
| `:s` | Open sidebar / view switcher |
| `:q` | Quit application |
| `:h` | Help (coming soon) |
| `Ctrl+C` | Force quit |

### Task View (Kanban)
| Key | Action |
|-----|--------|
| `n` | New task |
| `e` | Edit selected task |
| `d` / `Delete` | Delete selected task |
| `Enter` | Select/deselect task |
| `Tab` | Switch to next column |
| `Shift+Tab` | Switch to previous column |
| `j` / `↓` | Move down |
| `k` / `↑` | Move up |
| `←` | Move selected task to previous column |
| `→` | Move selected task to next column |
| `g` | Go to top |
| `G` | Go to bottom |
| `r` | Refresh |

### Form Navigation
| Key | Action |
|-----|--------|
| `Tab` | Next field |
| `Shift+Tab` | Previous field |
| `Enter` | Submit (when on button) |
| `Esc` | Cancel |
| `←` / `→` | Change priority |

## First Steps

### 1. Launch UniCLI
```bash
cd UniCLI
./unicli
```

### 2. Create Sample Tasks
```bash
# Run seed script to create example tasks
go run ./cmd/seed/main.go
```

### 3. Navigate to Tasks
- Press `:s` to open sidebar
- Select "Tasks" and press Enter
- You'll see a Kanban board with your tasks

### 4. Create Your First Task
1. Press `n` for new task
2. Fill in the form:
   - **Title**: e.g., "Study for Calculus exam"
   - **Description**: Optional details
   - **Due Date**: YYYY-MM-DD format (e.g., 2025-11-01)
   - **Priority**: Use `←` / `→` to change
3. Press `Tab` to navigate to Create button
4. Press `Enter` to submit

### 5. Manage Tasks
- Use `Enter` to select a task
- Press `→` to move it to "In Progress"
- Press `e` to edit
- Press `d` to delete
- Press `→` again to move to "Done"

## Data Storage

- **Database**: `~/.unicli/unicli.db` (SQLite)
- **Config**: `~/.unicli/config.json` (future)

## Coming Soon

- 📅 Calendar view with monthly/weekly display
- 🎒 Class schedule and timetable
- 📊 Grade tracking and GPA calculator
- 📝 Quick notes with markdown support
- ⏱️ Pomodoro timer
- 📈 Productivity statistics

## Troubleshooting

### Database Issues
```bash
# Check if database exists
ls -la ~/.unicli/unicli.db

# Reset database (⚠️ deletes all data)
rm ~/.unicli/unicli.db
go run ./cmd/seed/main.go
```

### Build Issues
```bash
# Clean and rebuild
rm unicli
go clean
go build -o unicli ./cmd/unicli
```

---

**Ready to be productive? Press `:s` and start organizing!** 🚀
