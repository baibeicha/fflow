package files

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/panjf2000/ants/v2"
)

var dirCache sync.Map

func makeDirAll(dir, relDir string) (string, error) {
	targetDir := filepath.Join(dir, relDir)
	if _, loaded := dirCache.LoadOrStore(targetDir, struct{}{}); !loaded {
		if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
			return targetDir, err
		}
	}
	return targetDir, nil
}

// TransferFiles copies or moves files to destination directories.
func TransferFiles(items []FileInfo, destDirs []string, isMove bool, rewrite bool, copySuffix string, onProgress func()) ([]FileInfo, *TransferStats, error) {
	var stats TransferStats

	newFiles := make([]FileInfo, len(items))

	for _, dir := range destDirs {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, nil, fmt.Errorf("failed to create destination directory %s: %w", dir, err)
		}
	}

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

			defer func() {
				if onProgress != nil {
					onProgress()
				}
			}()

			src := item.Path
			filename := item.Name
			size := item.Size
			relDir := item.RelDir

			var mutated FileInfo

			if isMove && len(destDirs) == 1 {
				targetDir, err := makeDirAll(destDirs[0], relDir)
				if err != nil {
					errorOnce.Do(func() { firstError = err })
					return
				}

				dst := getUniqueDestPath(targetDir, filename, copySuffix, rewrite)
				if absSrc, _ := filepath.Abs(src); absSrc != getAbs(dst) {
					if err := moveFile(src, dst); err != nil {
						errorOnce.Do(func() { firstError = err })
						return
					}
					atomic.AddInt64(&stats.Files, 1)
					atomic.AddInt64(&stats.Bytes, size)
				}

				mutated = item
				mutated.Path = dst
				mutated.Name = filepath.Base(dst)
			} else {
				successCount := 0
				var firstValidDst string

				for _, dir := range destDirs {
					targetDir, err := makeDirAll(dir, relDir)
					if err != nil {
						errorOnce.Do(func() { firstError = err })
						return
					}

					dst := getUniqueDestPath(targetDir, filename, copySuffix, rewrite)
					if absSrc, _ := filepath.Abs(src); absSrc == getAbs(dst) {
						successCount++
						continue
					}

					if err := copyFile(src, dst); err != nil {
						errorOnce.Do(func() { firstError = err })
						return
					}

					if firstValidDst == "" {
						firstValidDst = dst
					}

					atomic.AddInt64(&stats.Files, 1)
					atomic.AddInt64(&stats.Bytes, size)
					successCount++
				}

				if isMove && successCount == len(destDirs) {
					if err := os.Remove(src); err != nil {
						errorOnce.Do(func() { firstError = err })
						return
					}
				}

				if isMove && firstValidDst != "" {
					mutated = item
					mutated.Path = firstValidDst
					mutated.Name = filepath.Base(firstValidDst)
				} else {
					mutated = item
				}
			}

			newFiles[idx] = mutated
		})
	}

	wg.Wait()

	if firstError != nil {
		return nil, nil, firstError
	}

	return newFiles, &stats, nil
}

func getAbs(path string) string {
	abs, _ := filepath.Abs(path)
	return abs
}

func getUniqueDestPath(destDir, filename, suffix string, rewrite bool) string {
	target := filepath.Join(destDir, filename)

	if rewrite {
		return target
	}

	if _, err := os.Stat(target); os.IsNotExist(err) {
		return target
	}

	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, ext)

	target = filepath.Join(destDir, fmt.Sprintf("%s %s%s", base, suffix, ext))
	if _, err := os.Stat(target); os.IsNotExist(err) {
		return target
	}

	counter := 1
	for {
		target = filepath.Join(destDir, fmt.Sprintf("%s %s (%d)%s", base, suffix, counter, ext))
		if _, err := os.Stat(target); os.IsNotExist(err) {
			return target
		}
		counter++
	}
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	if info, err := os.Stat(src); err == nil {
		_ = os.Chtimes(dst, info.ModTime(), info.ModTime())
	}

	return out.Sync()
}

func moveFile(src, dst string) error {
	err := os.Rename(src, dst)
	if err != nil {
		if copyErr := copyFile(src, dst); copyErr != nil {
			return fmt.Errorf("move failed: %v", copyErr)
		}
		return os.Remove(src)
	}
	return nil
}
