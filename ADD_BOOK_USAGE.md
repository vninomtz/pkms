# Comando add-book

Comando interactivo para agregar libros a tu lista de lecturas en Obsidian.

## Compilación

```bash
cd /Users/vnino/github.com/vninomtz/pkms
go build -o bin/pkms main.go
```

## Uso

```bash
./bin/pkms add-book
```

## Flujo del comando

1. **Título del libro** - Campo obligatorio
2. **Autor** - Campo obligatorio
3. Validación de duplicados:
   - Si el libro existe: ofrece agregar una "relectura"
   - Si no existe: continúa con el flujo normal
4. **¿Ya lo completaste?** - Define si es completado o pendiente
5. **Idioma** - Default: "es" (español)
6. **Categorías** - Opcional, separadas por comas
7. **URL** - Opcional
8. **Notas** - Opcional
9. **Preview** - Muestra resumen de datos
10. **Confirmación** - Requiere confirmación antes de guardar

## Características

- ✅ Validación de duplicados
- ✅ Auto-generación de UUIDs (8 caracteres)
- ✅ Fecha automática (hoy) para libros completados
- ✅ Actualización automática de metadata
- ✅ Preview de datos antes de confirmar
- ✅ Manejo de relecturas

## Archivo modificado

El comando modifica:
`$HOME/Library/Mobile Documents/iCloud~md~obsidian/Documents/notes/reading-resources.yml`

## Estructura YAML

El comando mantiene la siguiente estructura:

```yaml
metadata:
  title: Reading Resources
  description: ...
  created: YYYY-MM-DD
  updated: YYYY-MM-DD
  total_resources: N
  breakdown:
    books:
      completed: N
      pending: N
    blogs:
      completed: N
      pending: N

resources:
  - id: xxxxxxxx
    type: book
    status: completed|pending
    title: "..."
    author: "..."
    language: es|en
    categories: [cat1, cat2, ...]
    completed_date: YYYY-MM-DD|null
    url: "..." (opcional)
    notes: "..."
```

## Ejemplo

```bash
$ ./bin/pkms add-book

📚 Agregar un nuevo libro a tus lecturas

Título del libro: El Quijote
Autor: Miguel de Cervantes
¿Ya lo completaste? (s/n): n
Idioma [es]: es
Categorías (separadas por comas, opcional): literature,classics,spanish
URL (opcional): 
Notas (opcional): Lectura recomendada por Pedro

==================================================
📋 PREVIEW
==================================================
Título: El Quijote
Autor: Miguel de Cervantes
Estado: Pendiente
Idioma: es
Categorías: literature, classics, spanish
Notas: Lectura recomendada por Pedro
==================================================

¿Confirmar? (s/n): s

✅ Libro agregado exitosamente

📚 Título: El Quijote
✍️  Autor: Miguel de Cervantes
📖 UUID: a1b2c3d4
🔖 Estado: Pendiente
```
