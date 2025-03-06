// Example of simple progress bar in Go + printing something above it too
package main

import (
	"fmt"
	"strings"
	"time"
)

const (
	width = 60
)

// Show a progress bar from 0-1 (0-100%).
func ProgressBar(progress float64) {
	count := int(progress*width + 0.5)
	bar := "[" + strings.Repeat("█", count) + strings.Repeat(" ", width-count) + "]"
	fmt.Printf("\r%s %.2f%%", bar, progress*100)
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
	for i := 0; i <= 1000; i++ {
		ProgressBar(float64(i) / 1000)
		if i%63 == 0 {
			MoveCursorUp(1)
			fmt.Printf("%d of 1000\n", i)
		}
		time.Sleep(10 * time.Millisecond)
	}
	fmt.Println()
}
