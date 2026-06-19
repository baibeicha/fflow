package cli

import (
	"fflow/internal/fflow/locale"
	"fflow/internal/fflow/ui"
	"fflow/pkg/config"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
	verbose bool
	quiet   bool

	appCfg *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "fflow",
	Short: "fflow",
	Long:  "fflow",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initConfig()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	msgs := locale.GetMessages()
	if msgs != nil {
		rootCmd.Short = msgs.CLI.Description
		rootCmd.Long = msgs.CLI.Long
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", locale.T("flags.config"))
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, locale.T("flags.verbose"))
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, locale.T("flags.quiet"))

	rootCmd.AddCommand(statsCmd)
	rootCmd.AddCommand(mergeCmd)
	rootCmd.AddCommand(pdfCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(localeCmd)
	rootCmd.AddCommand(copyCmd)
	rootCmd.AddCommand(moveCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(renameCmd)
	rootCmd.AddCommand(extCmd)
	rootCmd.AddCommand(zipCmd)
	rootCmd.AddCommand(cmdExecCmd)
	rootCmd.AddCommand(flowCmd)

	ui.Init(quiet, verbose)

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.InitDefaultHelpCmd()
	for _, c := range rootCmd.Commands() {
		if c.Name() == "help" {
			c.Short = locale.T("commands.help.short")
			c.Long = locale.T("commands.help.long")
		}
	}

	usageTmpl := fmt.Sprintf(`%s:
  {{.UseLine}}
{{if .HasAvailableSubCommands}}
%s:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

%s:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

%s:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}

%s
`,
		locale.T("cobra.usage"),
		locale.T("cobra.available_commands"),
		locale.T("cobra.flags"),
		locale.T("cobra.global_flags"),
		locale.T("cobra.help_info"),
	)
	rootCmd.SetUsageTemplate(usageTmpl)
}

func initConfig() error {
	if cfgFile != "" {
		appCfg = config.MustLoad(cfgFile)
	} else {
		appCfg = config.MustLoadUserConfig("fflow",
			"locale", "en",
		)
	}

	if verbose {
		ui.Info(locale.T("messages.success.config_loaded"))
	}

	return nil
}

func exitWithError(format string, args ...interface{}) {
	ui.Error(fmt.Sprintf(format, args...))
	os.Exit(1)
}
