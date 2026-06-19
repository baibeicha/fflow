package ui

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
)

var (
	quiet   bool
	verbose bool
)

func Init(quietMode, verboseMode bool) {
	quiet = quietMode
	verbose = verboseMode
}

func IsQuiet() bool {
	return quiet
}

func IsVerbose() bool {
	return verbose
}

type ProgressBar struct {
	bar         *progressbar.ProgressBar
	description string
	startTime   time.Time
}

func NewProgressBar(max int64, description string) *ProgressBar {
	if quiet {
		return &ProgressBar{description: description, startTime: time.Now()}
	}

	bar := progressbar.NewOptions64(
		max,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetDescription("[cyan]"+description+"[reset]"),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionSetWidth(40),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]█[reset]",
			SaucerHead:    "[green]█[reset]",
			SaucerPadding: "[dark_gray]░[reset]",
			BarStart:      "[blue]|[reset]",
			BarEnd:        "[blue]|[reset]",
		}),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionSetItsString("files"),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionSetElapsedTime(true),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionClearOnFinish(),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
	)

	return &ProgressBar{
		bar:         bar,
		description: description,
		startTime:   time.Now(),
	}
}

func NewSpinner(description string) *ProgressBar {
	if quiet {
		return &ProgressBar{description: description, startTime: time.Now()}
	}

	bar := progressbar.NewOptions64(
		-1,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetDescription("[cyan]"+description+"[reset]"),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionClearOnFinish(),
	)

	return &ProgressBar{
		bar:         bar,
		description: description,
		startTime:   time.Now(),
	}
}

func (p *ProgressBar) Add(n int) {
	if p.bar != nil {
		if err := p.bar.Add(n); err != nil {
			slog.Error("Error adding progressbar", "error", err)
			return
		}
	}
}

func (p *ProgressBar) Set(n int) {
	if p.bar != nil {
		if err := p.bar.Set(n); err != nil {
			slog.Error("Error setting progressbar", "error", err)
			return
		}
	}
}

func (p *ProgressBar) SetMax(max int) {
	if p.bar != nil {
		p.bar.ChangeMax(max)
	}
}

func (p *ProgressBar) Finish() {
	if p.bar != nil {
		if err := p.bar.Finish(); err != nil {
			slog.Error("Error finishing progressbar", "error", err)
			return
		}
	}
}

func (p *ProgressBar) Describe(description string) {
	if p.bar != nil {
		p.bar.Describe(description)
	}
	p.description = description
}

func (p *ProgressBar) ElapsedTime() time.Duration {
	return time.Since(p.startTime)
}
