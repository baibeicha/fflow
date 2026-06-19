package files

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/panjf2000/ants/v2"
)

// ExecuteCommandOnFiles executes a shell command per file. Context remains unchanged.
func ExecuteCommandOnFiles(items []FileInfo, cmdTemplate string, onProgress func()) ([]FileInfo, *TransferStats, error) {
	var stats TransferStats

	poolSize := runtime.NumCPU()
	p, err := ants.NewPool(poolSize, ants.WithPreAlloc(true))
	if err != nil {
		return nil, nil, err
	}
	defer p.Release()

	var wg sync.WaitGroup
	var firstError error
	var errorOnce sync.Once

	for _, item := range items {
		if item.IsDir {
			if onProgress != nil {
				onProgress()
			}
			continue
		}

		item := item
		wg.Add(1)

		p.Submit(func() {
			defer wg.Done()

			if firstError != nil {
				return
			}

			absPath, _ := filepath.Abs(item.Path)
			execStr := strings.ReplaceAll(cmdTemplate, "{}", fmt.Sprintf(`"%s"`, absPath))

			var cmd *exec.Cmd
			if runtime.GOOS == "windows" {
				cmd = exec.Command("cmd", "/c", execStr)
			} else {
				cmd = exec.Command("sh", "-c", execStr)
			}

			if err := cmd.Run(); err != nil {
				errorOnce.Do(func() {
					firstError = fmt.Errorf("error in file %s: %w", item.Name, err)
				})
				return
			}

			atomic.AddInt64(&stats.Files, 1)
			atomic.AddInt64(&stats.Bytes, item.Size)

			if onProgress != nil {
				onProgress()
			}
		})
	}

	wg.Wait()

	if firstError != nil {
		return nil, nil, firstError
	}

	return items, &stats, nil
}
