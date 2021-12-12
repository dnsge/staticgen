package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// setupOutputDirectory creates the output directory if it doesn't exist and
// optionally cleans it by removing all files and directories within.
func setupOutputDirectory(outPath string, doClean bool) error {
	if s, err := os.Stat(outPath); os.IsNotExist(err) {
		err = os.MkdirAll(outPath, 0o755)
		if err != nil {
			return fmt.Errorf("create output directory: %w", err)
		}
	} else if !s.IsDir() {
		return fmt.Errorf("output is not a directory")
	} else if doClean {
		// remove all existing files + directories within output
		return os.RemoveAll(outPath)
	}

	return nil
}

func makeAbs(pwd string, p string) string {
	if filepath.IsAbs(p) {
		return p
	} else if strings.HasPrefix(p, "~/") {
		usr, _ := user.Current()
		dir := usr.HomeDir
		return filepath.Join(dir, p[2:])
	} else {
		return filepath.Join(pwd, p)
	}
}
