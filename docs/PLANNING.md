# UniCLI - Planificación de Aplicación TUI para Organización Estudiantil

## 📋 Índice
1. [Visión General](#visión-general)
2. [Casos de Uso](#casos-de-uso)
3. [Tecnologías y Stack](#tecnologías-y-stack)
4. [Arquitectura](#arquitectura)
5. [Funcionalidades Principales](#funcionalidades-principales)
6. [Modelo de Datos](#modelo-de-datos)
7. [Estructura del Proyecto](#estructura-del-proyecto)
8. [Roadmap de Desarrollo](#roadmap-de-desarrollo)

---

## 🎯 Visión General

### Objetivo
Crear una aplicación TUI (Text User Interface) moderna y eficiente para la gestión de actividades estudiantiles, tareas, horarios y organización académica desde la terminal.

### Filosofía del Proyecto
- **Minimalista**: Interfaz limpia y enfocada en la productividad
- **Rápida**: Navegación mediante teclado, sin necesidad de mouse
- **Offline-first**: Funcionamiento sin internet, con sincronización opcional
- **Extensible**: Arquitectura modular para agregar funcionalidades
- **Cross-platform**: Compatible con Linux, macOS y Windows

---

## 💡 Casos de Uso

### Usuarios Objetivo
- Estudiantes universitarios que prefieren trabajar en la terminal
- Personas que buscan organización sin distracciones
- Usuarios de Vim/Emacs y entornos minimalistas
- Estudiantes de informática/ingeniería

### Escenarios de Uso
1. **Gestión de Tareas**: Crear, editar, completar tareas con deadlines
2. **Horarios de Clases**: Visualizar horarios semanales
3. **Seguimiento de Notas**: Registrar calificaciones y calcular promedios
4. **Calendario Académico**: Fechas importantes, exámenes, entregas
5. **Notas Rápidas**: Apuntes y recordatorios
6. **Pomodoro/Timer**: Técnica de estudio
7. **Gestión de Proyectos**: Proyectos grupales con subtareas

---

## 🛠 Tecnologías y Stack

### Opciones de Lenguajes y Frameworks

#### Opción 1: Python + Rich/Textual
**Pros:**
- Textual es un framework TUI moderno y reactivo
- Rich para renderizado avanzado de texto
- Desarrollo rápido y gran ecosistema
- Fácil manejo de datos (JSON, SQLite)
- Async/await nativo

**Contras:**
- Rendimiento moderado en operaciones intensivas
- Requiere Python instalado

**Librerías:**
```
- textual: Framework TUI principal
- rich: Renderizado de texto enriquecido
- pydantic: Validación de datos
- sqlalchemy o tinydb: Base de datos
- python-dateutil: Manejo de fechas
- click: CLI arguments
```

#### Opción 2: Rust + Ratatui
**Pros:**
- Alto rendimiento y bajo consumo de memoria
- Binario compilado, no requiere runtime
- Type safety fuerte
- Ratatui es maduro y estable

**Contras:**
- Curva de aprendizaje más pronunciada
- Desarrollo más lento inicialmente
- Manejo de errores más verbose

**Librerías:**
```
- ratatui: Framework TUI
- crossterm: Input handling cross-platform
- serde: Serialización de datos
- tokio: Runtime async
- sqlx o sled: Base de datos
```

#### Opción 3: Go + Bubble Tea
**Pros:**
- Excelente balance entre facilidad y rendimiento
- Bubble Tea es muy popular y bien mantenido
- Compilación rápida, binario único
- Goroutines para concurrencia

**Contras:**
- Ecosistema TUI más pequeño que Python
- Error handling con if err != nil

**Librerías:**
```
- bubbletea: Framework TUI (Elm Architecture)
- lipgloss: Styling
- bubbles: Componentes pre-hechos
- charm/log: Logging
- SQLite3/BoltDB: Base de datos
```

#### Opción 4: Node.js + Ink (React)
**Pros:**
- React para TUI (componentes familiares)
- Gran ecosistema npm
- TypeScript support
- Desarrollo web-like

**Contras:**
- Mayor consumo de memoria
- Requiere Node.js instalado
- Rendimiento inferior

**Librerías:**
```
- ink: React para terminal
- ink-ui: Componentes adicionales
- lowdb/sqlite3: Base de datos
- date-fns: Manejo de fechas
```

### ✅ Decisión Tomada
**Go + Bubble Tea** - Excelente balance entre rendimiento y facilidad de desarrollo, con una comunidad activa y componentes bien mantenidos. Estilo visual inspirado en lazygit con múltiples paneles y navegación intuitiva.

---

## 🏗 Arquitectura

### Patrón de Diseño Sugerido: MVC/MVU

```
┌─────────────────────────────────────────┐
│            Presentación (TUI)            │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐ │
│  │ Screens │  │ Widgets │  │  Input  │ │
│  └─────────┘  └─────────┘  └─────────┘ │
└──────────────────┬──────────────────────┘
                   │
┌──────────────────▼──────────────────────┐
│          Lógica de Negocio              │
│  ┌──────────┐  ┌──────────┐            │
│  │ Services │  │ Managers │            │
│  └──────────┘  └──────────┘            │
└──────────────────┬──────────────────────┘
                   │
┌──────────────────▼──────────────────────┐
│          Capa de Datos                  │
│  ┌──────────┐  ┌──────────┐            │
│  │   DAO    │  │  Models  │            │
│  └──────────┘  └──────────┘            │
└──────────────────┬──────────────────────┘
                   │
           ┌───────▼────────┐
           │   SQLite/JSON  │
           └────────────────┘
```

### Componentes Principales

1. **App Controller**: Punto de entrada, manejo de navegación
2. **Screen Manager**: Gestión de pantallas y transiciones
3. **Data Layer**: Persistencia y modelos de datos
4. **Service Layer**: Lógica de negocio (cálculos, validaciones)
5. **Config Manager**: Configuración y preferencias del usuario

---

## ✨ Funcionalidades Principales

### MVP (Minimum Viable Product)
1. ✅ Gestión básica de tareas (CRUD)
2. ✅ Visualización de lista de tareas
3. ✅ Marcar tareas como completadas
4. ✅ Fechas de vencimiento
5. ✅ Categorías/Tags
6. ✅ Persistencia local

### Fase 2
- Horario de clases semanal
- Vista de calendario mensual
- Búsqueda y filtros avanzados
- Estadísticas y reportes
- Temas de color personalizables

### Fase 3
- Gestión de notas/calificaciones
- Timer Pomodoro integrado
- Exportar datos (JSON, Markdown, PDF)
- Recordatorios/Notificaciones
- Sincronización opcional (Git/Cloud)

### Fase 4
- Modo colaborativo (proyectos grupales)
- Integración con calendarios externos (iCal)
- Plugin system
- Comandos personalizados
- CLI commands para scripting

---

## 📊 Modelo de Datos

### Entidades Principales

```yaml
Task:
  id: UUID
  title: String
  description: String (opcional)
  status: Enum[pending, in_progress, completed, cancelled]
  priority: Enum[low, medium, high, urgent]
  category: String (opcional)
  tags: List[String]
  due_date: DateTime (opcional)
  created_at: DateTime
  updated_at: DateTime
  completed_at: DateTime (opcional)
  estimated_time: Integer (minutos)
  actual_time: Integer (minutos)

Class:
  id: UUID
  name: String
  professor: String
  room: String
  schedule: List[Schedule]
  color: String (hex)
  semester: String
  credits: Integer

Schedule:
  day_of_week: Enum[monday...sunday]
  start_time: Time
  end_time: Time

Grade:
  id: UUID
  class_id: FK
  name: String
  score: Float
  max_score: Float
  weight: Float (porcentaje)
  date: Date
  type: Enum[exam, homework, project, quiz, participation]

Event:
  id: UUID
  title: String
  description: String
  start_datetime: DateTime
  end_datetime: DateTime
  type: Enum[exam, deadline, meeting, personal]
  related_to: FK (opcional, Task/Class)
  reminder: Boolean
  reminder_time: Integer (minutos antes)

Note:
  id: UUID
  title: String
  content: String (Markdown)
  tags: List[String]
  created_at: DateTime
  updated_at: DateTime

Config:
  theme: String
  default_view: String
  pomodoro_duration: Integer
  break_duration: Integer
  notifications_enabled: Boolean
  first_day_of_week: Enum[monday, sunday]
```

---

## 📁 Estructura del Proyecto

### Opción Python + Textual

```
unicli/
├── README.md
├── PLANNING.md
├── requirements.txt
├── setup.py
├── pyproject.toml
├── .gitignore
├── docs/
│   ├── user_guide.md
│   └── developer_guide.md
├── tests/
│   ├── __init__.py
│   ├── test_tasks.py
│   ├── test_ui.py
│   └── fixtures/
├── src/
│   └── unicli/
│       ├── __init__.py
│       ├── __main__.py
│       ├── app.py                 # Aplicación principal
│       ├── config.py              # Configuración
│       ├── cli.py                 # CLI entry point
│       │
│       ├── ui/                    # Interfaz TUI
│       │   ├── __init__.py
│       │   ├── screens/           # Pantallas principales
│       │   │   ├── __init__.py
│       │   │   ├── main.py
│       │   │   ├── tasks.py
│       │   │   ├── calendar.py
│       │   │   ├── classes.py
│       │   │   └── settings.py
│       │   ├── widgets/           # Componentes reutilizables
│       │   │   ├── __init__.py
│       │   │   ├── task_list.py
│       │   │   ├── calendar_view.py
│       │   │   ├── form.py
│       │   │   └── modal.py
│       │   └── theme.py           # Estilos y colores
│       │
│       ├── models/                # Modelos de datos
│       │   ├── __init__.py
│       │   ├── task.py
│       │   ├── class_model.py
│       │   ├── grade.py
│       │   ├── event.py
│       │   └── note.py
│       │
│       ├── services/              # Lógica de negocio
│       │   ├── __init__.py
│       │   ├── task_service.py
│       │   ├── class_service.py
│       │   ├── grade_service.py
│       │   ├── calendar_service.py
│       │   └── stats_service.py
│       │
│       ├── database/              # Capa de datos
│       │   ├── __init__.py
│       │   ├── connection.py
│       │   ├── repositories/
│       │   │   ├── __init__.py
│       │   │   ├── base.py
│       │   │   ├── task_repo.py
│       │   │   ├── class_repo.py
│       │   │   └── event_repo.py
│       │   └── migrations/
│       │       └── initial.sql
│       │
│       └── utils/                 # Utilidades
│           ├── __init__.py
│           ├── date_utils.py
│           ├── validators.py
│           └── exporters.py
│
└── data/                          # Datos del usuario (gitignored)
    ├── unicli.db
    └── config.json
```

---

## 🗺 Roadmap de Desarrollo

### Sprint 1: Fundamentos (1-2 semanas)
- [ ] Setup del proyecto y estructura
- [ ] Configuración de entorno de desarrollo
- [ ] Modelo de datos básico (Task)
- [ ] Database layer con SQLite
- [ ] Tests unitarios básicos

### Sprint 2: UI Básica (1-2 semanas)
- [ ] Pantalla principal con navegación
- [ ] Lista de tareas (visualización)
- [ ] Formulario para crear/editar tareas
- [ ] Keybindings básicos (j/k, enter, q)
- [ ] Sistema de temas

### Sprint 3: Funcionalidad Core (2 semanas)
- [ ] CRUD completo de tareas
- [ ] Filtros y búsqueda
- [ ] Categorías y tags
- [ ] Fechas de vencimiento
- [ ] Persistencia funcional

### Sprint 4: Experiencia de Usuario (1-2 semanas)
- [ ] Mejoras en navegación
- [ ] Atajos de teclado avanzados
- [ ] Mensajes de confirmación
- [ ] Manejo de errores elegante
- [ ] Ayuda contextual (?)

### Sprint 5: Calendario y Clases (2 semanas)
- [ ] Vista de calendario mensual
- [ ] Gestión de horario de clases
- [ ] Vista semanal
- [ ] Integración con tareas

### Sprint 6: Estadísticas y Extras (1-2 semanas)
- [ ] Dashboard con estadísticas
- [ ] Gestión de notas/calificaciones
- [ ] Exportación de datos
- [ ] Timer Pomodoro

---

## ✅ Decisiones Tomadas

### Confirmadas:
1. **Lenguaje/Framework**: ✅ **Go + Bubble Tea**
2. **Base de datos**: ✅ **SQLite** (con go-sqlite3 o modernc.org/sqlite)
3. **Alcance**: ✅ **App completa** con todas las funcionalidades planificadas
4. **Estilo de UI**: ✅ **Visual estilo Lazygit** - paneles, colores, navegación fluida
5. **Distribución**: Binario compilado (cross-platform)
6. **Sincronización**: Local primero, opcional para futuro

### Próximas Decisiones:
- Estructura de paneles específica
- Esquema de colores por defecto
- Comandos Git-style vs shortcuts tradicionales

---

## 📚 Referencias e Inspiración

### TUI Apps Similares:
- **taskwarrior-tui**: Gestión de tareas en Rust
- **lazygit**: UI intuitiva con paneles
- **k9s**: Navegación y shortcuts bien diseñados
- **glow**: Renderizado de Markdown
- **btop**: Widgets y visualización de datos

### Frameworks TUI:
- Textual (Python): https://textual.textualize.io/
- Ratatui (Rust): https://ratatui.rs/
- Bubble Tea (Go): https://github.com/charmbracelet/bubbletea
- Ink (Node): https://github.com/vadimdemedes/ink

---

## 🎨 Conceptos de UI

### Navegación Propuesta:
```
Ctrl+P: Command Palette
Tab: Cambiar entre paneles
j/k: Navegar arriba/abajo
h/l: Navegar izquierda/derecha
Enter: Seleccionar/Editar
Esc: Cancelar/Volver
/: Buscar
?: Ayuda
q: Salir
```

### Layout Principal:
```
┌─────────────────────────────────────────────────┐
│ UniCLI - Student Organizer     [?] Help [q] Quit│
├──────────┬──────────────────────────────────────┤
│          │                                      │
│  [T]asks │  ┌─ Today's Tasks ─────────────────┐│
│  [C]alendar  │ □ Estudiar para examen          ││
│  [S]chedule  │ ☑ Entregar proyecto             ││
│  [G]rades    │ □ Leer capítulo 5               ││
│  [N]otes     │                                 ││
│  ───────     │                                 ││
│  [St]ats     └─────────────────────────────────┘│
│  [Se]ttings                                     │
│          │  ┌─ Upcoming ─────────────────────┐ │
│          │  │ 📅 Examen Cálculo - 2 days     │ │
│          │  │ 📝 Proyecto Final - 1 week     │ │
│          │  └────────────────────────────────┘ │
│          │                                      │
└──────────┴──────────────────────────────────────┘
```

---

## ✅ Next Steps

1. **Revisar este documento** y agregar/modificar según tus necesidades
2. **Decidir stack tecnológico** basado en experiencia y objetivos
3. **Crear estructura inicial** del proyecto
4. **Setup de entorno** de desarrollo
5. **Comenzar con MVP**: Modelo de Task + UI básica

---

*Documento vivo - Actualizar conforme el proyecto evoluciona*
