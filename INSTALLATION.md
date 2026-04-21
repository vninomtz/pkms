# Instalación de PKMS

## Instalación global

Para instalar el comando `pkm` de forma global:

```bash
cd /Users/vnino/github.com/vninomtz/pkms
bash install.sh
```

O manualmente en `~/.local/bin`:

```bash
cd /Users/vnino/github.com/vninomtz/pkms
go build -o ~/.local/bin/pkm main.go
chmod +x ~/.local/bin/pkm
```

## Verificar instalación

```bash
pkm version
```

## Comandos disponibles

- `pkm add` - Agregar una nueva nota
- `pkm add-book` - Agregar un nuevo libro a la lista de lecturas
- `pkm search` - Buscar notas
- `pkm inspect` - Inspeccionar archivos
- `pkm publish` - Publicar contenido
- `pkm index` - Indexar notas
- `pkm version` - Mostrar versión

## Comando add-book

El comando `pkm add-book` proporciona una interfaz interactiva para agregar libros a tu lista de lecturas en Obsidian.

Ver `ADD_BOOK_USAGE.md` para más detalles.
