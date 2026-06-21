package cli

import (
	"fmt"
	"time"

	"github.com/baibeicha/fflow/pkg/telemetry"
	"github.com/spf13/cobra"

	"github.com/baibeicha/fflow/internal/fflow/locale"
	"github.com/baibeicha/fflow/internal/fflow/ui"
	"github.com/baibeicha/fflow/pkg/files"
)

var (
	mergePaths       []string
	mergeOutput      string
	mergeRecursive   bool
	mergeExtensions  []string
	mergeBlacklist   []string
	mergeMinSize     string
	mergeMaxSize     string
	mergeFast        bool
	mergeIncludePath bool
	mergeIncludeName bool
	mergeSeparator   string
	mergeSortBy      []string
	mergeCountLines  bool
	mergeCountWords  bool
	mergeCountChars  bool
	mergeCountNoSpc  bool
)

var mergeCmd = &cobra.Command{
	Use:     "merge",
	Short:   locale.T("commands.merge.short"),
	Long:    locale.T("commands.merge.long"),
	Example: "",
	RunE:    runMerge,
}

func init() {
	mergeCmd.Flags().StringSliceVarP(&mergePaths, "path", "p", []string{"."}, locale.T("flags.path"))
	mergeCmd.Flags().StringVarP(&mergeOutput, "output", "o", "merged.txt", locale.T("flags.output"))
	mergeCmd.Flags().BoolVarP(&mergeRecursive, "recursive", "r", false, locale.T("flags.recursive"))
	mergeCmd.Flags().StringSliceVarP(&mergeExtensions, "extensions", "e", nil, locale.T("flags.extensions"))
	mergeCmd.Flags().StringSliceVar(&mergeBlacklist, "blacklist", nil, locale.T("flags.blacklist"))
	mergeCmd.Flags().StringVar(&mergeMinSize, "min-size", "", locale.T("flags.min_size"))
	mergeCmd.Flags().StringVar(&mergeMaxSize, "max-size", "", locale.T("flags.max_size"))
	mergeCmd.Flags().BoolVarP(&mergeFast, "fast", "f", false, locale.T("flags.fast"))
	mergeCmd.Flags().BoolVar(&mergeIncludePath, "include-path", false, locale.T("flags.include_path"))
	mergeCmd.Flags().BoolVar(&mergeIncludeName, "include-name", true, locale.T("flags.include_name"))
	mergeCmd.Flags().StringVar(&mergeSeparator, "separator", "", locale.T("flags.separator"))
	mergeCmd.Flags().StringSliceVar(&mergeSortBy, "sort-by", []string{"name:asc"}, locale.T("flags.sort_by"))
	mergeCmd.Flags().BoolVar(&mergeCountLines, "count-lines", true, locale.T("flags.count_lines"))
	mergeCmd.Flags().BoolVar(&mergeCountWords, "count-words", false, locale.T("flags.count_words"))
	mergeCmd.Flags().BoolVar(&mergeCountChars, "count-chars", false, locale.T("flags.count_chars"))
	mergeCmd.Flags().BoolVar(&mergeCountNoSpc, "count-chars-no-space", false, locale.T("flags.count_chars_no_space"))
}

func runMerge(cmd *cobra.Command, args []string) (err error) {
	start := time.Now()
	var recordedFiles, recordedBytes int64

	defer func() {
		telemetry.Record("merge", recordedFiles, recordedBytes, start, err)
	}()

	ui.Title(locale.T("commands.merge.short"))

	fsc := files.NewFolderSearchConfig(mergeRecursive, mergePaths...)

	if len(mergeExtensions) > 0 {
		fsc.SearchForExtensions(mergeExtensions...)
	}

	if len(mergeBlacklist) > 0 {
		fsc.AddToBlackList(mergeBlacklist...)
	}

	if mergeMinSize != "" {
		size, unit := parseSize(mergeMinSize)
		fsc.SetMinSize(files.SizeFromUnit(size, unit))
	}
	if mergeMaxSize != "" {
		size, unit := parseSize(mergeMaxSize)
		fsc.SetMaxSize(files.SizeFromUnit(size, unit))
	}

	mergeCfg := &files.MergeConfig{
		IncludeFilePath: mergeIncludePath,
		IncludeFileName: mergeIncludeName,
		Separator:       mergeSeparator,
		CountLines:      mergeCountLines,
		CountWords:      mergeCountWords,
		CountChars:      mergeCountChars,
		CountCharsNoSpc: mergeCountNoSpc,
	}

	sorter := buildSorter(mergeSortBy)

	mode := locale.T("messages.labels.mode_full")
	if mergeFast {
		mode = locale.T("messages.labels.mode_fast")
	}

	start = time.Now()

	spinner := ui.NewSpinner(locale.T("messages.progress.collecting"))
	fileInfos, err := files.CollectFiles(fsc)
	spinner.Finish()

	if err != nil {
		return fmt.Errorf(locale.T("messages.errors.merging_files"), err)
	}

	if len(fileInfos) == 0 {
		ui.Warn(locale.T("messages.errors.no_files_found"))
		return nil
	}

	sorter.SortFiles(fileInfos)

	bar := ui.NewProgressBar(int64(len(fileInfos)), fmt.Sprintf(locale.T("messages.progress.merging"), mode))

	stats, err := files.MergeFiles(fileInfos, mergeOutput, mergeCfg, mergeFast, func() {
		bar.Add(1)
	})
	bar.Finish()

	if err != nil {
		return fmt.Errorf(locale.T("messages.errors.merging_files"), err)
	}

	if stats.Files == 0 {
		ui.Warn(locale.T("messages.errors.no_files_found"))
		return nil
	}

	recordedFiles = stats.Files
	recordedBytes = stats.Bytes

	ui.Success(fmt.Sprintf(locale.T("messages.success.files_found"), stats.Files))

	elapsed := time.Since(start)
	printMergeResult(stats, elapsed)

	return nil
}

func printMergeResult(stats *files.FileStats, elapsed time.Duration) {
	ui.PrintSection(locale.T("messages.results.merge_title"))

	mode := locale.T("messages.labels.mode_full")
	if mergeFast {
		mode = locale.T("messages.labels.mode_fast")
	}

	data := map[string]string{
		locale.T("messages.labels.output_file"): mergeOutput,
		locale.T("messages.labels.mode"):        mode,
		locale.T("messages.labels.files"):       ui.FormatNumber(stats.Files),
		locale.T("messages.labels.bytes"):       ui.FormatBytes(stats.Bytes),
		locale.T("messages.labels.elapsed"):     elapsed.Round(time.Millisecond).String(),
	}

	if !mergeFast {
		if mergeCountLines {
			data[locale.T("messages.labels.lines")] = ui.FormatNumber(stats.Lines)
		}
		if mergeCountWords {
			data[locale.T("messages.labels.words")] = ui.FormatNumber(stats.Words)
		}
		if mergeCountChars {
			data[locale.T("messages.labels.characters")] = ui.FormatNumber(stats.Characters)
		}
		if mergeCountNoSpc {
			data[locale.T("messages.labels.characters_no_space")] = ui.FormatNumber(stats.CharactersNoSpace)
		}
	}

	ui.PrintStatsTable(data)
	ui.Success(fmt.Sprintf(locale.T("messages.success.files_merged"), mergeOutput))
}
