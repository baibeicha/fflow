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
	zipPaths      []string
	zipRecursive  bool
	zipExtensions []string
	zipBlacklist  []string
	zipOutput     string
)

var zipCmd = &cobra.Command{
	Use:   "zip [paths...]",
	Short: locale.T("commands.zip.short"),
	Long:  locale.T("commands.zip.long"),
	Args:  cobra.ArbitraryArgs,
	RunE:  runZip,
}

func init() {
	zipCmd.Flags().StringSliceVarP(&zipPaths, "path", "p", []string{"."}, locale.T("flags.path"))
	zipCmd.Flags().BoolVarP(&zipRecursive, "recursive", "r", false, locale.T("flags.recursive"))
	zipCmd.Flags().StringSliceVarP(&zipExtensions, "extensions", "e", nil, locale.T("flags.extensions"))
	zipCmd.Flags().StringSliceVar(&zipBlacklist, "blacklist", nil, locale.T("flags.blacklist"))
	zipCmd.Flags().StringVarP(&zipOutput, "output", "o", "archive.zip", locale.T("flags.output"))
}

func runZip(cmd *cobra.Command, args []string) (err error) {
	start := time.Now()
	var recordedFiles, recordedBytes int64

	defer func() {
		telemetry.Record("zip", recordedFiles, recordedBytes, start, err)
	}()

	paths := zipPaths
	if len(args) > 0 {
		paths = args
	}

	fsc := files.NewFolderSearchConfig(zipRecursive, paths...)
	if len(zipExtensions) > 0 {
		fsc.SearchForExtensions(zipExtensions...)
	}
	if len(zipBlacklist) > 0 {
		fsc.AddToBlackList(zipBlacklist...)
	}
	fsc.CollectDirs = false

	start = time.Now()

	ui.Title(locale.T("commands.zip.short"))
	items, err := files.CollectFiles(fsc)
	if err != nil {
		return fmt.Errorf(locale.T("messages.errors.collecting_files"), err)
	}
	if len(items) == 0 {
		ui.Warn(locale.T("messages.errors.no_files_found"))
		return nil
	}

	bar := ui.NewProgressBar(int64(len(items)), locale.T("messages.progress.zipping"))

	_, stats, err := files.CreateZip(items, zipOutput, func() { bar.Add(1) })
	bar.Finish()

	if err != nil {
		return err
	}

	recordedFiles = stats.Files
	recordedBytes = stats.Bytes

	ui.PrintStatsTable(map[string]string{
		locale.T("messages.labels.files"):   ui.FormatNumber(stats.Files),
		locale.T("messages.labels.bytes"):   ui.FormatBytes(stats.Bytes),
		locale.T("messages.labels.elapsed"): time.Since(start).Round(time.Millisecond).String(),
	})

	ui.Success(locale.T("messages.success.archive_created"))

	return nil
}
