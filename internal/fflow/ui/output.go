package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

var (
	colorSuccess = color.New(color.FgGreen, color.Bold)
	colorError   = color.New(color.FgRed, color.Bold)
	colorInfo    = color.New(color.FgCyan)
	colorWarn    = color.New(color.FgYellow)
	colorTitle   = color.New(color.FgMagenta, color.Bold)
	colorMuted   = color.New(color.FgWhite, color.Faint)
)

func Success(format string, args ...interface{}) {
	if quiet {
		return
	}
	msg := fmt.Sprintf(format, args...)
	colorSuccess.Fprintf(os.Stdout, "[OK] %s\n", msg)
}

func Error(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	colorError.Fprintf(os.Stderr, "[ERROR] %s\n", msg)
}

func Info(format string, args ...interface{}) {
	if quiet {
		return
	}
	msg := fmt.Sprintf(format, args...)
	colorInfo.Fprintf(os.Stdout, "[INFO] %s\n", msg)
}

func Warn(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	colorWarn.Fprintf(os.Stderr, "[WARN] %s\n", msg)
}

func Title(format string, args ...interface{}) {
	if quiet {
		return
	}
	msg := fmt.Sprintf(format, args...)
	fmt.Println()
	colorTitle.Fprintf(os.Stdout, "%s\n", msg)
	colorMuted.Fprintf(os.Stdout, "%s\n", strings.Repeat("-", 50))
}

func PrintStatsTable(data map[string]string) {
	if quiet {
		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false
	t.Style().Format.HeaderAlign = text.AlignLeft
	t.Style().Format.RowAlign = text.AlignLeft

	for key, value := range data {
		t.AppendRow([]interface{}{
			colorMuted.Sprintf("  %s:", key),
			colorSuccess.Sprintf("%s", value),
		})
	}

	t.Render()
}

func FormatBytes(bytes int64) string {
	return humanize.IBytes(uint64(bytes))
}

func FormatNumber(n int64) string {
	return humanize.Comma(n)
}

func PrintKeyValue(key, value string) {
	if quiet {
		return
	}
	fmt.Printf("  %s %s\n", colorMuted.Sprintf("%s:", key), colorSuccess.Sprintf("%s", value))
}

func PrintSection(title string) {
	if quiet {
		return
	}
	fmt.Println()
	colorTitle.Fprintf(os.Stdout, "%s\n", title)
	colorMuted.Fprintf(os.Stdout, "%s\n", strings.Repeat("-", 50))
}
