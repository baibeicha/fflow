package cli

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"fflow/internal/fflow/locale"
	"fflow/internal/fflow/ui"
	"fflow/pkg/files"
)

var (
	flowPaths      []string
	flowRecursive  bool
	flowExtensions []string
	flowBlacklist  []string
	flowMinSize    string
	flowMaxSize    string
	flowYamlFile   string
	flowInlineCmd  string
	flowLoadEnv    bool
)

var flowCmd = &cobra.Command{
	Use:     "flow",
	Short:   locale.T("commands.flow.short"),
	Long:    locale.T("commands.flow.long"),
	Example: `fflow flow -re .go -c "var OUT=./bkp ; copy -d ${OUT} ; ext --to .txt"`,
	RunE:    runFlow,
}

func init() {
	flowCmd.Flags().StringSliceVarP(&flowPaths, "path", "p", []string{"."}, locale.T("flags.path"))
	flowCmd.Flags().BoolVarP(&flowRecursive, "recursive", "r", false, locale.T("flags.recursive"))
	flowCmd.Flags().StringSliceVarP(&flowExtensions, "extensions", "e", nil, locale.T("flags.extensions"))
	flowCmd.Flags().StringSliceVar(&flowBlacklist, "blacklist", nil, locale.T("flags.blacklist"))
	flowCmd.Flags().StringVar(&flowMinSize, "min-size", "", locale.T("flags.min_size"))
	flowCmd.Flags().StringVar(&flowMaxSize, "max-size", "", locale.T("flags.max_size"))

	flowCmd.Flags().StringVarP(&flowYamlFile, "file", "f", "", locale.T("flags.flow_file"))
	flowCmd.Flags().StringVarP(&flowInlineCmd, "cmd", "c", "", locale.T("flags.flow_cmd"))
	flowCmd.Flags().BoolVar(&flowLoadEnv, "env", false, locale.T("flags.flow_env"))
}

func runFlow(cmd *cobra.Command, args []string) error {
	if flowYamlFile == "" && flowInlineCmd == "" {
		return fmt.Errorf(locale.T("messages.errors.flow_no_pipeline"))
	}

	ui.Title(locale.T("commands.flow.short"))

	fsc := files.NewFolderSearchConfig(flowRecursive, flowPaths...)
	if len(flowExtensions) > 0 {
		fsc.SearchForExtensions(flowExtensions...)
	}
	if len(flowBlacklist) > 0 {
		fsc.AddToBlackList(flowBlacklist...)
	}
	if flowMinSize != "" {
		s, u := parseSize(flowMinSize)
		fsc.SetMinSize(files.SizeFromUnit(s, u))
	}
	if flowMaxSize != "" {
		s, u := parseSize(flowMaxSize)
		fsc.SetMaxSize(files.SizeFromUnit(s, u))
	}
	fsc.CollectDirs = false

	spinner := ui.NewSpinner(locale.T("messages.progress.collecting"))
	items, err := files.CollectFiles(fsc)
	spinner.Finish()

	if err != nil {
		return fmt.Errorf(locale.T("messages.errors.collecting_files"), err)
	}
	if len(items) == 0 {
		ui.Warn(locale.T("messages.errors.no_files_found"))
		return nil
	}

	var pipeline *files.Pipeline
	if flowYamlFile != "" {
		pipeline, err = parseYamlPipeline(flowYamlFile)
	} else {
		pipeline, err = parseInlinePipeline(flowInlineCmd)
	}

	if err != nil {
		return fmt.Errorf("pipeline error: %w", err)
	}

	if flowLoadEnv {
		if pipeline.Env == nil {
			pipeline.Env = make(map[string]string)
		}
		for _, e := range os.Environ() {
			pair := strings.SplitN(e, "=", 2)
			if len(pair) == 2 {
				pipeline.Env[pair[0]] = pair[1]
			}
		}
	}

	ctx := &files.FlowContext{
		Files: items,
		Vars:  make(map[string]string),
	}

	start := time.Now()
	var currentBar *ui.ProgressBar
	var lastAction string

	progressFn := func(action string, current, total int) {
		if currentBar == nil || action != lastAction {
			if currentBar != nil {
				currentBar.Finish()
				fmt.Println()
			}
			msg := fmt.Sprintf(locale.T("messages.progress.flow_step"), action)
			currentBar = ui.NewProgressBar(int64(total), msg)
			lastAction = action
		}
		currentBar.Add(1)
	}

	err = files.RunFlow(ctx, pipeline, progressFn)
	if currentBar != nil {
		currentBar.Finish()
		fmt.Println()
	}

	if err != nil {
		return fmt.Errorf("flow failed: %w", err)
	}

	ui.PrintStatsTable(map[string]string{
		locale.T("messages.labels.files"):   ui.FormatNumber(ctx.Stats.Files),
		locale.T("messages.labels.bytes"):   ui.FormatBytes(ctx.Stats.Bytes),
		locale.T("messages.labels.elapsed"): time.Since(start).Round(time.Millisecond).String(),
	})
	ui.Success(locale.T("messages.success.flow_completed"))

	return nil
}

func parseYamlPipeline(filepath string) (*files.Pipeline, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var p files.Pipeline
	err = yaml.Unmarshal(data, &p)
	return &p, err
}

func parseInlinePipeline(inline string) (*files.Pipeline, error) {
	p := &files.Pipeline{Env: make(map[string]string), Steps: make([]files.PipelineStep, 0)}
	commands := strings.Split(inline, ";")

	for _, cmd := range commands {
		cmd = strings.TrimSpace(cmd)
		if cmd == "" {
			continue
		}

		parts := strings.Fields(cmd)
		if len(parts) == 0 {
			continue
		}

		step := files.PipelineStep{Action: parts[0], Args: make(map[string]string)}

		if step.Action == "var" {
			if len(parts) >= 2 {
				eqIdx := strings.Index(parts[1], "=")
				if eqIdx != -1 {
					step.Args["name"] = parts[1][:eqIdx]
					step.Args["value"] = parts[1][eqIdx+1:]
				} else if len(parts) >= 3 {
					step.Args["name"] = parts[1]
					step.Args["value"] = parts[2]
				}
			}
			p.Steps = append(p.Steps, step)
			continue
		}

		for i := 1; i < len(parts); i++ {
			arg := parts[i]
			if arg == "-d" || arg == "--dest" {
				if i+1 < len(parts) {
					step.Args["dest"] = parts[i+1]
					i++
				}
			} else if arg == "-w" || arg == "--rewrite" {
				step.Args["rewrite"] = "true"
			} else if arg == "--to" {
				if i+1 < len(parts) {
					step.Args["to"] = parts[i+1]
					i++
				}
			} else if arg == "-o" || arg == "--output" {
				if i+1 < len(parts) {
					step.Args["output"] = parts[i+1]
					i++
				}
			} else if arg == "--exec" {
				if i+1 < len(parts) {
					step.Args["exec"] = parts[i+1]
					i++
				}
			} else if arg == "--search" {
				if i+1 < len(parts) {
					step.Args["search"] = parts[i+1]
					i++
				}
			} else if arg == "--replace" {
				if i+1 < len(parts) {
					step.Args["replace"] = parts[i+1]
					i++
				}
			} else if arg == "--prefix" {
				if i+1 < len(parts) {
					step.Args["prefix"] = parts[i+1]
					i++
				}
			} else if arg == "--suffix" {
				if i+1 < len(parts) {
					step.Args["suffix"] = parts[i+1]
					i++
				}
			}
		}
		p.Steps = append(p.Steps, step)
	}
	return p, nil
}
