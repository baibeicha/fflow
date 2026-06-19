[← Back to Main README](../README.md) | [🇷🇺 Русская версия](global-flags_RU.md)

# 🌐 Global Flags, Filtering & Sorting

Almost every command in `fflow` (`info`, `stats`, `merge`, `copy`, `move`, `delete`, `rename`, `ext`, `zip`, `cmd`, `pdf`, `flow`) shares a powerful underlying file discovery engine: the `FolderSearchConfig`. 

This engine uses the `fastwalk` library for blazing-fast, concurrent directory traversal and zero-allocation string parsing for extensions.

## 🎯 Targeting Files

| Flag | Short | Description | Deep Dive |
| :--- | :---: | :--- | :--- |
| `--path` | `-p` | Target directories. | Accepts multiple paths. You can use `-p ./src -p ./docs` or comma-separated `-p ./src,./docs`. |
| `--recursive` | `-r` | Traverse subdirectories. | If omitted, `fflow` only processes the immediate children of the target path(s). |
| `--extensions` | `-e` | Filter by file extension. | E.g., `-e .go,.md`. Internally uses a zero-allocation ASCII-optimized `fastExt()` parser. |
| `--blacklist` | | Exclude files or directories. | If a directory matches a blacklist item, `fastwalk` immediately skips it (`filepath.SkipDir`), saving massive amounts of I/O. |

## 📏 Size Filtering

You can filter files by exact byte size using human-readable units. `fflow` supports: `b`, `kb`, `mb`, `gb`, `tb`, `pb`.

| Flag | Description | Example |
| :--- | :--- | :--- |
| `--min-size` | Ignore files smaller than this. | `--min-size 500kb` |
| `--max-size` | Ignore files larger than this. | `--max-size 2gb` |

*Note: Size filters apply to individual files. When using the `info` command to calculate directory sizes, directories that only contain filtered-out files will correctly report a size of 0.*

## 🔄 Advanced Multi-Sorting

The `--sort-by` flag allows you to chain multiple sorting criteria. The syntax is `field:order`.

**Fields:**
- `name` (Alphabetical, case-insensitive, Unicode-aware)
- `size` (Bytes)
- `modtime` (Last modified timestamp)

**Orders:**
- `asc` (Ascending - Default)
- `desc` (Descending)

**Examples:**
```bash
# Sort by size descending, then by name ascending for files of the same size
fflow stats --sort-by size:desc,name:asc

# Sort by newest first
fflow info --sort-by modtime:desc
```
*Under the hood: `fflow` uses a custom `MultiSorter` with a `sync.Map` cache for Unicode lowercase conversion, ensuring sorting millions of files remains extremely fast.*

---
