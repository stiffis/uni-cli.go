# Nerd Fonts Icons Reference

## Icons Used in UniCLI

### App Title
-  - Graduation cap (Education/Student theme)

### Sidebar Navigation
-  - Tasks
-  - Calendar
- 󱉟 - Classes/Books
-  - Grades/Stars
-  - Notes
- 󰄨 - Statistics/Bar chart
-  - Settings/Cog

### Task Priorities (Kanban & Form)
-  - Urgent (exclamation circle) - Red
-  - High (arrow up) - Yellow/Orange
-  - Medium (minus/dash) - Blue
-  - Low (arrow down) - Gray

### Task Status/Dates
-  - Overdue warning (exclamation triangle)
-  - Due today (calendar)

### Form Title
-  - Tasks icon in "New Task" modal

## Color Mapping

### Priority Colors:
- Urgent:  → Red (#BF616A)
- High:  → Yellow/Orange (#EBCB8B)
- Medium:  → Blue (#81A1C1)
- Low:  → Gray (#4C566A)

### Date Status Colors:
- Overdue:  → Red (#BF616A)
- Today:  → Yellow (#EBCB8B)

## Complete Icon List

### Used in App (app.go):
```
  - graduation_cap (title)
  - tasks (sidebar)
  - calendar (sidebar)
󱉟 - book alt (sidebar)
  - star (sidebar)
  - sticky note (sidebar)
󰄨 - bar chart (sidebar)
  - cog (sidebar)
```

### Used in Kanban (tasks.go):
```
  - exclamation circle (urgent)
  - arrow up (high)
  - minus (medium)
  - arrow down (low)
  - exclamation triangle (overdue)
  - calendar (due today)
```

### Used in Form (taskform.go):
```
  - tasks (form title)
  - exclamation circle (urgent priority)
  - arrow up (high priority)
  - minus (medium priority)
  - arrow down (low priority)
```

## Font Requirements

**Required:** Nerd Fonts installed
- Recommended: `FiraCode Nerd Font`, `JetBrainsMono Nerd Font`, `Hack Nerd Font`
- Download: https://www.nerdfonts.com/

## Terminal Setup

Make sure your terminal is using a Nerd Font:
```bash
# Check if icons display correctly
echo "  󱉟 󰄨 "
```

If you see squares or missing characters, install Nerd Fonts.

## Testing All Icons

Run this to test all icon displays:
```bash
echo "App Title:  "
echo "Sidebar:  󱉟 󰄨 "
echo "Priority:    "
echo "Status:  "
```

All icons should display clearly without squares or missing glyphs.

## Icon Sources

All icons are from Nerd Fonts collection:
- Font Awesome icons (fa-*)
- Material Design Icons (md-*)
- Custom Nerd Fonts icons

Total unique icons used: **13 icons**
