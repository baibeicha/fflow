[← Back to Main README](../../README.md) | [🇷🇺 Русская версия](copy-move_RU.md)

# 📦 File Transfer (`copy` & `move`)

The transfer engine is designed for safety, concurrency, and multi-destination support. It uses an `ants` goroutine pool sized to your CPU cores to maximize throughput without thrashing the disk.

## 🛡 Collision Resolution
By default, `fflow` **will never overwrite** an existing file unless you explicitly pass the `-w` (`--rewrite`) flag. 
If a collision occurs, `fflow` automatically generates a safe unique name:
1. `file.txt` ➔ `file Copy.txt`
2. `file Copy.txt` ➔ `file Copy (1).txt`
3. `file Copy (1).txt` ➔ `file Copy (2).txt`

## 🚚 `copy`: Multi-Destination Copying
You can copy files to multiple directories simultaneously.
```bash
fflow copy -e .pdf -d ./backup_drive_1 -d ./backup_drive_2 -d ./cloud_sync
```

## 📦 `move`: Safe Multi-Destination Moving
Moving to multiple destinations is inherently risky. `fflow` handles this safely:
1. It **copies** the file to all specified destinations.
2. It tracks the success rate of every destination.
3. **Only if the copy succeeds in ALL destinations**, it deletes the original source file.
```bash
fflow move -e .mp4 -d ./archive_1 -d ./archive_2 -w
```

## 🚩 Flags
| Flag | Short | Description |
| :--- | :---: | :--- |
| `--dest` | `-d` | Destination directory (can be used multiple times). |
| `--rewrite` | `-w` | Force overwrite existing files. |

---
