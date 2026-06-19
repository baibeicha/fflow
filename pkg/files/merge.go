package files

import (
	"bufio"
	"fflow/pkg/files/countingwriter"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"
)

// MergeFiles merges the given slice of files and returns statistics
func MergeFiles(filesList []FileInfo, outputPath string, cfg *MergeConfig, fast bool, onProgress func()) (*FileStats, error) {
	out, err := os.Create(outputPath)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	writer := bufio.NewWriterSize(out, 256*1024)
	defer writer.Flush()

	var separator string
	if cfg.Separator != "" {
		separator = strings.ReplaceAll(cfg.Separator, "\\n", "\n")
	}

	if fast {
		return mergeFilesFast(filesList, writer, cfg, separator, onProgress)
	}
	return mergeFilesWithStats(filesList, writer, cfg, separator, onProgress)
}

var copyBufPool = sync.Pool{
	New: func() any {
		return new(make([]byte, 32*1024))
	},
}

func mergeFilesFast(filesList []FileInfo, writer *bufio.Writer, cfg *MergeConfig, separator string, onProgress func()) (*FileStats, error) {
	var stats FileStats

	for i, f := range filesList {
		stats.Files++

		if cfg.IncludeFilePath || cfg.IncludeFileName {
			var name string
			if cfg.IncludeFilePath {
				name = f.Path
			} else {
				name = f.Name
			}

			header := "==== " + name + " ====\n"
			if _, err := writer.WriteString(header); err != nil {
				return &stats, err
			}
			stats.Bytes += int64(len(header))
		}

		inFile, err := os.Open(f.Path)
		if err != nil {
			errMsg := "[Error reading file: " + err.Error() + "]\n"
			slog.Error(errMsg)
			if _, err := writer.WriteString(errMsg); err != nil {
				return &stats, err
			}
			stats.Bytes += int64(len(errMsg))

			if onProgress != nil {
				onProgress()
			}
			continue
		}

		bufPtr := copyBufPool.Get().(*[]byte)
		written, err := io.CopyBuffer(writer, inFile, *bufPtr)
		inFile.Close()
		copyBufPool.Put(bufPtr)

		if err != nil {
			return &stats, err
		}

		stats.Bytes += written

		if _, err := writer.WriteString("\n"); err != nil {
			return &stats, err
		}
		stats.Bytes++

		if i < len(filesList)-1 && separator != "" {
			if _, err := writer.WriteString(separator); err != nil {
				return &stats, err
			}
			stats.Bytes += int64(len(separator))
		}

		if onProgress != nil {
			onProgress()
		}
	}

	return &stats, nil
}

func mergeFilesWithStats(filesList []FileInfo, writer *bufio.Writer, cfg *MergeConfig, separator string, onProgress func()) (*FileStats, error) {
	var stats FileStats

	var opts []countingwriter.Option
	if cfg.CountLines {
		opts = append(opts, countingwriter.WithLines())
	}
	if cfg.CountWords {
		opts = append(opts, countingwriter.WithWords())
	}
	if cfg.CountChars {
		opts = append(opts, countingwriter.WithCharacters())
	}
	if cfg.CountCharsNoSpc {
		opts = append(opts, countingwriter.WithCharactersNoSpace())
	}

	cw := countingwriter.New(writer, opts...)

	for i, f := range filesList {
		stats.Files++

		if cfg.IncludeFilePath || cfg.IncludeFileName {
			var name string
			if cfg.IncludeFilePath {
				name = f.Path
			} else {
				name = f.Name
			}

			header := "==== " + name + " ====\n"
			if _, err := cw.Write([]byte(header)); err != nil {
				return &stats, err
			}
		}

		inFile, err := os.Open(f.Path)
		if err != nil {
			errMsg := "[Error reading file: " + err.Error() + "]\n"
			slog.Error(errMsg)
			if _, err := cw.Write([]byte(errMsg)); err != nil {
				return &stats, err
			}

			if onProgress != nil {
				onProgress()
			}
			continue
		}

		bufPtr := copyBufPool.Get().(*[]byte)
		_, err = io.CopyBuffer(cw, inFile, *bufPtr)
		inFile.Close()
		copyBufPool.Put(bufPtr)

		if _, err := cw.Write([]byte{'\n'}); err != nil {
			return &stats, err
		}

		if i < len(filesList)-1 && separator != "" {
			if _, err := cw.Write([]byte(separator)); err != nil {
				return &stats, err
			}
		}

		if onProgress != nil {
			onProgress()
		}
	}

	countingStats := cw.Stats()
	stats.Lines = countingStats.Lines
	stats.Words = countingStats.Words
	stats.Characters = countingStats.Characters
	stats.CharactersNoSpace = countingStats.CharactersNoSpc
	stats.Bytes = countingStats.Bytes

	return &stats, nil
}
