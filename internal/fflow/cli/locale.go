package cli

import (
	"fflow/internal/fflow/locale"
	"fflow/internal/fflow/ui"
	"fmt"

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

func runLocale(cmd *cobra.Command, args []string) error {
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
