package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/baibeicha/fflow/internal/fflow/locale"
	"github.com/baibeicha/fflow/internal/fflow/ui"
	"github.com/baibeicha/fflow/pkg/config"
	"github.com/baibeicha/fflow/pkg/telemetry"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	verbose bool
	quiet   bool
	appCfg  *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "fflow",
	Short: "fflow",
	Long:  "fflow",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		checkTelemetryPrompt(cmd)
		return nil
	},
}

// Execute starts the application CLI.
func Execute() error {
	return rootCmd.Execute()
}

func checkTelemetryPrompt(cmd *cobra.Command) {
	if cmd.Name() == "help" || cmd.Name() == "completion" || cmd.Name() == "telemetry" ||
		(cmd.Parent() != nil && cmd.Parent().Name() == "telemetry") {
		return
	}

	state := telemetry.LoadState()
	if state.HasPrompted {
		return
	}

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		state.Enabled = false
		state.HasPrompted = true
		telemetry.SaveState(state)
		return
	}

	fmt.Print(locale.T("messages.prompts.telemetry"))
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	if input == "y" || input == "yes" || input == "д" || input == "да" {
		state.Enabled = true
		fmt.Println(locale.T("messages.success.telemetry_on"))
	} else {
		state.Enabled = false
		fmt.Println(locale.T("messages.success.telemetry_off"))
	}

	state.HasPrompted = true
	telemetry.SaveState(state)
	fmt.Println()
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
	rootCmd.AddCommand(telemetryCmd)

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
