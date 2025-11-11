
## 🚀 Mejoras y extensiones sugeridas

### 🏷️ 1. **Etiquetas (tags o categorías)**

Permite agrupar tareas por materia o tipo:

```
📘 #calculus
📙 #discretemath
📗 #personal
```

Así puedes filtrar (`/tag calculus`) o crear vistas por curso.

💡 *Integración futura:* conectar etiquetas con la sección **Classes**.

---

### 🧠 2. **Nivel de prioridad numérico**

Además de los iconos de color, podrías guardar internamente un campo `priority` (1–3) o `urgency` calculado según fecha límite + prioridad.
Ejemplo de visual:

```
🔥 P1  Submit report
⭐ P2  Read paper
```

---

### 💬 5. **Descripción o subtareas**

Posibilidad de expandir una tarea para ver detalles:

```
Read Chapter 5
 ├─ pages 1–20
 ├─ highlight key terms
 └─ summarize in Notes
```

Podrías mostrarlo al presionar `enter`.

```
Study for calculus exam  [███░░░░░] 40%
```

---

### 🎨 9. **Personalización de vista**

Permitir alternar entre:

* **Kanban** (actual)
* **Lista simple** (compacta)
* **Vista por fecha** (agenda de deadlines)

---

### 🧩 10. **Filtros y búsqueda**

Comandos útiles:

```
/search "calculus"
/due next7d
/tag math
/show completed
```

Esto le da poder de CLI real, estilo `taskwarrior`.

---
