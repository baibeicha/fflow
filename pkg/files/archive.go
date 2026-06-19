package files

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// CreateZip creates a zip archive. Context remains unchanged.
func CreateZip(items []FileInfo, destZip string, onProgress func()) ([]FileInfo, *TransferStats, error) {
	var stats TransferStats

	outFile, err := os.Create(destZip)
	if err != nil {
		return nil, nil, err
	}
	defer outFile.Close()

	w := zip.NewWriter(outFile)
	defer w.Close()

	for _, item := range items {
		if item.IsDir {
			if onProgress != nil {
				onProgress()
			}
			continue
		}

		file, err := os.Open(item.Path)
		if err != nil {
			return nil, nil, err
		}

		zipPath := filepath.ToSlash(item.Path)
		f, err := w.Create(zipPath)
		if err != nil {
			file.Close()
			return nil, nil, err
		}

		if _, err = io.Copy(f, file); err != nil {
			file.Close()
			return nil, nil, err
		}
		file.Close()

		stats.Files++
		stats.Bytes += item.Size
		if onProgress != nil {
			onProgress()
		}
	}

	return items, &stats, nil
}
