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
	extPaths      []string
	extRecursive  bool
	extExtensions []string
	extBlacklist  []string
	extTo         string
)

var extCmd = &cobra.Command{
	Use:   "ext [paths...]",
	Short: locale.T("commands.ext.short"),
	Long:  locale.T("commands.ext.long"),
	Args:  cobra.ArbitraryArgs,
	RunE:  runExt,
}

func init() {
	extCmd.Flags().StringSliceVarP(&extPaths, "path", "p", []string{"."}, locale.T("flags.path"))
	extCmd.Flags().BoolVarP(&extRecursive, "recursive", "r", false, locale.T("flags.recursive"))
	extCmd.Flags().StringSliceVarP(&extExtensions, "extensions", "e", nil, locale.T("flags.extensions"))
	extCmd.Flags().StringSliceVar(&extBlacklist, "blacklist", nil, locale.T("flags.blacklist"))
	extCmd.Flags().StringVar(&extTo, "to", "", locale.T("flags.to"))
	extCmd.MarkFlagRequired("to")
}

func runExt(cmd *cobra.Command, args []string) (err error) {
	start := time.Now()
	var recordedFiles, recordedBytes int64

	defer func() {
		telemetry.Record("ext", recordedFiles, recordedBytes, start, err)
	}()

	paths := extPaths
	if len(args) > 0 {
		paths = args
	}

	fsc := files.NewFolderSearchConfig(extRecursive, paths...)
	if len(extExtensions) > 0 {
		fsc.SearchForExtensions(extExtensions...)
	}
	if len(extBlacklist) > 0 {
		fsc.AddToBlackList(extBlacklist...)
	}
	fsc.CollectDirs = false

	start = time.Now()

	ui.Title(locale.T("commands.ext.short"))
	items, err := files.CollectFiles(fsc)
	if err != nil {
		return fmt.Errorf(locale.T("messages.errors.collecting_files"), err)
	}
	if len(items) == 0 {
		ui.Warn(locale.T("messages.errors.no_files_found"))
		return nil
	}

	bar := ui.NewProgressBar(int64(len(items)), locale.T("messages.progress.renaming"))

	_, stats, err := files.ChangeExtension(items, extTo, func() { bar.Add(1) })
	bar.Finish()

	if err != nil {
		return err
	}

	recordedFiles = stats.Files
	recordedBytes = stats.Bytes

	ui.PrintStatsTable(map[string]string{
		locale.T("messages.labels.files"):   ui.FormatNumber(stats.Files),
		locale.T("messages.labels.elapsed"): time.Since(start).Round(time.Millisecond).String(),
	})

	ui.Success(locale.T("messages.success.files_renamed"))

	return nil
}
