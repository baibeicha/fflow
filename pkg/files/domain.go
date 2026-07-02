package files

import "strings"

// FileInfo stores file data for easy sorting, filtering, and pipeline processing.
type FileInfo struct {
	Path    string
	Name    string
	ModTime int64
	Size    int64
	IsDir   bool
	RelDir  string
}

// FolderSearchConfig stores rules on how to search for files.
type FolderSearchConfig struct {
	Paths           []string
	Recursively     bool
	MinSize         *int64
	MaxSize         *int64
	ValidExtensions map[string]bool
	BlackList       map[string]struct{}
	CollectDirs     bool
	Limit           *uint64
	Offset          *uint64
}

// FlowContext holds the state of the pipeline execution and mutates between steps.
type FlowContext struct {
	Files []FileInfo
	Vars  map[string]string
	Stats TransferStats
}

// PipelineStep represents a single operation in the flow.
type PipelineStep struct {
	Action string            `yaml:"action"`
	Args   map[string]string `yaml:"args"`
}

// Pipeline represents the full flow execution plan.
type Pipeline struct {
	Env   map[string]string `yaml:"env"`
	Steps []PipelineStep    `yaml:"steps"`
}

// TransferStats stores statistics for file transfer operations.
type TransferStats struct {
	Files int64
	Bytes int64
}

// Add combines two TransferStats.
func (s *TransferStats) Add(other *TransferStats) {
	s.Files += other.Files
	s.Bytes += other.Bytes
}

func NewFolderSearchConfig(recursively bool, paths ...string) *FolderSearchConfig {
	return &FolderSearchConfig{
		Paths:           paths,
		Recursively:     recursively,
		ValidExtensions: make(map[string]bool),
		MinSize:         nil,
		MaxSize:         nil,
		BlackList:       make(map[string]struct{}),
		CollectDirs:     false,
	}
}

func (fs *FolderSearchConfig) AddToBlackList(paths ...string) {
	for _, p := range paths {
		fs.BlackList[p] = struct{}{}
	}
}

func (fs *FolderSearchConfig) IsBlacklisted(path string, name string) bool {
	if len(fs.BlackList) == 0 {
		return false
	}

	if _, ok := fs.BlackList[name]; ok {
		return true
	}

	for bad := range fs.BlackList {
		if strings.Contains(path, bad) {
			return true
		}
	}
	return false
}

func (fs *FolderSearchConfig) SearchForExtensions(extensions ...string) {
	for _, ext := range extensions {
		if !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}
		fs.ValidExtensions[ext] = true
	}
}

func (fs *FolderSearchConfig) SetMinSize(size int64) {
	fs.MinSize = &size
}

func (fs *FolderSearchConfig) SetMaxSize(size int64) {
	fs.MaxSize = &size
}

func (fs *FolderSearchConfig) SetLimit(limit uint64) {
	fs.Limit = &limit
}

func (fs *FolderSearchConfig) SetOffset(offset uint64) {
	fs.Offset = &offset
}

func (fs *FolderSearchConfig) CheckSizeNotFits(size int64) bool {
	if fs.MinSize != nil && size < *fs.MinSize {
		return true
	}
	if fs.MaxSize != nil && size > *fs.MaxSize {
		return true
	}
	return false
}

func (fs *FolderSearchConfig) SetSize(minSize int64, maxSize int64) {
	fs.MinSize = &minSize
	fs.MaxSize = &maxSize
}

func (fs *FolderSearchConfig) SetRecursively(recursively bool) {
	fs.Recursively = recursively
}

func SizeFromUnit(size int64, unit string) int64 {
	switch unit {
	case "", "b":
		return size
	case "kb":
		return size * 1024
	case "mb":
		return size * 1024 * 1024
	case "gb":
		return size * 1024 * 1024 * 1024
	case "tb":
		return size * 1024 * 1024 * 1024 * 1024
	case "pb":
		return size * 1024 * 1024 * 1024 * 1024 * 1024
	default:
		return size
	}
}

type SortField byte

const (
	SortByName SortField = iota
	SortByModTime
	SortBySize
)

type SortOrder byte

const (
	Ascending SortOrder = iota
	Descending
)

// SortCriteria is a single sort condition
type SortCriteria struct {
	Field SortField
	Order SortOrder
}

// MultiSorter stores a chain of sorting criteria
type MultiSorter struct {
	Criteria []SortCriteria
}

func NewSortCriteria(field SortField, order SortOrder) *SortCriteria {
	return &SortCriteria{
		Field: field,
		Order: order,
	}
}

func NewMultiSorter(criteria ...SortCriteria) *MultiSorter {
	return &MultiSorter{
		Criteria: criteria,
	}
}

// StatsConfig defines which metrics should be collected
type StatsConfig struct {
	CountLines      bool
	CountWords      bool
	CountChars      bool
	CountCharsNoSpc bool
}

// DefaultStatsConfig includes all basic counters
var DefaultStatsConfig = &StatsConfig{
	CountLines: true,
	CountWords: true,
	CountChars: true,
}

// FileStats stores counting results for a single file or group of files.
type FileStats struct {
	Files             int64
	Lines             int64
	Words             int64
	Characters        int64
	CharactersNoSpace int64
	Bytes             int64
}

func (s *FileStats) Add(other *FileStats) {
	s.Files += other.Files
	s.Lines += other.Lines
	s.Words += other.Words
	s.Characters += other.Characters
	s.CharactersNoSpace += other.CharactersNoSpace
	s.Bytes += other.Bytes
}

func (s *FileStats) Reset() {
	s.Files = 0
	s.Lines = 0
	s.Words = 0
	s.Characters = 0
	s.CharactersNoSpace = 0
	s.Bytes = 0
}

func (s *FileStats) Clone() *FileStats {
	return &FileStats{
		Files:             s.Files,
		Lines:             s.Lines,
		Words:             s.Words,
		Characters:        s.Characters,
		CharactersNoSpace: s.CharactersNoSpace,
		Bytes:             s.Bytes,
	}
}

// MergeConfig stores merging, searching, and sorting settings
type MergeConfig struct {
	IncludeFilePath bool
	IncludeFileName bool
	Separator       string
	CountLines      bool
	CountWords      bool
	CountChars      bool
	CountCharsNoSpc bool
}
