[← Back to Main README](../../README.md) | [🇷🇺 Русская версия](zip-cmd_RU.md)

# 🗜 Archiving & Execution (`zip` & `cmd`)

## 🗜 `zip`: Archive Creation

Creates a standard `.zip` archive. Unlike basic archivers that flatten directories, `fflow` preserves the exact relative directory structure of the matched files using `filepath.ToSlash` for cross-platform compatibility.

**Example:**
```bash
fflow zip -r -e .go,.mod -o source_backup.zip
```

## ⚡ `cmd`: Shell Execution per File

Executes a custom shell command on **every matched file** concurrently. This is incredibly powerful for batch processing images, compiling code, or running linters.

**The `{}` Placeholder:**
Use `{}` in your `--exec` string. `fflow` will replace it with the **absolute path** of the current file, properly wrapped in quotes to handle spaces in filenames.

**Cross-Platform Execution:**
- On **Linux/macOS**, it wraps the command in `sh -c`.
- On **Windows**, it wraps the command in `cmd /c`.

**Examples:**

**1. Convert all PNGs to JPGs using ImageMagick:**
```bash
fflow cmd -e .png --exec "magick {} {}.jpg"
```

**2. Run a Go linter on all Go files:**
```bash
fflow cmd -r -e .go --exec "golangci-lint run {}"
```

**3. Extract all ZIP files in a directory:**
```bash
fflow cmd -e .zip --exec "unzip {} -d extracted/"
```

---
