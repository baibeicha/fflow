package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
	"unicode/utf8"

	"github.com/baibeicha/fflow/pkg/telemetry"
	"github.com/spf13/cobra"

	"github.com/baibeicha/fflow/internal/fflow/locale"
	"github.com/baibeicha/fflow/internal/fflow/ui"
	"github.com/baibeicha/fflow/pkg/files"
)

var (
	infoPaths      []string
	infoRecursive  bool
	infoExtensions []string
	infoBlacklist  []string
	infoMinSize    string
	infoMaxSize    string
	infoSortBy     []string
	infoLimit      string
	infoOffset     string
)

var infoCmd = &cobra.Command{
	Use:     "info [path...]",
	Short:   locale.T("commands.info.short"),
	Long:    locale.T("commands.info.long"),
	Example: "fflow info . -r\nfflow info ./src -e .go,.md --min-size 1kb\nfflow info ./logs --sort-by size:desc",
	Args:    cobra.ArbitraryArgs,
	RunE:    runInfo,
}

func init() {
	infoCmd.Flags().StringSliceVarP(&infoPaths, "path", "p", []string{"."}, locale.T("flags.path"))
	infoCmd.Flags().BoolVarP(&infoRecursive, "recursive", "r", false, locale.T("flags.recursive"))
	infoCmd.Flags().StringSliceVarP(&infoExtensions, "extensions", "e", nil, locale.T("flags.extensions"))
	infoCmd.Flags().StringSliceVar(&infoBlacklist, "blacklist", nil, locale.T("flags.blacklist"))
	infoCmd.Flags().StringVar(&infoMinSize, "min-size", "", locale.T("flags.min_size"))
	infoCmd.Flags().StringVar(&infoMaxSize, "max-size", "", locale.T("flags.max_size"))
	infoCmd.Flags().StringSliceVar(&infoSortBy, "sort-by", []string{"name:asc"}, locale.T("flags.sort_by"))
	infoCmd.Flags().StringVar(&infoLimit, "limit", "", locale.T("flags.limit"))
	infoCmd.Flags().StringVar(&infoOffset, "offset", "", locale.T("flags.offset"))
}

func runInfo(cmd *cobra.Command, args []string) (err error) {
	start := time.Now()
	var recordedFiles, recordedBytes int64

	defer func() {
		telemetry.Record("info", recordedFiles, recordedBytes, start, err)
	}()

	pathsToAnalyze := infoPaths
	if len(args) > 0 {
		pathsToAnalyze = args
	}

	if len(pathsToAnalyze) == 1 {
		ui.Title(fmt.Sprintf(locale.T("messages.progress.analyzing_dir"), pathsToAnalyze[0]))
	} else {
		ui.Title(locale.T("commands.info.short"))
	}

	spinner := ui.NewSpinner(locale.T("messages.progress.collecting"))

	fsc := files.NewFolderSearchConfig(true, pathsToAnalyze...)
	fsc.CollectDirs = true

	if len(infoExtensions) > 0 {
		fsc.SearchForExtensions(infoExtensions...)
	}

	if len(infoBlacklist) > 0 {
		fsc.AddToBlackList(infoBlacklist...)
	}

	if infoMinSize != "" {
		size, unit := parseSize(infoMinSize)
		fsc.SetMinSize(files.SizeFromUnit(size, unit))
	}
	if infoMaxSize != "" {
		size, unit := parseSize(infoMaxSize)
		fsc.SetMaxSize(files.SizeFromUnit(size, unit))
	}
	if infoLimit != "" {
		limit, err := strconv.ParseUint(infoLimit, 10, 64)
		if err != nil {
			fsc.SetLimit(limit)
		}
	}
	if infoOffset != "" {
		offset, err := strconv.ParseUint(infoOffset, 10, 64)
		if err != nil {
			fsc.SetLimit(offset)
		}
	}

	items, err := files.CollectFiles(fsc)
	items, _ = files.Paginate(items, fsc)
	spinner.Finish()

	if err != nil {
		return fmt.Errorf(locale.T("messages.errors.collecting_files"), err)
	}

	if len(items) == 0 {
		ui.Warn(locale.T("messages.errors.no_files_found"))
		return nil
	}

	var totalFiles int64
	var totalSize int64
	var displayItems []files.FileInfo

	for _, item := range items {
		if !item.IsDir {
			totalFiles++
			totalSize += item.Size
		}

		if !infoRecursive {
			parent := filepath.Dir(item.Path)
			if parent != "." {
				continue
			}
		}

		displayItems = append(displayItems, item)
	}

	sorter := files.NewMultiSorter(
		files.SortCriteria{Field: files.SortByName, Order: files.Ascending},
	)
	sorter.SortFilesWithDirs(displayItems)

	ui.PrintSection(locale.T("messages.results.file_list"))

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)

	colName := locale.T("messages.labels.file_name")
	colSize := locale.T("messages.labels.size")
	colDate := locale.T("messages.labels.modified_time")
	colPath := locale.T("messages.labels.relative_path")

	fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", colName, colSize, colDate, colPath)
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
		strings.Repeat("-", utf8.RuneCountInString(colName)),
		strings.Repeat("-", utf8.RuneCountInString(colSize)),
		strings.Repeat("-", 19),
		strings.Repeat("-", utf8.RuneCountInString(colPath)),
	)

	for _, item := range displayItems {
		name := item.Name
		if item.IsDir {
			name += string(os.PathSeparator)
		}

		sizeStr := ui.FormatBytes(item.Size)
		timeStr := time.Unix(item.ModTime, 0).Format("2006-01-02 15:04:05")

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", name, sizeStr, timeStr, item.Path)
	}
	w.Flush()
	fmt.Println()

	ui.PrintSection(locale.T("messages.results.dir_stats"))

	statsData := map[string]string{
		locale.T("messages.labels.files"): ui.FormatNumber(totalFiles),
		locale.T("messages.labels.bytes"): ui.FormatBytes(totalSize),
	}
	ui.PrintStatsTable(statsData)

	ui.Success(locale.T("messages.success.info_collected"))
	return nil
}
