package main

import (
	"context"
	"github.com/fatih/color"
	"os"
	"os/exec"
	"strings"
	"time"
)

type GenerationScript struct {
	Name  string   `json:"name"`
	Watch []string `json:"watch"`
	Build []string `json:"build"`
}

func runGenerationScripts(pwd string, scripts []GenerationScript) {
	for _, script := range scripts {
		if len(script.Build) == 0 {
			continue
		}

		baseColor := getNextColor()
		printFunc := baseColor.PrintfFunc()

		// Create os command
		cmd := makeCommand(script.Build)
		cmd.Dir = pwd
		cmd.Stdout = NewPrefixWriter(script, os.Stdout, baseColor)
		cmd.Stderr = NewPrefixWriter(script, os.Stdout, color.New(color.FgHiWhite).Add(color.BgRed))

		printFunc("[%s] Generating ... ", script.Name)
		start := time.Now()
		if err := cmd.Run(); err != nil {
			printFunc("error executing %q: %v\n", strings.Join(script.Build, " "), err)
			os.Exit(1)
		} else {
			end := time.Now()
			printFunc("done (%v)\n", end.Sub(start).Round(time.Millisecond))
		}
	}
}

func makeCommand(cmd []string) *exec.Cmd {
	if len(cmd) == 1 {
		return exec.Command(cmd[0])
	} else {
		return exec.Command(cmd[0], cmd[1:]...)
	}
}

func makeCommandContext(ctx context.Context, cmd []string) *exec.Cmd {
	if len(cmd) == 1 {
		return exec.CommandContext(ctx, cmd[0])
	} else {
		return exec.CommandContext(ctx, cmd[0], cmd[1:]...)
	}
}
