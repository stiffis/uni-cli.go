# 🧪 Instrucciones de Testing - Creación de Tareas

## Problema Reportado
La función de crear tareas no está guardando en la base de datos en Arch Linux.

## Cambios Realizados
1. ✅ Corregido error de sintaxis en `taskform.go` (método GetTask mal formateado)
2. ✅ Agregado logging de debugging en `/tmp/unicli_taskform_debug.log`
3. ✅ Recompilado el binario

## Cómo Probar

### Paso 1: Limpiar logs anteriores
```bash
rm -f /tmp/unicli_taskform_debug.log
```

### Paso 2: Ejecutar la aplicación
```bash
cd /home/stiff/UniCLI
./unicli
```

### Paso 3: Probar creación de tarea
1. Presiona `:s` para abrir el sidebar
2. Navega a "Tasks" con j/k y presiona Enter
3. Presiona `n` para crear nueva tarea
4. Llena el formulario:
   - **Title**: "Test Task Debug 123"
   - Presiona `Tab` para ir a Description
   - **Description**: "Testing task creation"
   - Presiona `Tab` para ir a Due Date (opcional, déjalo vacío o pon: 2025-11-01)
   - Presiona `Tab` para ir a Priority
   - Usa `←` / `→` para cambiar prioridad si quieres
   - Presiona `Tab` para ir al botón "Create"
   - **Presiona `Enter` para enviar**

5. Deberías ver la tarea aparecer en la columna "To Do"
6. Presiona `:q` para salir

### Paso 4: Verificar en la base de datos
```bash
sqlite3 ~/.unicli/unicli.db "SELECT title, status, created_at FROM tasks WHERE title LIKE '%Debug%' ORDER BY created_at DESC;"
```

### Paso 5: Revisar los logs
```bash
cat /tmp/unicli_taskform_debug.log
```

Deberías ver algo como:
```
[TASKFORM] SUBMIT: focusedField=4, titleValue='Test Task Debug 123'
[TASKFORM] Form submitted successfully
[TASKFORM] GetTask called: taskID='', titleValue='Test Task Debug 123'
[TASKFORM] Creating NEW task with ID: [algún-uuid]
[TASKFORM] Task prepared: ID=[uuid], Title='Test Task Debug 123', Status=pending, Priority=medium
```

## Verificaciones Adicionales

### Ver todas las tareas actuales:
```bash
sqlite3 ~/.unicli/unicli.db "SELECT id, title, status FROM tasks ORDER BY created_at DESC LIMIT 5;"
```

### Contar tareas totales:
```bash
sqlite3 ~/.unicli/unicli.db "SELECT COUNT(*) FROM tasks;"
```

### Ver permisos de la base de datos:
```bash
ls -l ~/.unicli/unicli.db
```

## Si No Funciona

1. Verifica que el log se esté creando:
```bash
ls -l /tmp/unicli_taskform_debug.log
```

2. Verifica que el binario esté actualizado:
```bash
ls -lh /home/stiff/UniCLI/unicli
```

3. Asegúrate de presionar Enter cuando el cursor esté en el botón "[Create]"
   - El botón debe estar resaltado (fondo verde)
   - Si no lo está, presiona Tab hasta que lo esté

4. Si el log muestra "Title is empty", entonces el problema es con el input field de Bubble Tea

## Reportar Resultados

Por favor copia y pega:
1. La salida del comando de verificación en DB
2. El contenido completo de `/tmp/unicli_taskform_debug.log`
3. Cualquier error que veas en la terminal

