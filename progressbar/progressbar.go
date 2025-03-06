// Example of simple progress bar in Go + printing something above it too
package main

import (
	"fmt"
	"strings"
	"time"
)

const (
	width = 40
	Space = " "
	Full  = "█"
)

// 1/8th of a full block to 7/8ths of a full block (ie fractional part of a block to
// get 8x resolution per character).
var FractionalBlocks = [...]string{"", "▏", "▎", "▍", "▌", "▋", "▊", "▉"}

// Show a progress bar from 0-1 (0-100%).
func ProgressBar(progress float64) {
	count := int(8*progress*width + 0.5)
	fullCount := count / 8
	remainder := count % 8
	spaceCount := width - fullCount - 1
	if remainder == 0 {
		spaceCount++
	}
	bar := "[" + strings.Repeat(Full, fullCount) + FractionalBlocks[remainder] + strings.Repeat(Space, spaceCount) + "]"
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
	fmt.Print("Progress bar example\n\n")
	n := width * 8 // just to demo every smooth step
	for i := 0; i <= n; i++ {
		ProgressBar(float64(i) / float64(n))
		if i%63 == 0 {
			MoveCursorUp(1)
			fmt.Printf("Just an extra demo print for %d\n", i)
		}
		time.Sleep(20 * time.Millisecond)
	}
	fmt.Println()
}
