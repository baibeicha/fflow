[← Back to Main README](../../README.md) | [🇷🇺 Русская версия](info_RU.md)

# 📂 Directory Analysis (`info`)

The `info` command is a high-performance directory analyzer. Unlike standard tools that require multiple passes to calculate folder sizes, `fflow` uses a **single-pass bottom-up accumulation algorithm**.

## 🧠 How it Works
1. **Single Pass Traversal**: Uses `fastwalk` to read every file and directory exactly once.
2. **Size Accumulation**: As it finds files, it adds their size to their immediate parent directory in a memory map.
3. **Bottom-Up Propagation**: It sorts directories by path depth (deepest first) and propagates sizes up the tree. This guarantees 100% accurate nested directory sizes without redundant `stat()` syscalls.

## 💻 Usage Examples

**Analyze current directory (non-recursive):**
```bash
fflow info .
```

**Find the heaviest directories in a project recursively:**
```bash
fflow info ./project -r --sort-by size:desc
```

**Analyze only specific file types (e.g., media files):**
```bash
fflow info ./media -r -e .mp4,.mkv,.avi --min-size 1gb
```

## 📊 Output Format
`info` outputs a beautifully formatted table using `tabwriter`, displaying:
- **Name** (Directories are suffixed with `/`)
- **Size** (Human-readable, e.g., `1.4 GB`)
- **Modified Time** (`YYYY-MM-DD HH:MM:SS`)
- **Relative Path**

---
