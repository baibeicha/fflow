package cli

import (
	"fmt"

	"github.com/baibeicha/fflow/internal/fflow/locale"
	"github.com/baibeicha/fflow/internal/fflow/ui"

	"github.com/spf13/cobra"
)

var localeCmd = &cobra.Command{
	Use:     "locale [en|ru]",
	Short:   locale.T("commands.locale.short"),
	Long:    locale.T("commands.locale.long"),
	Example: "",
	Args:    cobra.ExactArgs(1),
	RunE:    runLocale,
}

func runLocale(cmd *cobra.Command, args []string) (err error) {
	newLocale := args[0]

	if err := locale.SaveLocale(newLocale); err != nil {
		return fmt.Errorf(locale.T("messages.errors.invalid_locale"))
	}

	if err := locale.SetLocale(newLocale); err != nil {
		return fmt.Errorf(locale.T("messages.errors.invalid_locale"))
	}

	localeName := "English"
	if newLocale == "ru" {
		localeName = "Русский"
	}

	ui.Success(fmt.Sprintf(locale.T("messages.success.locale_changed"), localeName))
	return nil
}
