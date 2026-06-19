[← Back to Main README](../../README.md) | [🇷🇺 Русская версия](delete_RU.md)

# 🗑 Safe Deletion (`delete`)

The `delete` command permanently removes files matching your filters. It runs concurrently to process large batches instantly.

⚠️ **WARNING**: This action is irreversible. Always test your filters using the `info` command first to ensure you are targeting the correct files!

## 💻 Usage Examples

**Delete all temporary build files:**
```bash
fflow delete -r -e .tmp,.bak,.log
```

**Delete files larger than 5GB in the Downloads folder:**
```bash
fflow delete ~/Downloads --min-size 5gb
```

**Dry-Run Alternative:**
Since `fflow` doesn't have a built-in `--dry-run` flag for delete, use `info` to preview:
```bash
fflow info -r -e .tmp,.bak,.log
```

---
