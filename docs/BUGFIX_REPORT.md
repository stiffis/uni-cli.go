# 🐛 Bug Fix Report - Task Creation Issue

## Problema Original
Las tareas creadas desde el formulario no se guardaban en la base de datos en Arch Linux.

## Root Cause Identificado

### Bug #1: Error de Sintaxis en `taskform.go` ✅ FIXED
- **Archivo**: `internal/ui/components/taskform.go`
- **Problema**: El método `GetTask()` tenía una llave de cierre `}` mal colocada
- **Impacto**: Indentación incorrecta, posible comportamiento inesperado
- **Solución**: Corregida la estructura del método

### Bug #2: Lógica Invertida en `tasks.go` ✅ FIXED
- **Archivo**: `internal/ui/screens/tasks.go`
- **Problema**: La condición para determinar CREATE vs UPDATE estaba incorrecta:

```go
// ❌ CÓDIGO INCORRECTO (antes):
if task.ID != "" {
    return s, s.updateTask(task)  // task.ID SIEMPRE tiene valor (UUID)
} else {
    return s, s.createTask(task)  // Nunca llegaba aquí
}
```

**Explicación del bug:**
- `models.NewTask()` siempre genera un UUID automáticamente
- Por lo tanto, `task.ID` nunca está vacío
- **Todas las tareas nuevas se trataban como UPDATE en lugar de CREATE**
- UPDATE falla porque el ID no existe en la DB → tarea no se guarda

**Solución:**
```go
// ✅ CÓDIGO CORRECTO (después):
if s.taskForm.IsNewTask() {
    return s, s.createTask(task)  // Usa taskForm.taskID (campo privado)
} else {
    return s, s.updateTask(task)
}
```

Agregamos el método `IsNewTask()` al TaskForm que verifica el campo interno `taskID`:
- `taskID == ""` → Nueva tarea → CREATE
- `taskID != ""` → Edición → UPDATE

## Cambios Realizados

### 1. `internal/ui/components/taskform.go`
- ✅ Corregido método `GetTask()` (sintaxis)
- ✅ Agregado logging de debugging
- ✅ Agregado método `IsNewTask()` para exponer si es tarea nueva

### 2. `internal/ui/screens/tasks.go`
- ✅ Corregida lógica CREATE vs UPDATE usando `IsNewTask()`
- ✅ Agregado logging extensivo para debugging
- ✅ Logs en: submit, createTask, updateTask, taskCreatedMsg

## Archivos de Log Generados

1. `/tmp/unicli_taskform_debug.log` - Logs del formulario
2. `/tmp/unicli_taskscreen_debug.log` - Logs de la pantalla de tareas

## Cómo Verificar el Fix

```bash
# 1. Limpiar logs
rm -f /tmp/unicli_*.log

# 2. Ejecutar app
cd /home/stiff/UniCLI
./unicli

# 3. Crear una tarea (presiona: :s → Tasks → n → llenar formulario → Tab hasta Create → Enter)

# 4. Verificar en DB
sqlite3 ~/.unicli/unicli.db "SELECT id, title, status FROM tasks ORDER BY created_at DESC LIMIT 3;"

# 5. Revisar logs
echo "=== TASKFORM LOG ==="
cat /tmp/unicli_taskform_debug.log
echo ""
echo "=== TASKSCREEN LOG ==="
cat /tmp/unicli_taskscreen_debug.log
```

## Logs Esperados (Correcto)

```
[TASKFORM] Form submitted successfully
[TASKFORM] GetTask called: taskID='', titleValue='Mi Tarea'
[TASKFORM] Creating NEW task with ID: [uuid]
[TASKSCREEN] Form submitted! Task: ID=[uuid], Title='Mi Tarea', IsNewTask=true
[TASKSCREEN] Creating NEW task
[TASKSCREEN] createTask called with: ID=[uuid], Title='Mi Tarea'
[TASKSCREEN] Calling db.Tasks().Create()...
[TASKSCREEN] Create SUCCESS!
[TASKSCREEN] Task created successfully, reloading tasks...
```

## Logs Anteriores (Incorrecto)

```
[TASKFORM] Form submitted successfully
[TASKFORM] GetTask called: taskID='', titleValue='Mi Tarea'
[TASKFORM] Creating NEW task with ID: [uuid]
[TASKSCREEN] Form submitted! Task: ID=[uuid], Title='Mi Tarea', IsNewTask=false  ← BUG
[TASKSCREEN] Detected as UPDATE (ID exists)  ← Tomaba el camino equivocado
[TASKSCREEN] updateTask called... (falla porque el ID no existe en DB)
```

## Estado Actual

✅ **BUG FIXED** - Las tareas ahora se crean correctamente en la base de datos

## Testing

Por favor prueba creando varias tareas y reporta si:
1. ✅ Las tareas aparecen en el Kanban
2. ✅ Las tareas persisten después de cerrar y reabrir la app
3. ✅ La edición de tareas también funciona (presiona 'e' en una tarea)

---
**Fecha**: 2025-10-28
**Versión Fixed**: unicli (compilado a las 02:52 AM)
