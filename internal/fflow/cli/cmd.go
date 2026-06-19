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
	cmdPaths      []string
	cmdRecursive  bool
	cmdExtensions []string
	cmdBlacklist  []string
	cmdExec       string
)

var cmdExecCmd = &cobra.Command{
	Use:   "cmd [paths...]",
	Short: locale.T("commands.cmd.short"),
	Long:  locale.T("commands.cmd.long"),
	Args:  cobra.ArbitraryArgs,
	RunE:  runCmd,
}

func init() {
	cmdExecCmd.Flags().StringSliceVarP(&cmdPaths, "path", "p", []string{"."}, locale.T("flags.path"))
	cmdExecCmd.Flags().BoolVarP(&cmdRecursive, "recursive", "r", false, locale.T("flags.recursive"))
	cmdExecCmd.Flags().StringSliceVarP(&cmdExtensions, "extensions", "e", nil, locale.T("flags.extensions"))
	cmdExecCmd.Flags().StringSliceVar(&cmdBlacklist, "blacklist", nil, locale.T("flags.blacklist"))
	cmdExecCmd.Flags().StringVar(&cmdExec, "exec", "", locale.T("flags.exec"))
	cmdExecCmd.MarkFlagRequired("exec")
}

func runCmd(cmd *cobra.Command, args []string) error {
	paths := cmdPaths
	if len(args) > 0 {
		paths = args
	}

	fsc := files.NewFolderSearchConfig(cmdRecursive, paths...)
	if len(cmdExtensions) > 0 {
		fsc.SearchForExtensions(cmdExtensions...)
	}
	if len(cmdBlacklist) > 0 {
		fsc.AddToBlackList(cmdBlacklist...)
	}
	fsc.CollectDirs = false

	ui.Title(locale.T("commands.cmd.short"))
	items, err := files.CollectFiles(fsc)
	if err != nil {
		return fmt.Errorf(locale.T("messages.errors.collecting_files"), err)
	}
	if len(items) == 0 {
		ui.Warn(locale.T("messages.errors.no_files_found"))
		return nil
	}

	start := time.Now()
	bar := ui.NewProgressBar(int64(len(items)), locale.T("messages.progress.executing"))

	_, stats, err := files.ExecuteCommandOnFiles(items, cmdExec, func() { bar.Add(1) })
	bar.Finish()

	if err != nil {
		return err
	}
	ui.PrintStatsTable(map[string]string{
		locale.T("messages.labels.files"):   ui.FormatNumber(stats.Files),
		locale.T("messages.labels.elapsed"): time.Since(start).Round(time.Millisecond).String(),
	})
	ui.Success(locale.T("messages.success.cmd_executed"))
	return nil
}
