package cli

import (
	"fmt"

	"github.com/baibeicha/fflow/internal/fflow/locale"
	"github.com/baibeicha/fflow/internal/fflow/ui"
	"github.com/baibeicha/fflow/pkg/telemetry"
	"github.com/spf13/cobra"
)

var (
	telemetryEndpoint string
	telemetrySilent   bool
)

var telemetryCmd = &cobra.Command{
	Use:   "telemetry",
	Short: locale.T("commands.telemetry.short"),
	Long:  locale.T("commands.telemetry.long"),
}

var onCmd = &cobra.Command{
	Use:   "on",
	Short: locale.T("commands.telemetry_on.short"),
	RunE: func(cmd *cobra.Command, args []string) error {
		state := telemetry.LoadState()
		state.Enabled = true
		state.HasPrompted = true
		telemetry.SaveState(state)
		ui.Success(locale.T("messages.success.telemetry_on"))
		return nil
	},
}

var offCmd = &cobra.Command{
	Use:   "off",
	Short: locale.T("commands.telemetry_off.short"),
	RunE: func(cmd *cobra.Command, args []string) error {
		state := telemetry.LoadState()
		state.Enabled = false
		state.HasPrompted = true
		telemetry.SaveState(state)
		ui.Success(locale.T("messages.success.telemetry_off"))
		return nil
	},
}

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: locale.T("commands.telemetry_push.short"),
	RunE: func(cmd *cobra.Command, args []string) error {
		endpoint := telemetryEndpoint

		if endpoint == "" {
			state := telemetry.LoadState()
			endpoint = state.Endpoint
		}

		if endpoint == "" {
			return fmt.Errorf(locale.T("messages.errors.telemetry_no_endpoint"))
		}

		var spinner *ui.ProgressBar
		if !telemetrySilent {
			spinner = ui.NewSpinner(locale.T("messages.progress.pushing_telemetry"))
		}

		err := telemetry.Push(endpoint)

		if spinner != nil {
			spinner.Finish()
		}

		if err != nil {
			if !telemetrySilent {
				return err
			}
			return nil
		}

		if !telemetrySilent {
			ui.Success(locale.T("messages.success.telemetry_pushed"))
		}

		return nil
	},
}

func init() {
	pushCmd.Flags().StringVarP(&telemetryEndpoint, "endpoint", "e", "", locale.T("flags.endpoint"))
	pushCmd.Flags().BoolVar(&telemetrySilent, "silent", false, "Hide output (for background runs)")
	_ = pushCmd.Flags().MarkHidden("silent")

	telemetryCmd.AddCommand(onCmd)
	telemetryCmd.AddCommand(offCmd)
	telemetryCmd.AddCommand(pushCmd)
}
