package cli

import (
	"fmt"
	"strconv"
	"time"

	"github.com/baibeicha/fflow/pkg/telemetry"
	"github.com/spf13/cobra"

	"github.com/baibeicha/fflow/internal/fflow/locale"
	"github.com/baibeicha/fflow/internal/fflow/ui"
	"github.com/baibeicha/fflow/pkg/files"
)

var (
	movePaths      []string
	moveDests      []string
	moveRewrite    bool
	moveRecursive  bool
	moveExtensions []string
	moveBlacklist  []string
	moveMinSize    string
	moveMaxSize    string
	moveLimit      string
	moveOffset     string
)

var moveCmd = &cobra.Command{
	Use:     "move [paths...]",
	Short:   locale.T("commands.move.short"),
	Long:    locale.T("commands.move.long"),
	Example: "fflow move -r -e .mp4 --dest ./videos",
	Args:    cobra.ArbitraryArgs,
	RunE:    runMove,
}

func init() {
	moveCmd.Flags().StringSliceVarP(&movePaths, "path", "p", []string{"."}, locale.T("flags.path"))
	moveCmd.Flags().StringSliceVarP(&moveDests, "dest", "d", []string{"."}, locale.T("flags.dest"))
	moveCmd.Flags().BoolVarP(&moveRewrite, "rewrite", "w", false, locale.T("flags.rewrite"))
	moveCmd.Flags().BoolVarP(&moveRecursive, "recursive", "r", false, locale.T("flags.recursive"))
	moveCmd.Flags().StringSliceVarP(&moveExtensions, "extensions", "e", nil, locale.T("flags.extensions"))
	moveCmd.Flags().StringSliceVar(&moveBlacklist, "blacklist", nil, locale.T("flags.blacklist"))
	moveCmd.Flags().StringVar(&moveMinSize, "min-size", "", locale.T("flags.min_size"))
	moveCmd.Flags().StringVar(&moveMaxSize, "max-size", "", locale.T("flags.max_size"))
	moveCmd.Flags().StringVar(&moveLimit, "limit", "", locale.T("flags.limit"))
	moveCmd.Flags().StringVar(&moveOffset, "offset", "", locale.T("flags.offset"))
}

func runMove(cmd *cobra.Command, args []string) (err error) {
	start := time.Now()
	var recordedFiles, recordedBytes int64

	defer func() {
		telemetry.Record("move", recordedFiles, recordedBytes, start, err)
	}()

	pathsToSearch := movePaths
	if len(args) > 0 {
		pathsToSearch = args
	}

	ui.Title(locale.T("commands.move.short"))
	spinner := ui.NewSpinner(locale.T("messages.progress.collecting"))

	fsc := files.NewFolderSearchConfig(moveRecursive, pathsToSearch...)
	if len(moveExtensions) > 0 {
		fsc.SearchForExtensions(moveExtensions...)
	}
	if len(moveBlacklist) > 0 {
		fsc.AddToBlackList(moveBlacklist...)
	}
	if moveMinSize != "" {
		size, unit := parseSize(moveMinSize)
		fsc.SetMinSize(files.SizeFromUnit(size, unit))
	}
	if moveMaxSize != "" {
		size, unit := parseSize(moveMaxSize)
		fsc.SetMaxSize(files.SizeFromUnit(size, unit))
	}
	if moveLimit != "" {
		limit, err := strconv.ParseUint(moveLimit, 10, 64)
		if err != nil {
			fsc.SetLimit(limit)
		}
	}
	if moveOffset != "" {
		offset, err := strconv.ParseUint(moveOffset, 10, 64)
		if err != nil {
			fsc.SetLimit(offset)
		}
	}

	fsc.CollectDirs = false
	items, err := files.CollectFiles(fsc)
	items, amount := files.Paginate(items, fsc)
	spinner.Finish()

	if err != nil {
		return fmt.Errorf(locale.T("messages.errors.collecting_files"), err)
	}
	if len(items) == 0 {
		ui.Warn(locale.T("messages.errors.no_files_found"))
		return nil
	}

	start = time.Now()
	bar := ui.NewProgressBar(int64(amount), locale.T("messages.progress.moving"))
	suffix := locale.T("messages.labels.copy_suffix")

	_, stats, err := files.TransferFiles(items, moveDests, true, moveRewrite, suffix, func() {
		bar.Add(1)
	})
	bar.Finish()

	if err != nil {
		return fmt.Errorf("moving error: %w", err)
	}

	elapsed := time.Since(start)
	ui.PrintSection(locale.T("messages.results.move_title"))
	data := map[string]string{
		locale.T("messages.labels.files"):   ui.FormatNumber(stats.Files),
		locale.T("messages.labels.bytes"):   ui.FormatBytes(stats.Bytes),
		locale.T("messages.labels.elapsed"): elapsed.Round(time.Millisecond).String(),
	}
	ui.PrintStatsTable(data)

	ui.Success(locale.T("messages.success.files_moved"))

	return nil
}
