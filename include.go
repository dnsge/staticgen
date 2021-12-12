package main

import (
	cpy "github.com/otiai10/copy"
	"os"
	"path"
)

type IncludeSpec struct {
	// Name is the included item's representation in staticgen.json
	Name string
	// FullPath is the absolute path to the included item (file, directory)
	FullPath string
	// OutputPath is the absolute path to the corresponding output location
	OutputPath string
	// IsDirectory describes whether the included item is a directory
	IsDirectory bool
}

func (i *IncludeSpec) Copy() error {
	return cpy.Copy(i.FullPath, i.OutputPath)
}

// makeIncludeSpecs takes in the working directory, output directory, and a slice of
// relative include paths and generates the corresponding IncludeSpec entries.
func makeIncludeSpecs(pwd string, outPath string, includes []string) ([]IncludeSpec, error) {
	includeSpecs := make([]IncludeSpec, len(includes))
	for i, include := range includes {
		includeInPath := makeAbs(pwd, include)

		// check whether directory or file
		s, err := os.Stat(includeInPath)
		if err != nil {
			return nil, err
		}

		includeSpecs[i] = IncludeSpec{
			Name:        include,
			FullPath:    includeInPath,
			OutputPath:  path.Join(outPath, s.Name()),
			IsDirectory: s.IsDir(),
		}
	}

	return includeSpecs, nil
}

func CopyAll(specs []IncludeSpec) error {
	for _, spec := range specs {
		if err := spec.Copy(); err != nil {
			return err
		}
	}
	return nil
}
