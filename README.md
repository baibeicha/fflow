<p align="right">
  <strong>🇬🇧 English</strong> | <a href="README_RU.md">🇷🇺 Русская версия</a>
</p>


<p align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white" />
  <img src="https://img.shields.io/badge/License-MIT-green?style=for-the-badge" />
  <img src="https://img.shields.io/badge/CLI-Powered-blue?style=for-the-badge&logo=gnubash" />
  <img src="https://img.shields.io/badge/Languages-EN%20|%20RU-purple?style=for-the-badge" />
</p>

<h1 align="center">🌊 fflow</h1>
<p align="center">
  <b>The Ultimate File Management, Automation, and Processing CLI Utility.</b><br>
  <i>Blazing fast directory traversal, advanced PDF generation, and powerful pipeline automation.</i>
</p>

---

## 📖 Table of Contents
- [✨ Key Features](#-key-features)
- [🚀 Installation](#-installation)
- [⚡ Quick Start](#-quick-start)
- [🌐 Global Filters & Flags](#-global-filters--flags)
- [🛠 Commands Overview](#-commands-overview)
- [📄 Deep Dive: PDF Generation](#-deep-dive-pdf-generation)
- [🔄 Deep Dive: Flow Pipelines](#-deep-dive-flow-pipelines)
- [🌍 Localization](#-localization)
- [📂 Documentation](#-documentation)

---

## ✨ Key Features

- **🚀 Blazing Fast Engine**: Utilizes concurrent worker pools (`ants`) and zero-allocation optimized directory traversal (`fastwalk`) to process millions of files in seconds.
- **📄 Advanced PDF Generator**: Convert source code, markdown, CSV tables, and images into beautifully formatted PDFs. Includes 9 localized presets (Academic, Scientific, Cheatsheet, etc.).
- **🔄 Pipeline Automation (`flow`)**: Chain multiple file operations (copy, rename, zip, execute) using simple YAML configurations or inline CLI strings with variable interpolation.
- **📊 Deep Analytics**: Calculate exact directory sizes (including nested folders), line counts, word counts, and character statistics.
- **🛠 Bulk Operations**: Copy, move, delete, rename, change extensions, and create archives with surgical precision.
- **🌍 Built-in Localization**: Seamlessly switch between **English** and **Russian** interfaces.
- **⚙️ Highly Configurable**: Viper-based configuration, environment variable support, and extensive CLI flags for every edge case.

---

## 🚀 Installation

### From Source
Requires **Go 1.21** or higher.

```bash
git clone https://github.com/baibeicha/fflow.git
cd fflow
go build -o fflow ./cmd/fflow
# Move to your PATH (Linux/macOS)
sudo mv fflow /usr/local/bin/ 
```

### Go Install
```bash
go install github.com/baibeicha/fflow/cmd/fflow@latest
```

---

## ⚡ Quick Start

**1. Analyze a directory and sort by size:**
```bash
fflow info ./my_project --sort-by size:desc
```

**2. Merge all `.log` files into a single file with headers:**
```bash
fflow merge -e .log -o combined_logs.txt --include-path --separator "\n---\n"
```

**3. Generate a PDF from your Go source code:**
```bash
fflow pdf -e .go --preset SourceCode --lang en -o codebase.pdf
```

**4. Backup and rename files in one pipeline:**
```bash
fflow flow -c "var BACKUP=./backup ; copy -d ${BACKUP} ; rename --prefix backup_"
```

---

## 🌐 Global Filters & Flags

Almost all `fflow` commands support a powerful set of global filters to target exactly the files you need.

| Flag | Short | Description | Example |
| :--- | :---: | :--- | :--- |
| `--path` | `-p` | Target directories (comma-separated or multiple flags). | `-p ./src,./docs` |
| `--recursive` | `-r` | Traverse subdirectories. | `-r` |
| `--extensions` | `-e` | Filter by file extensions. | `-e .go,.md,.txt` |
| `--blacklist` | | Exclude specific files or directories. | `--blacklist node_modules,.git` |
| `--min-size` | | Minimum file size (supports `b, kb, mb, gb, tb, pb`). | `--min-size 1mb` |
| `--max-size` | | Maximum file size. | `--max-size 500kb` |
| `--sort-by` | | Sort results. Format: `field:order`. | `--sort-by size:desc,name:asc` |
| `--quiet` | `-q` | Suppress UI, progress bars, and colors. | `-q` |
| `--verbose` | `-v` | Enable verbose logging. | `-v` |

---

## 🛠 Commands Overview

| Command | Description | Documentation                                                  |
| :--- | :--- |:---------------------------------------------------------------|
| **`info`** | Analyze directories, calculate accumulated sizes, and list files. | [docs/commands/info.md](docs/commands/info.md)                 |
| **`stats`** | Count lines, words, and characters across multiple files. | [docs/commands/stats-merge.md](docs/commands/stats-merge.md)   |
| **`merge`** | Concatenate files with custom separators and headers. | [docs/commands/stats-merge.md](docs/commands/stats-merge.md)   |
| **`pdf`** | Generate highly customizable PDFs from code, text, and tables. | [docs/commands/pdf.md](docs/commands/pdf.md)                   |
| **`copy`** | Copy files to single or multiple destinations. | [docs/commands/copy-move.md](docs/commands/copy-move.md)       |
| **`move`** | Move files with collision resolution. | [docs/commands/copy-move.md](docs/commands/copy-move.md)       |
| **`delete`** | Safely delete files based on filters. | [docs/commands/delete.md](docs/commands/delete.md)             |
| **`rename`** | Bulk rename using search/replace, prefixes, and suffixes. | [docs/commands/rename-ext.md](docs/commands/rename-ext.md)     |
| **`ext`** | Change file extensions in bulk. | [docs/commands/rename-ext.md](docs/commands/rename-ext.md)     |
| **`zip`** | Create `.zip` archives preserving directory structures. | [docs/commands/zip-cmd.md](docs/commands/zip-cmd.md)           |
| **`cmd`** | Execute shell commands on every matched file. | [docs/commands/zip-cmd.md](docs/commands/zip-cmd.md)           |
| **`flow`** | Execute complex multi-step automation pipelines. | [docs/commands/flow.md](docs/commands/flow.md)                 |
| **`locale`** | Switch UI language (`en` or `ru`). | [docs/commands/localization.md](docs/commands/localization.md) |
---

## 📄 Deep Dive: PDF Generation

`fflow pdf` is not just a simple text-to-PDF converter. It is a full-fledged document generator that understands file types (Code, Markdown, CSV, Images) and applies professional typographic rules.

### 🎨 Built-in Presets
Use the `--preset` flag to instantly apply complex configurations. Every preset is localized for **EN** and **RU**.

| Preset Name | Best For | Key Characteristics |
| :--- | :--- | :--- |
| `A4Portrait` | General lists, text files | Standard margins, centered captions. |
| `A4Landscape` | Wide CSV tables | Horizontal layout, left-aligned data captions. |
| `AcademicReport` | University labs, drafts | 14pt Arial, 20mm margins, "Table - X" / "Figure - X" captions. |
| `ScientificArticle` | VAK / RINC publications | Strict formatting, right-aligned table captions, 12pt font. |
| `SourceCode` | Backend codebases, configs | 9pt Courier, gray backgrounds, disabled image/table captions. |
| `ServerLogs` | Extensive trace events | Dense landscape layout, 7pt font, minimal margins. |
| `Cheatsheet` | Quick reference guides | Maximum data density, 6pt font, markdown enabled. |
| `A3DataHeavy` | Massive datasets | A3 Landscape, "[EXPORT]" prefixes, stretched tables. |
| `PhotoAlbum` | Image galleries | Landscape, centered bottom captions, large image rendering. |

### ⚙️ Granular Customization
Override any preset parameter using CLI flags:
```bash
fflow pdf -e .md \
  --preset AcademicReport \
  --lang ru \
  --orientation landscape \
  --margin 15 \
  --font-size 12 \
  --code-bg-color "240,240,240" \
  --table-stretch-width \
  -o my_report.pdf
```
*(See [docs/commands/pdf.md](docs/commands/pdf.md) for the full list of 40+ PDF configuration flags).*

---

## 🔄 Deep Dive: Flow Pipelines

The `flow` command allows you to chain operations. The output of one step becomes the input of the next.

### 1. Inline CLI Pipelines
Use semicolons `;` to separate steps. Use `var` to define variables.
```bash
fflow flow -e .jpg \
  -c "var DEST=./backup/images ; copy -d ${DEST} ; rename --prefix vacation_ ; zip -o vacation_archive.zip"
```

### 2. YAML Pipelines
For complex logic, use a YAML file (`pipeline.yaml`):
```yaml
env:
  BACKUP_DIR: "./daily_backups"
  PREFIX: "prod_"

steps:
  - action: copy
    args:
      dest: "${BACKUP_DIR}"
      rewrite: "true"
      
  - action: rename
    args:
      prefix: "${PREFIX}"
      
  - action: ext
    args:
      to: ".bak"
      
  - action: zip
    args:
      output: "${BACKUP_DIR}/archive.zip"
```
Run it with:
```bash
fflow flow -f pipeline.yaml --env
```
*(The `--env` flag injects your system environment variables into the pipeline context).*

---

## 🌍 Localization

`fflow` supports full UI translation. By default, it detects your system locale or defaults to English.

**Change language to Russian:**
```bash
fflow locale ru
```
**Change language to English:**
```bash
fflow locale en
```
*Note: The language setting is persisted in your OS user config directory (e.g., `~/.config/fflow/locale.yaml`).*

---

## 📂 Documentation

For exhaustive details on every flag, edge case, and internal mechanic, please refer to the `/docs` directory:

- **[Global Flags & Sorting](docs/global-flags.md)**
- **[PDF Generation Guide](docs/commands/pdf.md)**
- **[Flow Automation Guide](docs/commands/flow.md)**
- **[File Operations (Copy/Move)](docs/commands/copy-move.md)**
- **[File Operations (Delete)](docs/commands/delete.md)**
- **[Text & Stats Operations](docs/commands/stats-merge.md)**

---
<p align="center">
  <b>Made with ❤️ and Go.</b>
</p>

---
