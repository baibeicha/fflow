package pdf

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/baibeicha/fflow/pkg/files"
)

type FileCategory byte

const (
	CategoryUnknown FileCategory = iota
	CategoryImage
	CategoryText
	CategoryTable
	CategoryCode
	CategoryAll
)

var ExtensionMap = map[string]FileCategory{
	".jpg":  CategoryImage,
	".jpeg": CategoryImage,
	".png":  CategoryImage,
	".gif":  CategoryImage,

	".csv": CategoryTable,
	".tsv": CategoryTable,
	".psv": CategoryTable,

	".txt": CategoryText,
	".md":  CategoryText,
	".log": CategoryText,
	".rtf": CategoryText,

	".json":       CategoryText,
	".xml":        CategoryText,
	".yaml":       CategoryText,
	".yml":        CategoryText,
	".toml":       CategoryText,
	".ini":        CategoryText,
	".env":        CategoryText,
	".conf":       CategoryText,
	".cfg":        CategoryText,
	".properties": CategoryText,

	".go":           CategoryCode,
	".mod":          CategoryCode,
	".sum":          CategoryCode,
	".proto":        CategoryCode,
	".sql":          CategoryCode,
	".graphql":      CategoryCode,
	".gql":          CategoryCode,
	".htm":          CategoryCode,
	".html":         CategoryCode,
	".css":          CategoryCode,
	".scss":         CategoryCode,
	".less":         CategoryCode,
	".js":           CategoryCode,
	".ts":           CategoryCode,
	".jsx":          CategoryCode,
	".tsx":          CategoryCode,
	".vue":          CategoryCode,
	".svg":          CategoryCode,
	".sh":           CategoryCode,
	".bash":         CategoryCode,
	".zsh":          CategoryCode,
	".bat":          CategoryCode,
	".cmd":          CategoryCode,
	".ps1":          CategoryCode,
	".py":           CategoryCode,
	".rb":           CategoryCode,
	".php":          CategoryCode,
	".java":         CategoryCode,
	".c":            CategoryCode,
	".cpp":          CategoryCode,
	".h":            CategoryCode,
	".hpp":          CategoryCode,
	".cs":           CategoryCode,
	".rs":           CategoryCode,
	".swift":        CategoryCode,
	".kt":           CategoryCode,
	".dart":         CategoryCode,
	".lua":          CategoryCode,
	".pl":           CategoryCode,
	".m":            CategoryCode,
	".r":            CategoryCode,
	".dockerfile":   CategoryCode,
	".makefile":     CategoryCode,
	".mk":           CategoryCode,
	".gitignore":    CategoryCode,
	".dockerignore": CategoryCode,
}

var allExtensions []string
var categoryExtensions map[FileCategory][]string

func init() {
	allExtensions = make([]string, 0, len(ExtensionMap))
	categoryExtensions = make(map[FileCategory][]string)

	for ext, cat := range ExtensionMap {
		allExtensions = append(allExtensions, ext)
		categoryExtensions[cat] = append(categoryExtensions[cat], ext)
	}
}

func GetCategory(fileName string) FileCategory {
	ext := strings.ToLower(filepath.Ext(fileName))
	if ext == "" {
		ext = "." + strings.ToLower(filepath.Base(fileName))
	}

	if cat, ok := ExtensionMap[ext]; ok {
		return cat
	}

	if filepath.Ext(fileName) == "" {
		return CategoryText
	}
	return CategoryUnknown
}

func ExportTypes(fsc *files.FolderSearchConfig, categories ...FileCategory) {
	if len(categories) == 0 {
		return
	}

	if slices.Contains(categories, CategoryAll) {
		fsc.SearchForExtensions(allExtensions...)
		return
	}

	var targetExts []string
	for _, cat := range categories {
		if exts, ok := categoryExtensions[cat]; ok {
			targetExts = append(targetExts, exts...)
		}
	}

	if len(targetExts) > 0 {
		fsc.SearchForExtensions(targetExts...)
	}
}
