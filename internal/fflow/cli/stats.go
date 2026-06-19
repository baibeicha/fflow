package cli

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"fflow/internal/fflow/locale"
	"fflow/internal/fflow/ui"
	"fflow/pkg/files"
)

var (
	statsPaths      []string
	statsRecursive  bool
	statsExtensions []string
	statsBlacklist  []string
	statsMinSize    string
	statsMaxSize    string
	statsCountLines bool
	statsCountWords bool
	statsCountChars bool
	statsCountNoSpc bool
	statsSortBy     []string
)

var statsCmd = &cobra.Command{
	Use:     "stats",
	Short:   locale.T("commands.stats.short"),
	Long:    locale.T("commands.stats.long"),
	Example: "",
	RunE:    runStats,
}

func init() {
	statsCmd.Flags().StringSliceVarP(&statsPaths, "path", "p", []string{"."}, locale.T("flags.path"))
	statsCmd.Flags().BoolVarP(&statsRecursive, "recursive", "r", false, locale.T("flags.recursive"))
	statsCmd.Flags().StringSliceVarP(&statsExtensions, "extensions", "e", nil, locale.T("flags.extensions"))
	statsCmd.Flags().StringSliceVar(&statsBlacklist, "blacklist", nil, locale.T("flags.blacklist"))
	statsCmd.Flags().StringVar(&statsMinSize, "min-size", "", locale.T("flags.min_size"))
	statsCmd.Flags().StringVar(&statsMaxSize, "max-size", "", locale.T("flags.max_size"))
	statsCmd.Flags().BoolVar(&statsCountLines, "count-lines", true, locale.T("flags.count_lines"))
	statsCmd.Flags().BoolVar(&statsCountWords, "count-words", true, locale.T("flags.count_words"))
	statsCmd.Flags().BoolVar(&statsCountChars, "count-chars", false, locale.T("flags.count_chars"))
	statsCmd.Flags().BoolVar(&statsCountNoSpc, "count-chars-no-space", false, locale.T("flags.count_chars_no_space"))
	statsCmd.Flags().StringSliceVar(&statsSortBy, "sort-by", []string{"name:asc"}, locale.T("flags.sort_by"))
}

func runStats(cmd *cobra.Command, args []string) error {
	ui.Title(locale.T("commands.stats.short"))

	fsc := files.NewFolderSearchConfig(statsRecursive, statsPaths...)

	if len(statsExtensions) > 0 {
		fsc.SearchForExtensions(statsExtensions...)
	}

	if len(statsBlacklist) > 0 {
		fsc.AddToBlackList(statsBlacklist...)
	}

	if statsMinSize != "" {
		size, unit := parseSize(statsMinSize)
		fsc.SetMinSize(files.SizeFromUnit(size, unit))
	}
	if statsMaxSize != "" {
		size, unit := parseSize(statsMaxSize)
		fsc.SetMaxSize(files.SizeFromUnit(size, unit))
	}

	statsCfg := &files.StatsConfig{
		CountLines:      statsCountLines,
		CountWords:      statsCountWords,
		CountChars:      statsCountChars,
		CountCharsNoSpc: statsCountNoSpc,
	}

	sorter := buildSorter(statsSortBy)

	startTime := time.Now()

	spinner := ui.NewSpinner(locale.T("messages.progress.collecting"))
	fileInfos, err := files.CollectFiles(fsc)
	spinner.Finish()

	if err != nil {
		return fmt.Errorf(locale.T("messages.errors.calculating_stats"), err)
	}

	if len(fileInfos) == 0 {
		ui.Warn(locale.T("messages.errors.no_files_found"))
		return nil
	}

	sorter.SortFiles(fileInfos)

	bar := ui.NewProgressBar(int64(len(fileInfos)), locale.T("messages.progress.calculating_stats"))

	stats, err := files.CalculateStats(fileInfos, statsCfg, func() {
		bar.Add(1)
	})
	bar.Finish()

	if err != nil {
		return fmt.Errorf(locale.T("messages.errors.calculating_stats"), err)
	}

	if stats.Files == 0 {
		ui.Warn(locale.T("messages.errors.no_files_found"))
		return nil
	}

	ui.Success(fmt.Sprintf(locale.T("messages.success.files_found"), stats.Files))

	elapsed := time.Since(startTime)
	printStatsResult(stats, elapsed)

	return nil
}

func printStatsResult(stats *files.FileStats, elapsed time.Duration) {
	ui.PrintSection(locale.T("messages.results.stats_title"))

	data := map[string]string{
		locale.T("messages.labels.files"):   ui.FormatNumber(stats.Files),
		locale.T("messages.labels.bytes"):   ui.FormatBytes(stats.Bytes),
		locale.T("messages.labels.elapsed"): elapsed.Round(time.Millisecond).String(),
	}

	if statsCountLines {
		data[locale.T("messages.labels.lines")] = ui.FormatNumber(stats.Lines)
	}
	if statsCountWords {
		data[locale.T("messages.labels.words")] = ui.FormatNumber(stats.Words)
	}
	if statsCountChars {
		data[locale.T("messages.labels.characters")] = ui.FormatNumber(stats.Characters)
	}
	if statsCountNoSpc {
		data[locale.T("messages.labels.characters_no_space")] = ui.FormatNumber(stats.CharactersNoSpace)
	}

	ui.PrintStatsTable(data)
	ui.Success(locale.T("messages.success.stats_collected"))
}

func buildSorter(sortBy []string) *files.MultiSorter {
	var criteria []files.SortCriteria

	for _, s := range sortBy {
		parts := splitSortCriteria(s)
		if len(parts) != 2 {
			continue
		}

		var field files.SortField
		switch parts[0] {
		case "name":
			field = files.SortByName
		case "modtime":
			field = files.SortByModTime
		case "size":
			field = files.SortBySize
		default:
			continue
		}

		var order files.SortOrder
		if parts[1] == "desc" {
			order = files.Descending
		} else {
			order = files.Ascending
		}

		criteria = append(criteria, files.SortCriteria{Field: field, Order: order})
	}

	return files.NewMultiSorter(criteria...)
}

func splitSortCriteria(s string) []string {
	for i, c := range s {
		if c == ':' {
			return []string{s[:i], s[i+1:]}
		}
	}
	return []string{s, "asc"}
}

func parseSize(s string) (int64, string) {
	s = strings.ToLower(strings.TrimSpace(s))

	unitIdx := -1
	for i, c := range s {
		if c < '0' || c > '9' {
			unitIdx = i
			break
		}
	}

	if unitIdx == -1 {
		size, _ := strconv.ParseInt(s, 10, 64)
		return size, ""
	}

	size, _ := strconv.ParseInt(s[:unitIdx], 10, 64)
	unit := s[unitIdx:]
	return size, unit
}
