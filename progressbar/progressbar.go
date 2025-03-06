// Example of simple progress bar in Go + printing something above it too
package main

import (
	"flag"
	"fmt"
	"strings"
	"time"
)

const (
	width = 40
	Space = " "
	Full  = "█"
	// Green FG, Grey BG
	Color = "\033[32;47m"
	Reset = "\033[0m"
)

// 1/8th of a full block to 7/8ths of a full block (ie fractional part of a block to
// get 8x resolution per character).
var FractionalBlocks = [...]string{"", "▏", "▎", "▍", "▌", "▋", "▊", "▉"}

// Show a progress bar from 0-1 (0-100%).
func ProgressBar(progress float64, useColors bool) {
	count := int(8*progress*width + 0.5)
	fullCount := count / 8
	remainder := count % 8
	spaceCount := width - fullCount - 1
	if remainder == 0 {
		spaceCount++
	}
	color := Color
	reset := Reset
	if !useColors {
		color = ""
		reset = ""
	}
	bar := color + strings.Repeat(Full, fullCount) + FractionalBlocks[remainder] + strings.Repeat(Space, spaceCount) + reset
	fmt.Printf("\r%s %.1f%%", bar, progress*100)
}

func MoveCursorUp(n int) {
	// ANSI escape codes used:
	// xA = move up x lines
	// 2K = clear entire line
	// G = move to the beginning of the line
	fmt.Printf("\033[%dA\033[2K\033[G", n)
}

func main() {
	colorFlag := flag.Bool("color", false, "Use color in the progress bar")
	flag.Parse()
	fmt.Print("Progress bar example\n\n")
	n := width * 8 // just to demo every smooth step
	for i := 0; i <= n; i++ {
		ProgressBar(float64(i)/float64(n), *colorFlag)
		if i%63 == 0 {
			MoveCursorUp(1)
			fmt.Printf("Just an extra demo print for %d\n", i)
		}
		time.Sleep(20 * time.Millisecond)
	}
	fmt.Println()
}
