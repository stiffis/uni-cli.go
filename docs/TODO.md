# Roadmap - UniCLI

## ✅ Completado

### Core Infrastructure
- [x] SQLite database with full schema
- [x] Repository pattern implementation
- [x] Task CRUD operations (10+ methods)
- [x] Base UI framework with Bubble Tea
- [x] Styling system (Nord/Zen theme)

### Task Management
- [x] Kanban board with 3 columns (To Do, In Progress, Done)
- [x] Create tasks with form (title, description, priority, due date)
- [x] Edit existing tasks
- [x] Delete tasks
- [x] Move tasks between columns
- [x] Task persistence to database
- [x] Tags support
- [x] Priority indicators (urgent, high, medium, low)
- [x] Due date indicators (overdue, today, upcoming)

### UI Components
- [x] Input field component
- [x] Textarea component
- [x] TaskForm (create/edit with validation)
- [x] Sidebar navigation
- [x] Command mode (:s, :q, :h)

## 🚧 En Progreso

### Bug Fixes & Improvements
- [ ] Remove debug logging from production
- [ ] Add keyboard shortcuts help overlay
- [ ] Improve error messages

## 📋 Próximas Features

### Calendar View (Prioridad Alta)
- [ ] Monthly calendar view
- [ ] Display events and deadlines
- [ ] Navigate between months
- [ ] Add events from calendar
- [ ] Integration with tasks due dates

### Classes & Schedule (Prioridad Alta)
- [ ] Weekly timetable view
- [ ] Class CRUD operations
- [ ] ClassRepository implementation
- [ ] Schedule conflict detection
- [ ] Color-coded classes

### Search & Filters (Prioridad Media)
- [ ] Fuzzy search in tasks
- [ ] Filter by status/priority/category/tags
- [ ] Filter by date range
- [ ] Save favorite filters

### Grades Management (Prioridad Media)
- [ ] Grades screen
- [ ] GradeRepository implementation
- [ ] CRUD operations for grades
- [ ] Calculate averages per class
- [ ] Calculate overall GPA
- [ ] Grade evolution chart

### Notes (Prioridad Baja)
- [ ] Notes screen with markdown support
- [ ] NoteRepository implementation
- [ ] Markdown preview
- [ ] Tags for notes
- [ ] Search in notes

### Statistics Dashboard (Prioridad Baja)
- [ ] Stats screen
- [ ] Tasks completed this week/month
- [ ] Productivity chart by day
- [ ] Priority distribution graph
- [ ] Upcoming deadlines summary

## 🎨 Polish & Features

### User Experience
- [ ] Pomodoro timer
- [ ] System notifications
- [ ] Themes (Dracula, Nord, Gruvbox, Solarized)
- [ ] Custom keyboard shortcuts
- [ ] Confirmation dialogs

### Data Management
- [ ] Export to JSON/Markdown/CSV
- [ ] Import from Todoist/Notion
- [ ] Automatic backups

### Development
- [ ] Unit tests
- [ ] Integration tests
- [ ] CI/CD pipeline
- [ ] Cross-platform builds (Linux, macOS, Windows)

---

**Last updated**: 2025-10-28  
**Status**: ~50% complete (Core task management fully functional)
