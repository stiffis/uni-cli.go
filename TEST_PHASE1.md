# 🚀 Testing FASE 1 - Quick Guide

## What just got implemented ✅

We now have **full database persistence** with a working TaskRepository!

## Quick Test

### 1. Add sample tasks:
```bash
go run ./cmd/seed/main.go
```

### 2. Run the app:
```bash
./unicli
```

### 3. Try these actions:

**Navigation:**
- `j` / `k` or arrows - Move between tasks
- `g` - Go to top
- `G` - Go to bottom

**Actions that NOW WORK:**
- `space` - Toggle task completion (✓ / ☐)
- `d` - Delete current task
- `r` - Refresh tasks from database

**Navigation:**
- `:s` - Open sidebar to switch views
- `:q` - Quit

## What to expect

You should see:
- 6 sample tasks loaded from real database
- Tasks with different priorities (colors)
- Some tasks marked as "due today" or "overdue"
- Tags displayed on tasks (study, math, project, etc.)

## Test the new features

1. **Toggle completion:** 
   - Navigate to a task and press `space`
   - Task should show ✓ and strikethrough
   - Press `space` again to mark as incomplete

2. **Delete a task:**
   - Navigate to a task and press `d`
   - Task should disappear immediately

3. **Persistence test:**
   - Toggle some tasks, delete others
   - Close the app (`:q`)
   - Open again (`./unicli`)
   - Your changes should be saved!

## If something goes wrong

### Reset everything:
```bash
rm ~/.unicli/unicli.db
go run ./cmd/seed/main.go
./unicli
```

### Recompile:
```bash
go build -o unicli ./cmd/unicli
```

## Next: FASE 2

Creating forms to add/edit tasks! 🎯
