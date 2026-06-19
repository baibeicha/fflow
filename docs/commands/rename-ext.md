[← Back to Main README](../../README.md) | [🇷🇺 Русская версия](rename-ext_RU.md)

# ✏️ Bulk Renaming (`rename` & `ext`)

These commands modify file names and extensions. Like the transfer engine, they include built-in collision avoidance to prevent accidental overwrites during the rename process.

## 🔄 `rename`: Search, Replace, Prefix & Suffix

Modifies the base name of the file (ignoring the extension).

**Flags:**
- `--search`: String to find in the filename.
- `--replace`: String to replace the search term with.
- `--prefix`: String to prepend to the filename.
- `--suffix`: String to append to the filename (before the extension).

**Example: Sanitize filenames**
```bash
# Replaces spaces with underscores and adds a date prefix
fflow rename -r -e .jpg --search " " --replace "_" --prefix "2026_"
```

## 📎 `ext`: Extension Changer

Changes the file extension in bulk. You do not need to include the dot (`.`); `fflow` will add it automatically if missing.

**Example: Convert all CSVs to TSVs**
```bash
fflow ext -r -e .csv --to tsv
```

---
