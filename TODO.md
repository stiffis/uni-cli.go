# TODO - UniCLI Development Tasks

## 🔴 Alta Prioridad (MVP)

### Database Layer
- [ ] Implementar `base.go` con interfaz Repository base
- [ ] Implementar `TaskRepository` completo
  - [ ] Create
  - [ ] FindByID
  - [ ] FindAll
  - [ ] Update
  - [ ] Delete
  - [ ] FindByStatus
  - [ ] FindDueToday
  - [ ] FindUpcoming
- [ ] Conectar TaskRepository con TaskScreen
- [ ] Reemplazar datos de muestra con datos reales de DB

### Task Management (CRUD completo)
- [ ] Crear componente Form para nueva tarea
- [ ] Crear componente Modal para confirmaciones
- [ ] Implementar crear tarea (tecla `n`)
- [ ] Implementar editar tarea (tecla `e`)
- [ ] Implementar eliminar tarea (tecla `d`)
- [ ] Implementar toggle completado (tecla `space`)
- [ ] Validación de formularios

### UI Components
- [ ] Input field component
- [ ] Textarea component
- [ ] Select/dropdown component
- [ ] Date picker component
- [ ] Modal/dialog component
- [ ] Confirmation dialog

## 🟡 Media Prioridad

### Calendar View
- [ ] Crear `calendar.go` screen
- [ ] Implementar vista mensual
- [ ] Resaltar días con eventos
- [ ] Navegación entre meses (h/l o arrows)
- [ ] Vista de detalles del día seleccionado
- [ ] Agregar eventos desde calendario

### Classes & Schedule
- [ ] Crear `classes.go` screen
- [ ] Implementar ClassRepository
- [ ] Vista de lista de clases
- [ ] CRUD de clases
- [ ] Vista de horario semanal
- [ ] Detección de conflictos de horario
- [ ] Asignación de colores

### Search & Filters
- [ ] Implementar búsqueda fuzzy en tareas
- [ ] Filtrar por estado
- [ ] Filtrar por prioridad
- [ ] Filtrar por categoría
- [ ] Filtrar por tags
- [ ] Filtrar por rango de fechas
- [ ] Guardar filtros favoritos

## 🟢 Baja Prioridad

### Grades Management
- [ ] Crear `grades.go` screen
- [ ] Implementar GradeRepository
- [ ] CRUD de calificaciones
- [ ] Calcular promedio por clase
- [ ] Calcular GPA general
- [ ] Gráfico de evolución de notas

### Notes
- [ ] Crear `notes.go` screen
- [ ] Implementar NoteRepository
- [ ] CRUD de notas
- [ ] Editor markdown básico
- [ ] Preview de markdown
- [ ] Sistema de tags para notas
- [ ] Búsqueda en notas

### Statistics Dashboard
- [ ] Crear `stats.go` screen
- [ ] Implementar StatsService
- [ ] Tareas completadas esta semana
- [ ] Tareas completadas este mes
- [ ] Gráfico de productividad por día
- [ ] Distribución por prioridad (gráfico de barras ASCII)
- [ ] Próximos deadlines importantes
- [ ] Tiempo promedio para completar tareas

### Pomodoro Timer
- [ ] Crear componente Timer
- [ ] Configuración de duración (trabajo/descanso)
- [ ] Notificación al terminar
- [ ] Historial de sesiones
- [ ] Estadísticas de tiempo de estudio

## 🔵 Features Adicionales

### Configuration
- [ ] Crear `settings.go` screen
- [ ] Selector de temas
- [ ] Configurar atajos de teclado
- [ ] Preferencias de visualización
- [ ] Configurar notificaciones
- [ ] Configurar formato de fechas

### Import/Export
- [ ] Exportar a JSON
- [ ] Exportar a Markdown
- [ ] Exportar a CSV
- [ ] Importar desde JSON
- [ ] Importar desde Todoist
- [ ] Backup automático

### Themes
- [ ] Implementar sistema de temas
- [ ] Tema Dracula
- [ ] Tema Nord
- [ ] Tema Gruvbox
- [ ] Tema Solarized
- [ ] Editor de temas personalizado

### Help System
- [ ] Crear `help.go` screen
- [ ] Documentación de atajos por vista
- [ ] Tutorial interactivo (primera vez)
- [ ] Tips rápidos

## 🧪 Testing

### Unit Tests
- [ ] Tests para models (Task, Class, etc.)
- [ ] Tests para services
- [ ] Tests para repositories
- [ ] Tests para validators

### Integration Tests
- [ ] Test de flujo completo de tareas
- [ ] Test de migraciones de DB
- [ ] Test de importación/exportación

### UI Tests
- [ ] Golden tests para screenshots de UI
- [ ] Tests de navegación

## 📦 Distribution

### Build & Release
- [ ] Script de build cross-platform
- [ ] GitHub Actions CI/CD
- [ ] Release automation
- [ ] Generar binarios para Linux/Mac/Windows
- [ ] Crear instaladores

### Package Managers
- [ ] Homebrew formula
- [ ] AUR package (Arch Linux)
- [ ] apt repository (Debian/Ubuntu)
- [ ] Snap package
- [ ] Chocolatey (Windows)

### Documentation
- [ ] User guide completo
- [ ] Developer guide
- [ ] API documentation (si se agrega)
- [ ] Video tutorial
- [ ] Screenshots/GIFs para README

## 🐛 Bug Fixes & Improvements

### Known Issues
- [ ] Manejo de terminal resize en tiempo real
- [ ] Scroll en listas largas
- [ ] Performance con muchos items (1000+)
- [ ] Manejo de errores más robusto

### Code Quality
- [ ] Agregar linter (golangci-lint)
- [ ] Mejorar manejo de errores
- [ ] Agregar logging apropiado
- [ ] Refactorizar código duplicado
- [ ] Mejorar comentarios y documentación

## 🎨 UI/UX Improvements

### Visual
- [ ] Animaciones suaves (fade in/out)
- [ ] Loading indicators
- [ ] Progress bars
- [ ] Better error messages
- [ ] Icons más descriptivos

### Navigation
- [ ] Command palette (Ctrl+P)
- [ ] Quick actions menu
- [ ] Breadcrumbs
- [ ] Recent views history
- [ ] Bookmarks/favorites

### Accessibility
- [ ] Soporte para lectores de pantalla
- [ ] Alto contraste mode
- [ ] Keyboard navigation mejorado
- [ ] Configuración de tamaño de texto

---

## Orden Sugerido de Implementación

1. **Semana 1**: Database Repositories + CRUD de Tareas
2. **Semana 2**: Forms, Modals e Inputs
3. **Semana 3**: Vista de Calendario
4. **Semana 4**: Vista de Clases y Horario
5. **Semana 5**: Search & Filters
6. **Semana 6**: Grades Management
7. **Semana 7**: Notes & Markdown
8. **Semana 8**: Statistics Dashboard
9. **Semana 9**: Polish + Testing
10. **Semana 10**: Documentation & Release

---

*Actualizar este archivo conforme avanza el desarrollo*
