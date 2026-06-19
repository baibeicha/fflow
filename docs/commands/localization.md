[← Back to Main README](../../README.md) | [🇷🇺 Русская версия](localization_RU.md)

# 🌍 Localization System

`fflow` features a fully localized CLI interface, currently supporting **English (`en`)** and **Russian (`ru`)**. 

## 🔄 Changing the Language

Use the `locale` command to switch the UI language. This setting is persisted across sessions.

```bash
# Switch to Russian
fflow locale ru

# Switch to English
fflow locale en
```

## ⚙️ How it Works Under the Hood

1. **Embedded Messages**: All translation strings are stored in YAML files and embedded directly into the Go binary using `//go:embed`. This means `fflow` requires zero external dependencies or configuration files to run.
2. **Persistent Config**: When you run `fflow locale ru`, the utility saves your preference to a `locale.yaml` file in your OS's standard user configuration directory:
    - **Linux**: `~/.config/fflow/locale.yaml`
    - **macOS**: `~/Library/Application Support/fflow/locale.yaml`
    - **Windows**: `%AppData%\fflow\locale.yaml`
3. **Auto-Detection**: On startup, if no config file exists, `fflow` defaults to `en`.

## 📝 Translating PDF Presets
The localization system also extends to the `pdf` command. Built-in presets (like `AcademicReport` or `ScientificArticle`) will automatically use the correct terminology for captions (e.g., "Table - " vs "Таблица - ") based on the `--lang` flag passed to the PDF generator, independent of the UI locale.

---
