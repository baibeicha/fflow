package files

import (
	"fmt"
	"strings"
	"sync"
)

// RunFlow executes a file processing pipeline.
func RunFlow(ctx *FlowContext, pipeline *Pipeline, onProgress func(action string, current, total int)) error {
	for k, v := range pipeline.Env {
		ctx.Vars[k] = interpolate(v, ctx.Vars)
	}

	for _, step := range pipeline.Steps {
		args := make(map[string]string)
		for k, v := range step.Args {
			args[k] = interpolate(v, ctx.Vars)
		}

		total := len(ctx.Files)
		progressCount := 0
		var mu sync.Mutex

		onProg := func() {
			mu.Lock()
			progressCount++
			cur := progressCount
			mu.Unlock()
			if onProgress != nil {
				onProgress(step.Action, cur, total)
			}
		}

		var err error
		var stats *TransferStats
		var newFiles []FileInfo

		switch step.Action {
		case "copy":
			destDirs := strings.Split(args["dest"], ",")
			rewrite := args["rewrite"] == "true"
			newFiles, stats, err = TransferFiles(ctx.Files, destDirs, false, rewrite, "Copy", onProg)
			if err != nil {
				return err
			}
			ctx.Stats.Add(stats)

		case "move":
			destDirs := strings.Split(args["dest"], ",")
			rewrite := args["rewrite"] == "true"
			newFiles, stats, err = TransferFiles(ctx.Files, destDirs, true, rewrite, "Copy", onProg)
			if err != nil {
				return err
			}
			ctx.Stats.Add(stats)
			ctx.Files = newFiles

		case "delete":
			newFiles, stats, err = DeleteFiles(ctx.Files, onProg)
			if err != nil {
				return err
			}
			ctx.Stats.Add(stats)
			ctx.Files = newFiles

		case "rename":
			newFiles, stats, err = RenameFiles(ctx.Files, args["search"], args["replace"], args["prefix"], args["suffix"], onProg)
			if err != nil {
				return err
			}
			ctx.Stats.Add(stats)
			ctx.Files = newFiles

		case "ext":
			newFiles, stats, err = ChangeExtension(ctx.Files, args["to"], onProg)
			if err != nil {
				return err
			}
			ctx.Stats.Add(stats)
			ctx.Files = newFiles

		case "zip":
			newFiles, stats, err = CreateZip(ctx.Files, args["output"], onProg)
			if err != nil {
				return err
			}
			ctx.Stats.Add(stats)
			ctx.Files = newFiles

		case "cmd":
			newFiles, stats, err = ExecuteCommandOnFiles(ctx.Files, args["exec"], onProg)
			if err != nil {
				return err
			}
			ctx.Stats.Add(stats)
			ctx.Files = newFiles

		case "var":
			ctx.Vars[args["name"]] = args["value"]
			if onProgress != nil {
				onProgress("var", 1, 1)
			}

		default:
			return fmt.Errorf("unknown flow action: %s", step.Action)
		}

		if len(ctx.Files) == 0 && step.Action != "var" && step.Action != "delete" {
			break
		}
	}

	return nil
}

func interpolate(val string, vars map[string]string) string {
	res := val
	for k, v := range vars {
		res = strings.ReplaceAll(res, "${"+k+"}", v)
	}
	return res
}
