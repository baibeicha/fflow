package files

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/charlievieth/fastwalk"
)

// CollectFiles traverses a directory and collects information about files.
// If CollectDirs is true, it also gathers directories and computes their accumulated sizes.
// Otherwise, it uses a single-pass algorithm that avoids redundant syscalls.
func CollectFiles(fsc *FolderSearchConfig) ([]FileInfo, error) {
	if fsc.CollectDirs {
		return collectFilesAndDirs(fsc)
	}
	return collectFilesOnly(fsc)
}

// collectFilesOnly uses a single-pass WalkDir algorithm.
func collectFilesOnly(fsc *FolderSearchConfig) ([]FileInfo, error) {
	files := make([]FileInfo, 0, 1024)

	for _, curPath := range fsc.Paths {
		absPath, err := filepath.Abs(curPath)
		if err != nil {
			return nil, err
		}

		info, err := os.Stat(absPath)
		if err != nil {
			return nil, err
		}

		if !info.IsDir() {
			if hasValidExtension(info.Name(), fsc.ValidExtensions) && !fsc.CheckSizeNotFits(info.Size()) {
				files = append(files, FileInfo{
					Path:    absPath,
					Name:    info.Name(),
					ModTime: info.ModTime().Unix(),
					Size:    info.Size(),
					IsDir:   false,
					RelDir:  "",
				})
			}
			continue
		}

		var mu sync.Mutex
		err = fastwalk.Walk(nil, absPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}

			if fsc.IsBlacklisted(path, d.Name()) {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}

			if d.IsDir() {
				if path != absPath && !fsc.Recursively {
					return filepath.SkipDir
				}
				return nil
			}

			if !hasValidExtension(d.Name(), fsc.ValidExtensions) {
				return nil
			}

			info, err := d.Info()
			if err != nil {
				return nil
			}

			if fsc.CheckSizeNotFits(info.Size()) {
				return nil
			}

			dir := parentDir(path)
			relDir := getRelPath(absPath, dir)
			if relDir == "." {
				relDir = ""
			}

			mu.Lock()
			files = append(files, FileInfo{
				Path:    path,
				Name:    d.Name(),
				ModTime: info.ModTime().Unix(),
				Size:    info.Size(),
				IsDir:   false,
				RelDir:  relDir,
			})
			mu.Unlock()
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return files, nil
}

// collectFilesAndDirs collects both files and directories with accumulated sizes.
func collectFilesAndDirs(fsc *FolderSearchConfig) ([]FileInfo, error) {
	var result []FileInfo

	for _, targetPath := range fsc.Paths {
		absPath, err := filepath.Abs(targetPath)
		if err != nil {
			return nil, err
		}

		fi, err := os.Stat(absPath)
		if err != nil || !fi.IsDir() {
			continue
		}

		var allItems []FileInfo
		dirSizes := make(map[string]int64)

		var mu sync.Mutex
		err = fastwalk.Walk(nil, absPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}

			if path == absPath {
				return nil
			}

			if fsc.IsBlacklisted(path, d.Name()) {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}

			info, err := d.Info()
			if err != nil {
				return nil
			}

			isDir := d.IsDir()
			isValidFile := !isDir && hasValidExtension(d.Name(), fsc.ValidExtensions) && !fsc.CheckSizeNotFits(info.Size())

			if isValidFile {
				dir := parentDir(path)
				mu.Lock()
				dirSizes[dir] += info.Size()
				mu.Unlock()
			}

			if isDir || isValidFile {
				var relDir string
				if isDir {
					relDir = getRelPath(absPath, path)
				} else {
					relDir = getRelPath(absPath, parentDir(path))
				}
				if relDir == "." {
					relDir = ""
				}

				mu.Lock()
				allItems = append(allItems, FileInfo{
					Name:    d.Name(),
					Path:    getRelPath(absPath, path),
					Size:    info.Size(),
					ModTime: info.ModTime().Unix(),
					IsDir:   isDir,
					RelDir:  relDir,
				})
				mu.Unlock()
			}

			return nil
		})

		if err != nil {
			return nil, err
		}

		var dirPaths []string
		for _, item := range allItems {
			if item.IsDir {
				absItemPath := filepath.Join(absPath, item.Path)
				if item.Path == "." || item.Path == "" {
					absItemPath = absPath
				}
				dirPaths = append(dirPaths, absItemPath)
			}
		}

		sort.Slice(dirPaths, func(i, j int) bool {
			return len(dirPaths[i]) > len(dirPaths[j])
		})

		for _, dir := range dirPaths {
			parent := parentDir(dir)
			if parent != "." && parent != "/" && parent != "" && parent != dir {
				dirSizes[parent] += dirSizes[dir]
			}
		}

		for _, item := range allItems {
			if item.IsDir {
				absItemPath := filepath.Join(absPath, item.Path)
				if item.Path == "." || item.Path == "" {
					absItemPath = absPath
				}
				item.Size = dirSizes[absItemPath]
				if fsc.CheckSizeNotFits(item.Size) {
					continue
				}
			}
			result = append(result, item)
		}
	}

	return result, nil
}

// hasValidExtension checks if the filename has a valid extension.
func hasValidExtension(filename string, validExts map[string]bool) bool {
	if len(validExts) == 0 {
		return true
	}
	ext := fastExt(filename)
	if ext == "" {
		return false
	}
	return validExts[ext]
}

// fastExt is a zero-alloc, ASCII-optimized alternative to filepath.Ext.
func fastExt(name string) string {
	for i := len(name) - 1; i >= 0; i-- {
		c := name[i]
		if c == '.' {
			return name[i:]
		}
		if c == '/' || c == '\\' {
			break
		}
	}
	return ""
}

// getRelPath is a zero-allocation alternative to filepath.Rel for paths inside absPath.
func getRelPath(absPath, path string) string {
	if path == absPath {
		return "."
	}
	if len(path) > len(absPath) && (path[len(absPath)] == os.PathSeparator || path[len(absPath)] == '/') {
		return path[len(absPath)+1:]
	}
	rel, _ := filepath.Rel(absPath, path)
	return rel
}

// parentDir is a zero-allocation alternative to filepath.Dir.
func parentDir(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == os.PathSeparator || path[i] == '/' {
			if i == 0 {
				return string(os.PathSeparator)
			}
			return path[:i]
		}
	}
	return "."
}
