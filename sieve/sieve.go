// Super minimal demo of ansipixels with Eratosthenes' sieve for prime numbers
package main

import (
	"flag"
	"fmt"

	"fortio.org/terminal/ansipixels"
)

// returns the N and width of each cell and how many cell per line - eg 700 and 4 (3+1) and p, or 2300 and 5 (4+1) and q.
func calcN(height, width int) (int, int, int) {
	// Find number of digits (start with 3 and try 4)
	spaceUsedPer := 3 + 1 // + 1 space at least
	perLine := width / spaceUsedPer
	total := perLine * height
	if total < 1000 {
		return total, spaceUsedPer, perLine
	}
	// try 4
	spaceUsedPer++
	nPerLine := width / spaceUsedPer
	ntotal := nPerLine * height
	// could be back to 3 digits on border cases
	if ntotal < 1000 {
		return 999, spaceUsedPer - 1, perLine
	}
	// if somehow huge terminal that can do 5 digits... use 4 digits anyway for now
	if ntotal >= 10000 {
		return 9999, spaceUsedPer, nPerLine
	}
	return ntotal, spaceUsedPer, nPerLine
}

func main() {
	fps := flag.Float64("fps", 30.0, "Frames per second")
	ap := ansipixels.NewAnsiPixels(*fps)
	err := ap.Open()
	if err != nil {
		panic(err)
	}
	defer ap.Restore()
	ap.OnResize = func() error {
		ap.ClearScreen()
		// Calculate sieve dimensions:
		n, space, perLine := calcN(ap.H, ap.W)
		for i := 1; i <= n; i++ {
			if i != 1 && (i-1)%perLine == 0 {
				fmt.Fprintln(ap.Out, "\r")
			}
			fmt.Fprintf(ap.Out, "%*d", space, i)
		}
		return nil
	}
	_ = ap.OnResize()
	for {
		ap.WriteBoxed(ap.H/2, "Resize me, Q or ^C to quit")
		err := ap.ReadOrResizeOrSignal()
		if err != nil {
			panic(err)
		}
		if len(ap.Data) > 0 {
			c := ap.Data[0]
			switch c {
			case 'q', 3, 'Q':
				return
			default:
				// nothing else to do for now
			}
		}
	}
}
