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
	renPaths      []string
	renRecursive  bool
	renExtensions []string
	renBlacklist  []string
	renSearch     string
	renReplace    string
	renPrefix     string
	renSuffix     string
)

var renameCmd = &cobra.Command{
	Use:   "rename [paths...]",
	Short: locale.T("commands.rename.short"),
	Long:  locale.T("commands.rename.long"),
	Args:  cobra.ArbitraryArgs,
	RunE:  runRename,
}

func init() {
	renameCmd.Flags().StringSliceVarP(&renPaths, "path", "p", []string{"."}, locale.T("flags.path"))
	renameCmd.Flags().BoolVarP(&renRecursive, "recursive", "r", false, locale.T("flags.recursive"))
	renameCmd.Flags().StringSliceVarP(&renExtensions, "extensions", "e", nil, locale.T("flags.extensions"))
	renameCmd.Flags().StringSliceVar(&renBlacklist, "blacklist", nil, locale.T("flags.blacklist"))
	renameCmd.Flags().StringVar(&renSearch, "search", "", locale.T("flags.search"))
	renameCmd.Flags().StringVar(&renReplace, "replace", "", locale.T("flags.replace"))
	renameCmd.Flags().StringVar(&renPrefix, "prefix", "", locale.T("flags.prefix"))
	renameCmd.Flags().StringVar(&renSuffix, "suffix", "", locale.T("flags.suffix"))
}

func runRename(cmd *cobra.Command, args []string) error {
	paths := renPaths
	if len(args) > 0 {
		paths = args
	}

	fsc := files.NewFolderSearchConfig(renRecursive, paths...)
	if len(renExtensions) > 0 {
		fsc.SearchForExtensions(renExtensions...)
	}
	if len(renBlacklist) > 0 {
		fsc.AddToBlackList(renBlacklist...)
	}
	fsc.CollectDirs = false

	ui.Title(locale.T("commands.rename.short"))
	items, err := files.CollectFiles(fsc)
	if err != nil {
		return fmt.Errorf(locale.T("messages.errors.collecting_files"), err)
	}
	if len(items) == 0 {
		ui.Warn(locale.T("messages.errors.no_files_found"))
		return nil
	}

	start := time.Now()
	bar := ui.NewProgressBar(int64(len(items)), locale.T("messages.progress.renaming"))

	_, stats, err := files.RenameFiles(items, renSearch, renReplace, renPrefix, renSuffix, func() { bar.Add(1) })
	bar.Finish()

	if err != nil {
		return err
	}
	ui.PrintStatsTable(map[string]string{
		locale.T("messages.labels.files"):   ui.FormatNumber(stats.Files),
		locale.T("messages.labels.elapsed"): time.Since(start).Round(time.Millisecond).String(),
	})
	ui.Success(locale.T("messages.success.files_renamed"))
	return nil
}
