# 🚀 Guía de Inicio Rápido - UniCLI

## ✅ Estado Actual del Proyecto

### Lo que está funcionando:
- ✅ Estructura completa del proyecto
- ✅ Base de datos SQLite con schema completo
- ✅ Modelos de dominio (Task, Class)
- ✅ Sistema de estilos inspirado en lazygit
- ✅ Aplicación principal con navegación entre vistas
- ✅ Vista de tareas con datos de ejemplo
- ✅ Sidebar de navegación estilo lazygit
- ✅ Sistema de colores y temas
- ✅ El proyecto compila correctamente

### Total de código:
- **~600 líneas de Go**
- **8 archivos principales**
- **Estructura modular y escalable**

## 🎯 Cómo Ejecutar

### Opción 1: Script rápido
```bash
./run.sh
```

### Opción 2: Manual
```bash
# Compilar
go build -o unicli ./cmd/unicli

# Ejecutar
./unicli
```

### Opción 3: Desarrollo (sin compilar)
```bash
go run ./cmd/unicli
```

## 🎮 Controles Actuales

### Navegación:
- `1` o `t` - Vista de Tareas
- `2` o `c` - Vista de Calendario (placeholder)
- `3` o `s` - Vista de Clases (placeholder)
- `4` o `g` - Vista de Calificaciones (placeholder)
- `5` o `n` - Vista de Notas (placeholder)
- `6` - Estadísticas (placeholder)
- `7` - Configuración (placeholder)
- `q` o `Ctrl+C` - Salir

### En Vista de Tareas:
- `j` / `k` o `↓` / `↑` - Navegar entre tareas
- `g` - Ir al inicio
- `G` - Ir al final
- `space` - Toggle completado (placeholder)
- `n` - Nueva tarea (placeholder)
- `e` - Editar tarea (placeholder)
- `d` - Eliminar tarea (placeholder)

## 📋 Próximos Pasos

### Paso 1: Implementar Repositories (Alta prioridad)
Los repositories permitirán guardar y cargar datos reales de la base de datos.

**Archivos a crear:**
```
internal/database/repositories/base.go
internal/database/repositories/task_repo.go
internal/database/repositories/class_repo.go
```

**Qué hacer:**
1. Crear `base.go` con interfaz común
2. Implementar `TaskRepository` con métodos CRUD:
   - `Create(task *Task) error`
   - `FindByID(id string) (*Task, error)`
   - `FindAll() ([]Task, error)`
   - `Update(task *Task) error`
   - `Delete(id string) error`
   - `FindByStatus(status TaskStatus) ([]Task, error)`
   - `FindDueToday() ([]Task, error)`

3. Conectar TaskRepository con TaskScreen
4. Reemplazar datos de ejemplo con datos reales

### Paso 2: Formularios y Modales
Crear componentes para agregar/editar tareas.

**Archivos a crear:**
```
internal/ui/components/form.go
internal/ui/components/modal.go
internal/ui/components/input.go
```

**Funcionalidades:**
- Form para crear/editar tareas
- Modal de confirmación para eliminar
- Inputs con validación

### Paso 3: Vista de Calendario
Implementar la vista de calendario mensual.

**Archivo:**
```
internal/ui/screens/calendar.go
```

**Funcionalidades:**
- Vista mensual con días del mes
- Resaltar días con eventos/tareas
- Navegar entre meses
- Ver detalles del día seleccionado

### Paso 4: Vista de Clases y Horario
Gestión de horario de clases.

**Archivo:**
```
internal/ui/screens/classes.go
```

**Funcionalidades:**
- Vista semanal del horario
- CRUD de clases
- Asignar colores a clases
- Detectar conflictos de horario

### Paso 5: Búsqueda y Filtros
Sistema de búsqueda y filtrado avanzado.

**Funcionalidades:**
- Filtrar tareas por:
  - Estado (pending, completed, etc.)
  - Prioridad (urgent, high, medium, low)
  - Categoría
  - Tags
  - Rango de fechas
- Búsqueda fuzzy en títulos
- Guardar filtros favoritos

### Paso 6: Estadísticas y Dashboard
Vista con métricas y gráficos.

**Archivo:**
```
internal/ui/screens/stats.go
internal/services/stats_service.go
```

**Funcionalidades:**
- Tareas completadas esta semana/mes
- Productividad por día
- Distribución por prioridad
- Próximos deadlines
- Gráficos ASCII

## 🔧 Comandos Útiles

