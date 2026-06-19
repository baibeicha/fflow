[← Back to Main README](../../README.md) | [🇷🇺 Русская версия](flow_RU.md)

# 🔄 Flow Automation Guide

The `flow` command is the heart of `fflow`'s automation capabilities. It allows you to create complex data pipelines where the output of one operation feeds directly into the next.

## 🏗 Architecture
1. **Context Initialization**: `fflow` collects the initial list of files based on your global flags (`-e`, `-r`, etc.).
2. **Variable Interpolation**: Environment variables and custom `var` steps are resolved.
3. **Step Execution**: Each step processes the current `[]FileInfo` array and returns a mutated array for the next step.
4. **Progress Tracking**: A dynamic progress bar updates for every step in the pipeline.

## 📝 Available Actions

| Action | Description | Required Args | Optional Args |
| :--- | :--- | :--- | :--- |
| `copy` | Copies files to destination(s). | `dest` | `rewrite` (true/false) |
| `move` | Moves files to destination(s). | `dest` | `rewrite` |
| `delete` | Deletes files from disk. | *None* | *None* |
| `rename` | Renames files. | *None* | `search`, `replace`, `prefix`, `suffix` |
| `ext` | Changes file extensions. | `to` | *None* |
| `zip` | Creates a zip archive. | `output` | *None* |
| `cmd` | Executes shell command per file. | `exec` | *None* |
| `var` | Defines a pipeline variable. | `name`, `value` | *None* |

## 🧪 Examples

### Example 1: The "Backup & Sanitize" Pipeline (Inline)
Find all `.env` files, copy them to a secure folder, rename them to hide the original name, and zip them.
```bash
fflow flow -e .env -r \
-c "var SECURE=./secure_envs ; copy -d ${SECURE} ; rename --prefix sanitized_ ; ext --to .txt ; zip -o env_backup.zip"
```

### Example 2: Image Processing Pipeline (YAML)
Create a file named `images.yaml`:
```yaml
env:
  OUT_DIR: "./processed_images"

steps:
- action: var
  args:
    name: "PREFIX"
    value: "vacation_"

- action: copy
  args:
    dest: "${OUT_DIR}"

- action: rename
  args:
    prefix: "${PREFIX}"

- action: cmd
  args:
    # Use {} as a placeholder for the file path
    exec: "exiftool -all= {}"
```
Run it:
```bash
fflow flow -e .jpg,.png -f images.yaml
```

## ⚠️ Important Notes
- **Destructive Actions**: The `delete` and `move` actions permanently alter your filesystem. Always test pipelines with `copy` first.
- **Variable Syntax**: Use `${VAR_NAME}` to reference variables in YAML or inline strings.
- **System Env**: Pass the `--env` flag to inject your OS environment variables (e.g., `$HOME`, `$USER`) into the pipeline context.

---
