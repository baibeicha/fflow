package files

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/panjf2000/ants/v2"
)

// RenameFiles renames files and returns their updated metadata.
func RenameFiles(items []FileInfo, search, replace, prefix, suffix string, onProgress func()) ([]FileInfo, *TransferStats, error) {
	var stats TransferStats

	newFiles := make([]FileInfo, len(items))

	poolSize := runtime.NumCPU()
	p, err := ants.NewPool(poolSize, ants.WithPreAlloc(true))
	if err != nil {
		return nil, nil, err
	}
	defer p.Release()

	var wg sync.WaitGroup
	var firstError error
	var errorOnce sync.Once

	for idx, item := range items {
		if item.IsDir {
			if onProgress != nil {
				onProgress()
			}
			continue
		}

		idx := idx
		item := item

		wg.Add(1)
		p.Submit(func() {
			defer wg.Done()

			if firstError != nil {
				return
			}

			dir := filepath.Dir(item.Path)
			ext := filepath.Ext(item.Name)
			base := strings.TrimSuffix(item.Name, ext)

			if search != "" {
				base = strings.ReplaceAll(base, search, replace)
			}

			newName := prefix + base + suffix + ext
			newPath := filepath.Join(dir, newName)

			mutated := item
			if item.Path != newPath {
				safePath := getUniqueDestPath(dir, newName, "", false)
				if err := os.Rename(item.Path, safePath); err != nil {
					errorOnce.Do(func() { firstError = err })
					return
				}
				mutated.Path = safePath
				mutated.Name = filepath.Base(safePath)
			}

			newFiles[idx] = mutated

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

	return newFiles, &stats, nil
}

// ChangeExtension changes the extension of files and returns updated metadata.
func ChangeExtension(items []FileInfo, newExt string, onProgress func()) ([]FileInfo, *TransferStats, error) {
	if newExt != "" && !strings.HasPrefix(newExt, ".") {
		newExt = "." + newExt
	}

	var stats TransferStats
	newFiles := make([]FileInfo, len(items))

	poolSize := runtime.NumCPU()
	p, err := ants.NewPool(poolSize, ants.WithPreAlloc(true))
	if err != nil {
		return nil, nil, err
	}
	defer p.Release()

	var wg sync.WaitGroup
	var firstError error
	var errorOnce sync.Once

	for idx, item := range items {
		if item.IsDir {
			if onProgress != nil {
				onProgress()
			}
			continue
		}

		idx := idx
		item := item

		wg.Add(1)
		p.Submit(func() {
			defer wg.Done()

			if firstError != nil {
				return
			}

			dir := filepath.Dir(item.Path)
			oldExt := filepath.Ext(item.Name)
			base := strings.TrimSuffix(item.Name, oldExt)
			newName := base + newExt
			newPath := filepath.Join(dir, newName)

			mutated := item
			if item.Path != newPath {
				safePath := getUniqueDestPath(dir, newName, "", false)
				if err := os.Rename(item.Path, safePath); err != nil {
					errorOnce.Do(func() { firstError = err })
					return
				}
				mutated.Path = safePath
				mutated.Name = filepath.Base(safePath)
			}

			newFiles[idx] = mutated

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

	return newFiles, &stats, nil
}
