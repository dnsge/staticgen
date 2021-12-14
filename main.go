package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

const (
	configFileName = "staticgen.json"
)

type Configuration struct {
	Clean   bool               `json:"clean"`
	Include []string           `json:"include"`
	Output  string             `json:"output"`
	Scripts []GenerationScript `json:"scripts"`
}

var (
	watchFlag = flag.Bool("watch", false, "Run watch generation configuration")
	debugFlag = flag.Bool("debug", false, "Print debug information")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "staticgen is a tool that generates file bundles and watches directories for changes\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Arguments:\n")
		flag.PrintDefaults()

		fmt.Fprintf(flag.CommandLine.Output(), "\nConfiguration:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "staticgen reads in staticgen.json in your current directory for configuration.\n")
		fmt.Fprintf(flag.CommandLine.Output(), "See https://github.com/dnsge/staticgen for documentation.")
	}

	flag.Parse()
}

func signalInterrupterContext() context.Context {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGKILL)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer func() {
			cancel()
		}()
		<-c
	}()

	return ctx
}

func main() {
	ctx := signalInterrupterContext()

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("get pwd error: %v\n", err)
		os.Exit(1)
	}

	f, err := os.Open(filepath.Join(pwd, configFileName))
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	var config Configuration
	if err := json.NewDecoder(f).Decode(&config); err != nil {
		fmt.Printf("configuration error: %v\n", err)
		os.Exit(1)
	}

	if *debugFlag {
		fmt.Printf("%#v\n", config)
	}

	if config.Output == "" {
		fmt.Printf("configuration error: output must be set\n")
		os.Exit(1)
	}

	outPath := makeAbs(pwd, config.Output)
	if err := setupOutputDirectory(outPath, config.Clean); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	includeSpecs, err := makeIncludeSpecs(pwd, outPath, config.Include)
	if err != nil {
		fmt.Printf("error reading includes: %v\n", err)
		os.Exit(1)
	}

	if err := CopyAll(includeSpecs); err != nil {
		fmt.Printf("error copying includes: %v\n", err)
		os.Exit(1)
	}

	if *watchFlag {
		runWatch(ctx, pwd, includeSpecs, config.Scripts)
	} else {
		runGenerationScripts(ctx, pwd, config.Scripts)
	}
}
