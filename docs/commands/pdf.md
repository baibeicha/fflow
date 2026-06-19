[← Back to Main README](../../README.md) | [🇷🇺 Русская версия](pdf_RU.md)

# 📄 PDF Generation Guide

The `fflow pdf` command is a powerhouse for document generation. It reads files, categorizes them (Code, Text, Markdown, Table, Image), and renders them into a cohesive PDF document using the `go-pdf/fpdf` library.

## 🧠 How it Works
1. **File Collection**: Gathers files based on your global filters (`-e`, `-r`, `--min-size`, etc.).
2. **Categorization**: Maps extensions to types (e.g., `.go` -> Code, `.csv` -> Table, `.md` -> Markdown).
3. **Parallel Reading**: Reads file contents concurrently using worker pools to maximize I/O throughput.
4. **Rendering**: Applies the selected Preset or Custom Configuration to render pages, handling page breaks, captions, and typography automatically.

## 🎨 Presets & Localization
Presets are defined in `pkg/files/pdf/presets`. Each preset contains an `EN` and `RU` configuration.
*Example: The `AcademicReport` preset in Russian will automatically use "Таблица - " for table captions and "Рисунок - " for images, while the English version uses "Table - " and "Figure - ".*

## ⚙️ Complete Flag Reference

### General Layout
| Flag | Description | Default |
| :--- | :--- | :--- |
| `--preset` | Choose a built-in template. | `A4Portrait` |
| `--lang` | Localization for captions (`en`, `ru`). | `ru` |
| `--output`, `-o` | Output PDF file path. | `output.pdf` |
| `--orientation` | `portrait` or `landscape`. | Preset default |
| `--page-size` | `A4`, `A3`, `A5`. | Preset default |
| `--margin` | Page margins in mm. | Preset default |
| `--font-name` | Base font (must be embedded or standard). | `Arial` |
| `--font-size` | Base font size in pt. | `10.0` |
| `--font-path` | Path to custom `.ttf` font file. | Empty (Uses embedded Arial) |

### 📊 Table Configuration (CSV/TSV/PSV)
| Flag | Description |
| :--- | :--- |
| `--table-separator` | Force separator (`,` , `;`, `\t`). Auto-detected if empty. |
| `--table-enable-caption` | Add a caption above/below the table. |
| `--table-caption-position` | `top` or `bottom`. |
| `--table-caption-alignment`| `left`, `center`, `right`. |
| `--table-caption-margin` | Space between caption and table (mm). |
| `--table-caption-prefix` | Text before the filename (e.g., "Table 1: "). |
| `--table-stretch-width` | Force table to span the entire page width. |

### 🖼 Image Configuration
| Flag | Description |
| :--- | :--- |
| `--image-enable-caption` | Add filename caption. |
| `--image-caption-position` | `top` or `bottom`. |
| `--image-caption-alignment`| `left`, `center`, `right`. |
| `--image-caption-margin` | Space between image and caption. |
| `--image-caption-prefix` | Text prefix (e.g., "Fig. "). |

### 💻 Code Configuration
| Flag | Description |
| :--- | :--- |
| `--code-style-as-block` | Render with a background color block. |
| `--code-disable-header` | Hide the filename/language header bar. |
| `--code-font-name` | Monospace font (Default: `Courier`). |
| `--code-font-size` | Code font size. |
| `--code-bg-color` | RGB background color (e.g., `245,246,248`). |
| `--code-text-color` | RGB text color (e.g., `36,41,47`). |

### 📝 Text & Markdown Configuration
| Flag | Description |
| :--- | :--- |
| `--text-render-markdown` | Force Markdown parsing for ALL `.txt` files. |

## 💡 Pro-Tips
- **Emoji Stripping**: `fflow` automatically strips emojis from text and code to prevent PDF encoding errors with standard fonts.
- **BOM Handling**: CSV parser automatically detects and skips UTF-8 BOM (Byte Order Mark).
- **Large Files**: Files larger than 10MB are read sequentially during rendering rather than loaded into RAM all at once, preventing OOM (Out of Memory) crashes.

---