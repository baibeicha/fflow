package cli

import (
	"fmt"
	"time"

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
	defaultPaths := []string{"."}

	moveCmd.Flags().StringSliceVarP(&movePaths, "path", "p", defaultPaths, locale.T("flags.path"))
	moveCmd.Flags().StringSliceVarP(&moveDests, "dest", "d", []string{"."}, locale.T("flags.dest"))
	moveCmd.Flags().BoolVarP(&moveRewrite, "rewrite", "w", false, locale.T("flags.rewrite"))

	moveCmd.Flags().BoolVarP(&moveRecursive, "recursive", "r", false, locale.T("flags.recursive"))
	moveCmd.Flags().StringSliceVarP(&moveExtensions, "extensions", "e", nil, locale.T("flags.extensions"))
	moveCmd.Flags().StringSliceVar(&moveBlacklist, "blacklist", nil, locale.T("flags.blacklist"))
	moveCmd.Flags().StringVar(&moveMinSize, "min-size", "", locale.T("flags.min_size"))
	moveCmd.Flags().StringVar(&moveMaxSize, "max-size", "", locale.T("flags.max_size"))
}

func runMove(cmd *cobra.Command, args []string) error {
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

	fsc.CollectDirs = false
	items, err := files.CollectFiles(fsc)
	spinner.Finish()

	if err != nil {
		return fmt.Errorf(locale.T("messages.errors.collecting_files"), err)
	}
	if len(items) == 0 {
		ui.Warn(locale.T("messages.errors.no_files_found"))
		return nil
	}

	startTime := time.Now()
	bar := ui.NewProgressBar(int64(len(items)), locale.T("messages.progress.moving"))
	suffix := locale.T("messages.labels.copy_suffix")

	_, stats, err := files.TransferFiles(items, moveDests, true, moveRewrite, suffix, func() {
		bar.Add(1)
	})
	bar.Finish()

	if err != nil {
		return fmt.Errorf("перемещение завершилось с ошибкой: %w", err)
	}

	elapsed := time.Since(startTime)
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
