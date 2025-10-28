# 🎨 Kanban View - Implemented!

## ✅ Vista Kanban Implementada

Se ha transformado completamente la vista de Tasks a un formato **Kanban Board** profesional con 3 columnas.

---

## 📊 Layout

```
┌─────────────────────────────────────────────────────────────┐
│  To Do (3)     │   In Progress (2)   │      Done (1)        │
├────────────────┼─────────────────────┼──────────────────────┤
│                │                     │                      │
│ 🔴 Task 1      │ 🟡 Task 4          │ 🔵 Task 6           │
│ 🟡 Task 2      │ 🔵 Task 5          │                      │
│ 🔵 Task 3      │                     │                      │
│                │                     │                      │
└────────────────┴─────────────────────┴──────────────────────┘

Total: 6 tasks
```

---

## 🎮 Controles Implementados

### Navegación entre columnas:
- **`Tab`** - Ir a la siguiente columna (To Do → In Progress → Done)
- **`Shift+Tab`** - Ir a la columna anterior (Done → In Progress → To Do)

### Navegación dentro de columna:
- **`j` / `↓`** - Bajar cursor
- **`k` / `↑`** - Subir cursor
- **`g`** - Ir al inicio de la columna
- **`G`** - Ir al final de la columna

### Selección y movimiento:
- **`Enter`** - Seleccionar tarea (se marca con highlight)
- **`Enter`** (en tarea seleccionada) - Deseleccionar
- **`←` / `h`** - Mover tarea seleccionada a columna anterior
- **`→` / `l`** - Mover tarea seleccionada a columna siguiente

### Acciones:
- **`Delete` / `Backspace`** - Eliminar tarea seleccionada
- **`n`** - Nueva tarea (TODO: implementar formulario)
- **`e`** - Editar tarea seleccionada (TODO: implementar formulario)
- **`r`** - Refrescar datos

---

## 🎨 Características Visuales

### Indicadores de Prioridad:
- 🔴 Urgent (Urgente)
- 🟡 High (Alta)
- 🔵 Medium (Media)
- ⚪ Low (Baja)

### Indicadores de Fecha:
- ⚠ Overdue (Vencida)
- 📅 Today (Hoy)

### Estados de Tarea:
- **Cursor** - Fondo gris claro (navegación)
- **Seleccionada** - Fondo verde (lista para mover/editar/eliminar)
- **Completada** - Tachado

### Columnas:
- **Activa** - Borde y título resaltado
- **Inactiva** - Borde normal

### Headers:
- Muestran título y cantidad: "To Do (3)"

---

## 💾 Mapeo de Estados

| Columna       | Estado en DB             |
|---------------|--------------------------|
| To Do         | `TaskStatusPending`      |
| In Progress   | `TaskStatusInProgress`   |
| Done          | `TaskStatusCompleted`    |

---

## 🔄 Flujo de Trabajo

### Crear nueva tarea:
1. Presionar `n` (se abrirá formulario - TODO)
2. La tarea se crea en columna **To Do**
3. Usuario puede moverla con Enter + flechas

### Mover tarea entre columnas:
1. Navegar hasta la tarea con `j/k`
2. Presionar `Enter` para seleccionar (se marca)
3. Usar `←` o `→` para mover entre columnas
4. La tarea cambia automáticamente de estado en la DB

### Eliminar tarea:
1. Seleccionar tarea con `Enter`
2. Presionar `Delete` o `Backspace`
3. Tarea se elimina permanentemente

---

## 🎯 Ventajas de Kanban

✅ **Visual** - Estado de tareas a simple vista
✅ **Intuitivo** - Drag & drop con teclado
✅ **Organizado** - Separación clara de trabajo
✅ **Productivo** - Flujo de trabajo visible
✅ **Profesional** - Estilo moderno como Trello/Jira

---

## 📝 Mejoras Futuras

### Corto plazo:
- [ ] Formulario para crear tareas (tecla `n`)
- [ ] Formulario para editar tareas (tecla `e`)
- [ ] Confirmación antes de eliminar
- [ ] Soporte para categorías/tags visuales

### Medio plazo:
- [ ] Filtros por prioridad
- [ ] Búsqueda de tareas
- [ ] Ordenamiento (por fecha, prioridad)
- [ ] Vista compacta/expandida

### Largo plazo:
- [ ] Columnas personalizables
- [ ] Límite WIP (Work In Progress)
- [ ] Tiempo en cada columna (cycle time)
- [ ] Estadísticas por columna

---

## 🐛 Testing

Para probar la vista Kanban:

```bash
# 1. Asegurar que hay tareas en diferentes estados
go run ./cmd/seed/main.go

# 2. Ejecutar la app
./unicli

# 3. Ir a Tasks (presionar :s, luego seleccionar Tasks)

# 4. Probar:
# - Tab para cambiar columnas
# - j/k para navegar
# - Enter para seleccionar
# - ← → para mover entre columnas
# - Delete para eliminar
```

---

## 🔧 Cambios en el Código

### Archivos modificados:
- `internal/ui/screens/tasks.go` - Reescrito completamente para Kanban
  - Agregado: Column enum (Todo, InProgress, Done)
  - Agregado: cursors por columna
  - Agregado: selectedTaskID
  - Reescrito: View() con renderColumn()
  - Reescrito: Update() con nueva navegación
  - Nuevo: moveTaskToColumn()
  - Nuevo: getTasksForColumn()

- `internal/database/repositories/task_repo.go` - Mejorado Update()
  - Auto-set CompletedAt cuando status = completed
  - Auto-clear CompletedAt cuando status != completed

### Líneas de código:
- **~500 líneas** reescritas/agregadas
- Arquitectura mantenida limpia
- Separación de concerns

---

## 📚 Inspiración

Diseño inspirado en:
- **Trello** - Columnas y movimiento de cards
- **Jira** - Estados de workflow
- **GitHub Projects** - Kanban simple
- **lazygit** - Navegación con teclado

---

*Vista Kanban completamente funcional! 🎉*
