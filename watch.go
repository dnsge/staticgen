package main

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/rjeczalik/notify"
	"os"
	"path/filepath"
	"strings"
)

type IncludeWatcher struct {
	c           chan notify.EventInfo
	spec        IncludeSpec
	stdoutColor *color.Color
}

func NewIncludeWatcher(spec IncludeSpec) (*IncludeWatcher, error) {
	c := make(chan notify.EventInfo, 1)

	var err error
	if spec.IsDirectory {
		// If the spec describes a directory, recursively watch the whole directory
		err = notify.Watch(filepath.Join(spec.FullPath, "..."), c, notify.All)
	} else {
		// Otherwise, just watch the single file for only writes (renames, removes are invalid operations)
		err = notify.Watch(spec.FullPath, c, notify.Write)
	}

	if err != nil {
		return nil, err
	}

	return &IncludeWatcher{
		c:           c,
		spec:        spec,
		stdoutColor: getNextColor(),
	}, nil
}

func (iw *IncludeWatcher) Run(ctx context.Context) {
	printFunc := iw.stdoutColor.PrintfFunc()
	printFunc("[%s] Watching for changes\n", iw.spec.Name)
	for {
		select {
		case event := <-iw.c:
			if strings.HasSuffix(event.Path(), "~") {
				// Ignore any temporary files created by editors
				continue
			}

			printFunc("[%s] Syncing changes\n", iw.spec.Name)
			if err := iw.spec.Copy(); err != nil {
				fmt.Println(err)
			}
		case <-ctx.Done():
			notify.Stop(iw.c)
			close(iw.c)
			return
		}
	}
}

type GenerationWatcher struct {
	script      GenerationScript
	pwd         string
	stdoutColor *color.Color
}

func NewGenerationWatcher(script GenerationScript, pwd string) *GenerationWatcher {
	return &GenerationWatcher{
		script:      script,
		pwd:         pwd,
		stdoutColor: getNextColor(),
	}
}

func (gw *GenerationWatcher) Run(ctx context.Context) {
	cmd := makeCommandContext(ctx, gw.script.Watch)

	errColor := *gw.stdoutColor
	errColor.Add(color.BgRed)
	cmd.Stdout = NewPrefixWriter(gw.script, os.Stdout, gw.stdoutColor)
	cmd.Stderr = NewPrefixWriter(gw.script, os.Stderr, &errColor)
	cmd.Dir = gw.pwd
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func runWatch(ctx context.Context, pwd string, specs []IncludeSpec, scripts []GenerationScript) {
	for _, spec := range specs {
		watcher, err := NewIncludeWatcher(spec)
		if err != nil {
			panic(err)
		}

		go watcher.Run(ctx)
	}

	for _, genScript := range scripts {
		if len(genScript.Watch) == 0 {
			continue
		}

		watcher := NewGenerationWatcher(genScript, pwd)
		go watcher.Run(ctx)
	}

	<-ctx.Done()
}
