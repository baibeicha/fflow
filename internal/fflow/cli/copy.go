package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"fflow/internal/fflow/locale"
	"fflow/internal/fflow/ui"
	"fflow/pkg/files"
)

var (
	copyPaths      []string
	copyDests      []string
	copyRewrite    bool
	copyRecursive  bool
	copyExtensions []string
	copyBlacklist  []string
	copyMinSize    string
	copyMaxSize    string
)

var copyCmd = &cobra.Command{
	Use:     "copy [paths...]",
	Short:   locale.T("commands.copy.short"),
	Long:    locale.T("commands.copy.long"),
	Example: "fflow copy -r -e .txt --dest ./backup1 --dest ./backup2",
	Args:    cobra.ArbitraryArgs,
	RunE:    runCopy,
}

func init() {
	defaultPaths := []string{"."}

	copyCmd.Flags().StringSliceVarP(&copyPaths, "path", "p", defaultPaths, locale.T("flags.path"))
	copyCmd.Flags().StringSliceVarP(&copyDests, "dest", "d", []string{"."}, locale.T("flags.dest"))
	copyCmd.Flags().BoolVarP(&copyRewrite, "rewrite", "w", false, locale.T("flags.rewrite"))

	copyCmd.Flags().BoolVarP(&copyRecursive, "recursive", "r", false, locale.T("flags.recursive"))
	copyCmd.Flags().StringSliceVarP(&copyExtensions, "extensions", "e", nil, locale.T("flags.extensions"))
	copyCmd.Flags().StringSliceVar(&copyBlacklist, "blacklist", nil, locale.T("flags.blacklist"))
	copyCmd.Flags().StringVar(&copyMinSize, "min-size", "", locale.T("flags.min_size"))
	copyCmd.Flags().StringVar(&copyMaxSize, "max-size", "", locale.T("flags.max_size"))
}

func runCopy(cmd *cobra.Command, args []string) error {
	pathsToSearch := copyPaths
	if len(args) > 0 {
		pathsToSearch = args
	}

	ui.Title(locale.T("commands.copy.short"))
	spinner := ui.NewSpinner(locale.T("messages.progress.collecting"))

	fsc := files.NewFolderSearchConfig(copyRecursive, pathsToSearch...)
	if len(copyExtensions) > 0 {
		fsc.SearchForExtensions(copyExtensions...)
	}
	if len(copyBlacklist) > 0 {
		fsc.AddToBlackList(copyBlacklist...)
	}
	if copyMinSize != "" {
		size, unit := parseSize(copyMinSize)
		fsc.SetMinSize(files.SizeFromUnit(size, unit))
	}
	if copyMaxSize != "" {
		size, unit := parseSize(copyMaxSize)
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
	bar := ui.NewProgressBar(int64(len(items)), locale.T("messages.progress.copying"))
	suffix := locale.T("messages.labels.copy_suffix")

	_, stats, err := files.TransferFiles(items, copyDests, false, copyRewrite, suffix, func() {
		bar.Add(1)
	})
	bar.Finish()

	if err != nil {
		return fmt.Errorf("передача завершилась с ошибкой: %w", err)
	}

	elapsed := time.Since(startTime)
	ui.PrintSection(locale.T("messages.results.copy_title"))
	data := map[string]string{
		locale.T("messages.labels.files"):   ui.FormatNumber(stats.Files),
		locale.T("messages.labels.bytes"):   ui.FormatBytes(stats.Bytes),
		locale.T("messages.labels.elapsed"): elapsed.Round(time.Millisecond).String(),
	}
	ui.PrintStatsTable(data)
	ui.Success(locale.T("messages.success.files_copied"))
	return nil
}
