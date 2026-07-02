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
	delPaths      []string
	delRecursive  bool
	delExtensions []string
	delBlacklist  []string
	delMinSize    string
	delMaxSize    string
	delLimit      string
	delOffset     string
)

var deleteCmd = &cobra.Command{
	Use:   "delete [paths...]",
	Short: locale.T("commands.delete.short"),
	Long:  locale.T("commands.delete.long"),
	Args:  cobra.ArbitraryArgs,
	RunE:  runDelete,
}

func init() {
	deleteCmd.Flags().StringSliceVarP(&delPaths, "path", "p", []string{"."}, locale.T("flags.path"))
	deleteCmd.Flags().BoolVarP(&delRecursive, "recursive", "r", false, locale.T("flags.recursive"))
	deleteCmd.Flags().StringSliceVarP(&delExtensions, "extensions", "e", nil, locale.T("flags.extensions"))
	deleteCmd.Flags().StringSliceVar(&delBlacklist, "blacklist", nil, locale.T("flags.blacklist"))
	deleteCmd.Flags().StringVar(&delMinSize, "min-size", "", locale.T("flags.min_size"))
	deleteCmd.Flags().StringVar(&delMaxSize, "max-size", "", locale.T("flags.max_size"))
	deleteCmd.Flags().StringVar(&delLimit, "limit", "", locale.T("flags.limit"))
	deleteCmd.Flags().StringVar(&delOffset, "offset", "", locale.T("flags.offset"))
}

func runDelete(cmd *cobra.Command, args []string) (err error) {
	start := time.Now()
	var recordedFiles, recordedBytes int64

	defer func() {
		telemetry.Record("delete", recordedFiles, recordedBytes, start, err)
	}()

	paths := delPaths
	if len(args) > 0 {
		paths = args
	}

	fsc := files.NewFolderSearchConfig(delRecursive, paths...)
	if len(delExtensions) > 0 {
		fsc.SearchForExtensions(delExtensions...)
	}
	if len(delBlacklist) > 0 {
		fsc.AddToBlackList(delBlacklist...)
	}
	if delMinSize != "" {
		s, u := parseSize(delMinSize)
		fsc.SetMinSize(files.SizeFromUnit(s, u))
	}
	if delMaxSize != "" {
		s, u := parseSize(delMaxSize)
		fsc.SetMaxSize(files.SizeFromUnit(s, u))
	}
	if delLimit != "" {
		limit, err := strconv.ParseUint(delLimit, 10, 64)
		if err != nil {
			fsc.SetLimit(limit)
		}
	}
	if delOffset != "" {
		offset, err := strconv.ParseUint(delOffset, 10, 64)
		if err != nil {
			fsc.SetLimit(offset)
		}
	}
	fsc.CollectDirs = false

	start = time.Now()

	ui.Title(locale.T("commands.delete.short"))
	items, err := files.CollectFiles(fsc)
	items, amount := files.Paginate(items, fsc)

	if err != nil {
		return fmt.Errorf(locale.T("messages.errors.collecting_files"), err)
	}
	if len(items) == 0 {
		ui.Warn(locale.T("messages.errors.no_files_found"))
		return nil
	}

	bar := ui.NewProgressBar(int64(amount), locale.T("messages.progress.deleting"))

	_, stats, err := files.DeleteFiles(items, func() { bar.Add(1) })
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

	ui.Success(locale.T("messages.success.files_deleted"))

	return nil
}
