package files

import (
	"cmp"
	"slices"
	"sync"
	"unicode"
	"unicode/utf8"
)

var lowerCache sync.Map

func getLowerRune(r rune) rune {
	if r < utf8.RuneSelf {
		if r >= 'A' && r <= 'Z' {
			return r + ('a' - 'A')
		}
		return r
	}

	if cached, ok := lowerCache.Load(r); ok {
		return cached.(rune)
	}

	lower := unicode.ToLower(r)
	lowerCache.Store(r, lower)
	return lower
}

func compareLower(a, b string) int {
	i, j := 0, 0

	for i < len(a) && j < len(b) {
		if a[i] < utf8.RuneSelf && b[j] < utf8.RuneSelf {
			byteA := a[i]
			byteB := b[j]

			if byteA >= 'A' && byteA <= 'Z' {
				byteA += 'a' - 'A'
			}
			if byteB >= 'A' && byteB <= 'Z' {
				byteB += 'a' - 'A'
			}

			if byteA < byteB {
				return -1
			}
			if byteA > byteB {
				return 1
			}

			i++
			j++
			continue
		}

		runeA, sizeA := utf8.DecodeRuneInString(a[i:])
		runeB, sizeB := utf8.DecodeRuneInString(b[j:])

		lowerA := getLowerRune(runeA)
		lowerB := getLowerRune(runeB)

		if lowerA < lowerB {
			return -1
		}
		if lowerA > lowerB {
			return 1
		}

		i += sizeA
		j += sizeB
	}

	remainingA := len(a) - i
	remainingB := len(b) - j

	if remainingA < remainingB {
		return -1
	}
	if remainingA > remainingB {
		return 1
	}

	return 0
}

func (ms *MultiSorter) SortFiles(files []FileInfo) {
	slices.SortFunc(files, func(a, b FileInfo) int {
		for _, c := range ms.Criteria {
			var cmpResult int

			switch c.Field {
			case SortByName:
				cmpResult = compareLower(a.Name, b.Name)
			case SortByModTime:
				cmpResult = cmp.Compare(a.ModTime, b.ModTime)
			case SortBySize:
				cmpResult = cmp.Compare(a.Size, b.Size)
			}

			if cmpResult != 0 {
				if c.Order == Descending {
					return -cmpResult
				}
				return cmpResult
			}
		}
		return 0
	})
}

func (ms *MultiSorter) SortFilesWithDirs(files []FileInfo) {
	slices.SortFunc(files, func(a, b FileInfo) int {
		if a.IsDir != b.IsDir {
			if a.IsDir {
				return -1
			}
			return 1
		}

		for _, c := range ms.Criteria {
			var cmpResult int

			switch c.Field {
			case SortByName:
				cmpResult = compareLower(a.Name, b.Name)
			case SortByModTime:
				cmpResult = cmp.Compare(a.ModTime, b.ModTime)
			case SortBySize:
				cmpResult = cmp.Compare(a.Size, b.Size)
			}

			if cmpResult != 0 {
				if c.Order == Descending {
					return -cmpResult
				}
				return cmpResult
			}
		}
		return 0
	})
}
