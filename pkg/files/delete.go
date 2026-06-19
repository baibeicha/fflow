package files

import (
	"os"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/panjf2000/ants/v2"
)

// DeleteFiles deletes provided files and clears them from the flow context.
func DeleteFiles(items []FileInfo, onProgress func()) ([]FileInfo, *TransferStats, error) {
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

			if err := os.Remove(item.Path); err != nil {
				errorOnce.Do(func() { firstError = err })
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

	return nil, &stats, nil
}
