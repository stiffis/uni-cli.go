# 🏗️ Arquitectura de UniCLI

## Visión General

UniCLI está construido con **Go** y **Bubble Tea**, siguiendo el patrón **Model-View-Update (MVU)** inspirado en la arquitectura Elm. La aplicación está diseñada para ser modular, testeable y escalable.

## Stack Tecnológico

### Core
- **Go 1.21+**: Lenguaje principal
- **Bubble Tea**: Framework TUI (Model-View-Update)
- **Lipgloss**: Sistema de estilos y layout
- **Bubbles**: Componentes TUI reutilizables

### Base de Datos
- **modernc.org/sqlite**: SQLite puro en Go (sin CGo)
- Schema relacional normalizado
- Migraciones automáticas

### Utilidades
- **google/uuid**: Generación de IDs únicos
- **Go standard library**: Para manejo de tiempo, archivos, etc.

## Estructura del Proyecto

```
UniCLI/
├── cmd/
│   └── unicli/
│       └── main.go              # Entry point de la aplicación
│
├── internal/                     # Código privado de la aplicación
│   ├── app/
│   │   └── app.go               # Modelo principal de Bubble Tea
│   │
│   ├── config/
│   │   └── config.go            # Configuración y temas
│   │
│   ├── database/
│   │   ├── database.go          # Conexión y migraciones
│   │   ├── repositories/        # Data access layer
│   │   │   ├── base.go          # Repository base
│   │   │   ├── task_repo.go     # CRUD de tareas
│   │   │   ├── class_repo.go    # CRUD de clases
│   │   │   └── ...
│   │   └── migrations/          # SQL migrations
│   │
│   ├── models/                   # Modelos de dominio
│   │   ├── task.go              # Modelo Task
│   │   ├── class.go             # Modelo Class
│   │   ├── grade.go             # Modelo Grade
│   │   ├── event.go             # Modelo Event
│   │   └── note.go              # Modelo Note
│   │
│   ├── services/                 # Lógica de negocio
│   │   ├── task_service.go      # Operaciones de tareas
│   │   ├── calendar_service.go  # Lógica de calendario
│   │   ├── stats_service.go     # Cálculos y estadísticas
│   │   └── ...
│   │
│   └── ui/                       # Interfaz de usuario
│       ├── screens/              # Pantallas principales
│       │   ├── tasks.go         # Vista de tareas
│       │   ├── calendar.go      # Vista de calendario
│       │   ├── classes.go       # Vista de clases
│       │   ├── grades.go        # Vista de calificaciones
│       │   ├── notes.go         # Vista de notas
│       │   └── settings.go      # Vista de configuración
│       │
│       ├── components/           # Componentes reutilizables
│       │   ├── form.go          # Formularios
│       │   ├── modal.go         # Modales/diálogos
│       │   ├── list.go          # Listas personalizadas
│       │   ├── calendar.go      # Widget de calendario
│       │   └── chart.go         # Gráficos ASCII
│       │
│       └── styles/
│           └── styles.go        # Definición de estilos y colores
│
├── pkg/                          # Código reutilizable (público)
│   └── utils/
│       ├── date.go              # Utilidades de fechas
│       ├── format.go            # Formateo de texto
│       └── export.go            # Exportación de datos
│
├── data/                         # Datos del usuario (gitignored)
│   ├── unicli.db                # Base de datos SQLite
│   └── config.json              # Configuración del usuario
│
├── go.mod                        # Dependencias
├── go.sum
├── README.md
├── PLANNING.md
├── ARCHITECTURE.md               # Este archivo
└── .gitignore
```

## Flujo de Datos (MVU Pattern)

```
┌─────────────────────────────────────────────────┐
│                    Main                         │
│  1. Carga config                                │
│  2. Inicializa DB                               │
│  3. Crea Model                                  │
│  4. Ejecuta tea.Program                         │
└──────────────────┬──────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────┐
│              App Model (MVU)                    │
│                                                 │
│  Init()                                         │
│    └─> Comando inicial (si necesario)          │
│                                                 │
│  Update(msg)                                    │
│    ├─> Procesa mensaje                         │
│    ├─> Actualiza estado                        │
│    └─> Retorna nuevo modelo + comando          │
│                                                 │
│  View()                                         │
│    └─> Renderiza UI actual                     │
│                                                 │
└─────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────┐
│                  Messages                       │
│                                                 │
│  • tea.KeyMsg        - Teclas presionadas       │
│  • tea.WindowSizeMsg - Resize de terminal      │
│  • tea.MouseMsg      - Eventos de mouse        │
│  • CustomMsg         - Mensajes personalizados │
│                                                 │
└─────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────┐
│                  Commands                       │
│                                                 │
│  • Database queries (async)                     │
│  • File operations                              │
│  • Timer/Tick events                            │
│  • Custom async operations                      │
│                                                 │
└─────────────────────────────────────────────────┘
```

## Capas de la Aplicación

### 1. Capa de Presentación (UI Layer)

