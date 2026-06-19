package files

import (
	"bufio"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"unicode"
	"unicode/utf8"

	"github.com/panjf2000/ants/v2"
)

// CalculateStats iterates through the passed list of files and calculates the selected metrics
func CalculateStats(filesList []FileInfo, cfg *StatsConfig, onProgress func()) (*FileStats, error) {
	var totalStats FileStats

	poolSize := runtime.NumCPU()
	p, err := ants.NewPool(poolSize, ants.WithPreAlloc(true))
	if err != nil {
		return nil, err
	}
	defer p.Release()

	var wg sync.WaitGroup
	var firstError error
	var errorOnce sync.Once

	for _, f := range filesList {
		fileInfo := f
		wg.Add(1)

		p.Submit(func() {
			defer wg.Done()

			if firstError != nil {
				return
			}

			if onProgress != nil {
				defer onProgress()
			}

			stats, err := calculateFileStats(fileInfo, cfg)
			if err != nil {
				errorOnce.Do(func() { firstError = err })
				return
			}

			atomic.AddInt64(&totalStats.Files, stats.Files)
			atomic.AddInt64(&totalStats.Lines, stats.Lines)
			atomic.AddInt64(&totalStats.Words, stats.Words)
			atomic.AddInt64(&totalStats.Characters, stats.Characters)
			atomic.AddInt64(&totalStats.CharactersNoSpace, stats.CharactersNoSpace)
			atomic.AddInt64(&totalStats.Bytes, stats.Bytes)
		})
	}

	wg.Wait()

	if firstError != nil {
		return &totalStats, firstError
	}

	return &totalStats, nil
}

var bufferPool = sync.Pool{
	New: func() any {
		return new(make([]byte, 64*1024))
	},
}

func calculateFileStats(f FileInfo, cfg *StatsConfig) (FileStats, error) {
	var stats FileStats
	stats.Files = 1

	file, err := os.Open(f.Path)
	if err != nil {
		return stats, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return stats, err
	}
	stats.Bytes = info.Size()

	bufPtr := bufferPool.Get().(*[]byte)
	buf := *bufPtr
	defer func() {
		buf = buf[:0]
		*bufPtr = buf
		bufferPool.Put(bufPtr)
	}()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		if cfg.CountLines {
			stats.Lines++
		}

		lineBytes := scanner.Bytes()

		if cfg.CountWords || cfg.CountChars || cfg.CountCharsNoSpc {
			countLineStats(lineBytes, cfg, &stats)
		}
	}

	if err := scanner.Err(); err != nil {
		return stats, err
	}

	return stats, nil
}

func countLineStats(line []byte, cfg *StatsConfig, stats *FileStats) {
	var (
		words        int64
		chars        int64
		charsNoSpace int64
		inWord       bool
	)

	needChars := cfg.CountChars || cfg.CountCharsNoSpc
	needWords := cfg.CountWords

	if !needWords && !needChars {
		return
	}

	for i := 0; i < len(line); {
		b := line[i]

		if b < utf8.RuneSelf {
			i++

			if needWords {
				if b == ' ' || b == '\t' || b == '\n' || b == '\r' {
					inWord = false
				} else if !inWord {
					inWord = true
					words++
				}
			}

			if needChars {
				chars++
				if cfg.CountCharsNoSpc && b != ' ' && b != '\t' {
					charsNoSpace++
				}
			}

			continue
		}

		r, size := utf8.DecodeRune(line[i:])
		i += size

		if needWords {
			if unicode.IsSpace(r) {
				inWord = false
			} else if !inWord {
				inWord = true
				words++
			}
		}

		if needChars {
			chars++
			if cfg.CountCharsNoSpc && !unicode.IsSpace(r) {
				charsNoSpace++
			}
		}
	}

	if needWords {
		stats.Words += words
	}
	if cfg.CountChars {
		stats.Characters += chars + 1
	}
	if cfg.CountCharsNoSpc {
		stats.CharactersNoSpace += charsNoSpace
	}
}
