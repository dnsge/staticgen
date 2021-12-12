package main

import (
	"github.com/fatih/color"
	"sync"
)

var (
	availableColors = []*color.Color{
		color.New(color.FgHiGreen),
		color.New(color.FgRed),
		color.New(color.FgHiBlue),
		color.New(color.FgCyan),
		color.New(color.FgHiYellow),
		color.New(color.FgHiCyan),
		color.New(color.FgHiWhite),
		color.New(color.FgHiRed),
	}
	currentColorIndex = 0

	mu sync.Mutex
)

func getNextColor() *color.Color {
	mu.Lock()
	defer mu.Unlock()

	ret := availableColors[currentColorIndex]
	currentColorIndex = (currentColorIndex + 1) % len(availableColors)
	return ret
}