**Responsabilidad**: Renderizar la interfaz y manejar input del usuario.

- **Screens**: Vistas completas (tasks, calendar, etc.)
- **Components**: Widgets reutilizables (forms, modals, lists)
- **Styles**: Definiciones de colores y estilos

**Principios:**
- Las screens NO acceden directamente a la DB
- Toda la lógica de negocio va en Services
- Usan Bubble Tea messages para comunicación async

### 2. Capa de Servicio (Business Logic)

**Responsabilidad**: Implementar reglas de negocio y orquestar operaciones.

**Ejemplos:**
- `TaskService`: Validar fechas, filtrar tareas, calcular estadísticas
- `GradeService`: Calcular promedios ponderados, GPAs
- `CalendarService`: Detectar conflictos de horario, generar vistas

**Principios:**
- Sin dependencias de UI
- Testeable independientemente
- Orquesta llamadas a múltiples repositories

### 3. Capa de Datos (Data Layer)

**Responsabilidad**: Acceso a datos y persistencia.

**Repositories**: Abstracción sobre la DB
- CRUD operations
- Queries específicas
- Transacciones

**Principios:**
- Un repository por entidad principal
- Retorna modelos de dominio, no structs SQL
- Maneja errores de DB apropiadamente

### 4. Modelos de Dominio

**Responsabilidad**: Representar entidades del dominio.

**Características:**
- Structs con lógica de negocio simple
- Métodos de validación
- Métodos helper (IsOverdue, IsDueToday, etc.)
- Sin dependencias de DB o UI

## Schema de Base de Datos

### Tablas Principales

```sql
tasks
├── id (TEXT PRIMARY KEY)
├── title (TEXT NOT NULL)
├── description (TEXT)
├── status (TEXT)
├── priority (TEXT)
├── category (TEXT)
├── due_date (DATETIME)
├── created_at (DATETIME)
├── updated_at (DATETIME)
└── completed_at (DATETIME)

classes
├── id (TEXT PRIMARY KEY)
├── name (TEXT NOT NULL)
├── professor (TEXT)
├── room (TEXT)
├── color (TEXT)
├── semester (TEXT)
├── credits (INTEGER)
└── created_at (DATETIME)

schedules
├── id (INTEGER PRIMARY KEY)
├── class_id (TEXT FK → classes)
├── day_of_week (INTEGER)
├── start_time (TEXT)
└── end_time (TEXT)

grades
├── id (TEXT PRIMARY KEY)
├── class_id (TEXT FK → classes)
├── name (TEXT NOT NULL)
├── score (REAL)
├── max_score (REAL)
├── weight (REAL)
├── date (DATE)
├── type (TEXT)
└── created_at (DATETIME)

events
├── id (TEXT PRIMARY KEY)
├── title (TEXT NOT NULL)
├── description (TEXT)
├── start_datetime (DATETIME)
├── end_datetime (DATETIME)
├── type (TEXT)
└── created_at (DATETIME)

notes
├── id (TEXT PRIMARY KEY)
├── title (TEXT NOT NULL)
├── content (TEXT)
├── created_at (DATETIME)
└── updated_at (DATETIME)

tags
├── id (INTEGER PRIMARY KEY)
└── name (TEXT UNIQUE)

task_tags
├── task_id (TEXT FK → tasks)
└── tag_id (INTEGER FK → tags)
```

## Diseño de UI (Estilo Lazygit)

### Layout Principal

```
┌────────────────────────────────────────────────────────┐
│ 🎓 UniCLI - Student Organization Manager       [?][q]  │ <- Title Bar
├───────────┬────────────────────────────────────────────┤
│           │                                            │
│ 📋 Tasks  │  ┌─ Today's Tasks ──────────────────────┐ │
│ 📅 Calendar│  │ ○ 🔴 Study for exam         (Today) │ │
│ 🎒 Classes│  │ ✓ 🟡 Submit homework      (Complete) │ │
│ 📊 Grades │  │ ○ 🔵 Read chapter              (...)  │ │
│ 📝 Notes  │  └─────────────────────────────────────┘ │
│ 📈 Stats  │                                            │
│ ⚙️ Settings  ┌─ Upcoming ───────────────────────────┐ │
│           │  │ 📅 Midterm exam - 3 days            │ │
│           │  │ 📝 Project due - 1 week             │ │
│           │  └─────────────────────────────────────┘ │
│           │                                            │
└───────────┴────────────────────────────────────────────┘
│ [n] new [e] edit [d] delete [space] toggle    j/k nav  │ <- Status Bar
└────────────────────────────────────────────────────────┘
```

### Navegación

#### Global
- `1-7` o `t,c,s,g,n,⚙`: Cambiar de vista
- `Tab` / `Shift+Tab`: Cambiar entre paneles
- `?`: Mostrar ayuda
- `q` / `Ctrl+C`: Salir

