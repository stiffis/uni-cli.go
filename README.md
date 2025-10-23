# 🎓 UniCLI - Student Organization TUI - In Development

A modern, beautiful Terminal User Interface (TUI) application for student organization and productivity. Built with Go and Bubble Tea, inspired by lazygit's intuitive design.

## ✨ Features

### Core Functionality

- 📋 **Task Management**: Create, edit, complete, and organize tasks with priorities and deadlines
- 📅 **Calendar View**: Monthly and weekly calendar with all your events and deadlines
- 🎒 **Class Schedule**: Weekly timetable with class information and locations
- 📊 **Grade Tracking**: Record grades and automatically calculate averages
- 📝 **Quick Notes**: Markdown-based notes with tags
- ⏱️ **Pomodoro Timer**: Built-in focus timer
- 📈 **Statistics**: Visual insights into your productivity

### UI Features

- Beautiful multi-panel interface inspired by lazygit
- Intuitive keyboard navigation
- Customizable color themes
- Real-time updates
- Context-sensitive help

## 🚀 Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/unicli.git
cd unicli

# Build
go build -o unicli ./cmd/unicli

# Run
./unicli
```

### Requirements

- Go 1.21 or higher
- A terminal with true color support (recommended)

## ⌨️ Keyboard Shortcuts

| Key                    | Action                    |
| ---------------------- | ------------------------- |
| `Tab` / `Shift+Tab`    | Navigate between panels   |
| `j` / `k` or `↓` / `↑` | Navigate up/down in lists |
| `Enter`                | Select/Open item          |
| `n`                    | New item                  |
| `e`                    | Edit item                 |
| `d`                    | Delete item               |
| `Space`                | Toggle complete           |
| `/`                    | Search/Filter             |
| `?`                    | Show help                 |
| `q` or `Ctrl+C`        | Quit                      |

## 📦 Project Structure

```
unicli/
├── cmd/unicli/          # Main application entry point
├── internal/
│   ├── app/            # Application core
│   ├── ui/             # UI components and screens
│   ├── models/         # Data models
│   ├── services/       # Business logic
│   ├── database/       # Database layer
│   └── config/         # Configuration
├── pkg/                # Reusable packages
└── data/               # User data (gitignored)
```

## 󰚓 ScreenShots

![UniCLI Interface](assets/welcomeview.png)
![Taks Management](assets/tasksview.png)

## 🛠️ Development

```bash
# Run in development mode
go run ./cmd/unicli

# Run tests
go test ./...

# Build for production
go build -ldflags="-s -w" -o unicli ./cmd/unicli
```

## 📝 License

MIT License - see LICENSE file for details

## 🙏 Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Style definitions
- [lazygit](https://github.com/jesseduffield/lazygit) - UI inspiration
