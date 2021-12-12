package main

import (
	"github.com/fatih/color"
	"io"
)

type PrefixWriter struct {
	color      *color.Color
	script     GenerationScript
	underlying io.Writer
}

func (w *PrefixWriter) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	i := 0
	j := 0
	for j < len(p) {
		if p[j] == '\n' {
			// Write nicely formatted prefix
			if _, err := w.color.Fprintf(w.underlying, "[%s]", w.script.Name); err != nil {
				return 0, err
			}
			if _, err := w.underlying.Write([]byte{' '}); err != nil {
				return 0, err
			}

			// Write actual passed data
			_, err = w.underlying.Write(p[i : j+1])
			if err != nil {
				return 0, err
			}

			i = j + 1
		}
		j++
	}

	// Check if we need a newline at the end
	// This may break the formatting, but it prevents multiple writers from causing issues.
	if p[len(p)-1] != '\n' {
		_, err = w.underlying.Write([]byte{'\n'})
		if err != nil {
			return 0, err
		}
	}

	return len(p), nil
}

func NewPrefixWriter(script GenerationScript, underlying io.Writer, color *color.Color) *PrefixWriter {
	return &PrefixWriter{
		script:     script,
		underlying: underlying,
		color:      color,
	}
}
