[← Back to Main README](../../README.md) | [🇷🇺 Русская версия](stats-merge_RU.md)

# 📊 Text Analytics & Merging (`stats` & `merge`)

These commands are built around a highly optimized `CountingWriter` and `bufio.Scanner` that accurately handles UTF-8 Unicode characters, ensuring exact word and character counts across any language.

## 📈 `stats`: Deep Text Analytics

Calculates lines, words, characters, and characters (excluding spaces) across thousands of files concurrently using an `ants` worker pool.

**Flags:**
- `--count-lines` (Default: `true`)
- `--count-words` (Default: `true`)
- `--count-chars` (Default: `false`)
- `--count-chars-no-space` (Default: `false`)

**Example:**
```bash
# Count lines and words in all Go files, excluding vendor folders
fflow stats -r -e .go --blacklist vendor
```

## 🔗 `merge`: File Concatenation

Combines multiple files into a single output file. It supports two distinct modes:

### 1. Full Mode (Default)
Reads files through the `CountingWriter`. This allows `fflow` to output exact statistics (lines, words, chars) about the merged result while writing it to disk.
```bash
fflow merge -e .md -o book.md --include-path --separator "\n\n---\n\n"
```

### 2. Fast Mode (`-f` / `--fast`)
Bypasses the counting logic and uses `io.CopyBuffer` with a `sync.Pool` of 32KB byte buffers. This achieves **maximum possible disk I/O throughput**, ideal for merging gigabytes of log files where you don't need text statistics.
```bash
fflow merge -e .log -o all_logs.txt -f --include-name
```

**Header Flags:**
- `--include-path`: Adds `==== /absolute/path/to/file.txt ====` before each file.
- `--include-name`: Adds `==== filename.txt ====` before each file.

---