### Desarrollo:
```bash
# Ver cambios en tiempo real (requiere watchexec o similar)
watchexec -e go -r go run ./cmd/unicli

# Formatear código
go fmt ./...

# Verificar errores
go vet ./...

# Ejecutar tests (cuando los agreguemos)
go test ./...
```

### Base de Datos:
```bash
# Ver la base de datos (se crea en ~/.unicli/unicli.db)
sqlite3 ~/.unicli/unicli.db

# Queries útiles:
sqlite> .tables
sqlite> .schema tasks
sqlite> SELECT * FROM tasks;
```

## 📚 Recursos de Aprendizaje

### Bubble Tea:
- Tutorial oficial: https://github.com/charmbracelet/bubbletea/tree/master/tutorials
- Ejemplos: https://github.com/charmbracelet/bubbletea/tree/master/examples
- Lista de apps con Bubble Tea: https://github.com/charmbracelet/bubbletea#bubble-tea-in-the-wild

### Lipgloss (estilos):
- Docs: https://github.com/charmbracelet/lipgloss
- Ejemplos: https://github.com/charmbracelet/lipgloss/tree/master/examples

### SQLite en Go:
- modernc.org/sqlite docs: https://pkg.go.dev/modernc.org/sqlite

## 🎨 Personalización

### Cambiar colores:
Editar `internal/ui/styles/styles.go`:
```go
var (
    Primary   = lipgloss.Color("#7C3AED") // Tu color
    Secondary = lipgloss.Color("#06B6D4") // Tu color
    // ...
)
```

### Agregar nueva vista:
1. Crear archivo en `internal/ui/screens/`
2. Implementar interfaz `tea.Model`
3. Agregar en `internal/app/app.go`:
   - Enum `View`
   - Field en `Model`
   - Case en `Update()`
   - Case en `View()`
   - Item en sidebar

## 🐛 Debugging

### Ver logs:
```bash
# Ejecutar con output a archivo
./unicli 2> debug.log

# En otra terminal
tail -f debug.log
```

### Agregar logs en código:
```go
import "log"

log.Printf("Debug: %+v\n", variable)
```

## 📊 Arquitectura del Código

```
User Input → Bubble Tea → Update() → Services → Repositories → SQLite
                ↓                                              ↓
            View() ← Models ← ───────────────────────────────┘
```

**Flujo:**
1. Usuario presiona tecla
2. Bubble Tea genera mensaje (`tea.KeyMsg`)
3. `Update()` procesa mensaje
4. Si necesita datos, llama Service
5. Service usa Repository para acceder DB
6. Repository ejecuta SQL y retorna Models
7. `Update()` actualiza estado y retorna comando
8. `View()` renderiza nuevo estado
9. Usuario ve resultado

## ✨ Ideas Futuras

### Features avanzados:
- [ ] Plugin system (agregar features sin modificar core)
- [ ] Sincronización con Google Calendar/Outlook
- [ ] Modo colaborativo (compartir tareas con compañeros)
- [ ] Integración con Git (commits, PRs como tareas)
- [ ] Generación de reportes PDF
- [ ] Notificaciones de sistema
- [ ] Soporte para Markdown avanzado en notas
- [ ] Vim-mode avanzado (comandos :w, :q, etc.)
- [ ] Atajos personalizables
- [ ] Múltiples temas predefinidos (Dracula, Nord, Gruvbox)

### Integraciones:
- [ ] Import desde Notion
- [ ] Import desde Todoist
- [ ] Export a Markdown/HTML
- [ ] Webhook support
- [ ] API REST (opcional)

## 🤝 Contribuir

Para mantener calidad del código:

1. **Seguir la estructura** definida en ARCHITECTURE.md
2. **Escribir código idiomático Go**
3. **Agregar comentarios** en funciones públicas
4. **Testear** cambios importantes
5. **Mantener consistencia** de estilos

## 🆘 Ayuda

Si tienes problemas:

1. **Revisar ARCHITECTURE.md** - Explica cómo funciona todo
2. **Ver ejemplos de Bubble Tea** - Link arriba
3. **Revisar código de lazygit** - Gran referencia
4. **Leer docs de Go** - https://go.dev/doc/

---

## 🎉 ¡Felicidades!

Has creado las bases de una aplicación TUI completa y profesional. El proyecto está estructurado de forma que es fácil:
- ✅ Agregar nuevas vistas
- ✅ Extender funcionalidad
- ✅ Mantener el código organizado
- ✅ Escalar sin reescribir

**El siguiente paso más importante es implementar los Repositories para tener persistencia real de datos.**

¡A codear! 🚀