#### Listas
- `j` / `k` o `↓` / `↑`: Navegar
- `g` / `G`: Ir al inicio/final
- `Enter`: Seleccionar/abrir
- `Space`: Toggle/seleccionar

#### Acciones
- `n`: Nuevo item
- `e`: Editar item
- `d`: Eliminar item
- `/`: Buscar/filtrar
- `r`: Refrescar

## Manejo de Estado

### Estado Global (App Model)
```go
type Model struct {
    db          *database.DB
    cfg         *config.Config
    currentView View
    width       int
    height      int
    
    // Screens
    taskScreen    tea.Model
    calendarScreen tea.Model
    // ... otros screens
    
    ready bool
    err   error
}
```

### Estado de Pantalla (Task Screen Example)
```go
type TaskScreen struct {
    db       *database.DB
    tasks    []models.Task
    cursor   int
    selected map[int]struct{}
    filter   string
    
    width  int
    height int
}
```

## Patrones de Diseño Utilizados

### 1. Repository Pattern
Abstrae el acceso a datos, facilitando testing y cambios de DB.

```go
type TaskRepository interface {
    Create(task *models.Task) error
    FindByID(id string) (*models.Task, error)
    FindAll() ([]models.Task, error)
    Update(task *models.Task) error
    Delete(id string) error
}
```

### 2. Service Layer Pattern
Encapsula lógica de negocio compleja.

```go
type TaskService struct {
    repo TaskRepository
}

func (s *TaskService) GetTasksDueToday() ([]models.Task, error) {
    // Lógica de filtrado compleja
}
```

### 3. Model-View-Update (MVU/Elm Architecture)
Patrón central de Bubble Tea:
- **Model**: Estado inmutable
- **Update**: Función pura que transforma estado
- **View**: Función pura que renderiza estado

### 4. Command Pattern
Para operaciones asíncronas (DB queries, timers, etc.)

```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    return m, m.loadTasksCmd()
}

func (m Model) loadTasksCmd() tea.Cmd {
    return func() tea.Msg {
        tasks, err := m.db.Tasks().FindAll()
        return tasksLoadedMsg{tasks: tasks, err: err}
    }
}
```

## Extensibilidad

### Agregar una Nueva Vista

1. Crear screen en `internal/ui/screens/nueva_vista.go`
2. Implementar interfaz `tea.Model` (Init, Update, View)
3. Agregar View enum en `internal/app/app.go`
4. Agregar case en Update y View de App
5. Agregar en sidebar

### Agregar una Nueva Entidad

1. Modelo en `internal/models/entidad.go`
2. Tabla en `internal/database/database.go` (migrations)
3. Repository en `internal/database/repositories/entidad_repo.go`
4. Service en `internal/services/entidad_service.go`
5. UI en `internal/ui/screens/entidad.go`

## Testing Strategy

### Unit Tests
- Modelos: Lógica de validación
- Services: Reglas de negocio
- Repositories: CRUD operations (con DB en memoria)

### Integration Tests
- Flujos completos (crear tarea → guardar → recuperar)
- Migraciones de DB

### UI Tests (futuros)
- Bubble Tea tiene soporte para testing de TUI
- Golden tests para snapshots de UI

## Performance Considerations

### Base de Datos
- Índices en columnas frecuentemente consultadas
- Prepared statements para queries repetidas
- Transacciones para operaciones múltiples

### UI
- Lazy loading para listas grandes
- Virtualización de listas (solo renderizar visible)
- Debouncing para búsqueda/filtros

### Memoria
- Paginación para datasets grandes
- Cleanup de recursos al cambiar de vista

## Roadmap Técnico

### Fase 1: Foundation ✅
- [x] Estructura del proyecto
- [x] Schema de DB
- [x] Modelos básicos
- [x] App skeleton con Bubble Tea
- [x] Vista de tareas básica

### Fase 2: Core Features
- [ ] Implementar repositories
- [ ] CRUD completo de tareas
- [ ] Formularios y modales
- [ ] Vista de calendario
- [ ] Vista de clases

### Fase 3: Advanced Features
- [ ] Estadísticas y gráficos
- [ ] Búsqueda y filtros avanzados
- [ ] Timer Pomodoro
- [ ] Sistema de tags mejorado

### Fase 4: Polish
- [ ] Temas personalizables
- [ ] Exportación de datos
- [ ] Importación desde otras apps
- [ ] Documentación completa
- [ ] Tests comprehensivos

### Fase 5: Distribution
- [ ] CI/CD pipeline
- [ ] Cross-compilation
- [ ] Package managers (Homebrew, AUR, etc.)
- [ ] Auto-updates

## Referencias

- [Bubble Tea Documentation](https://github.com/charmbracelet/bubbletea)
- [Lipgloss Documentation](https://github.com/charmbracelet/lipgloss)
- [The Elm Architecture](https://guide.elm-lang.org/architecture/)
- [lazygit](https://github.com/jesseduffield/lazygit) - Inspiración UI

---

*Documento en evolución - actualizar conforme el proyecto crece*
