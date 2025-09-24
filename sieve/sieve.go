// Super minimal demo of ansipixels with Eratosthenes' sieve for prime numbers
package main

import (
	"flag"

	"fortio.org/terminal/ansipixels"
	"fortio.org/terminal/ansipixels/tcolor"
)

// returns the N and width of each cell and how many cell per line - eg 700 and 4 (3+1) and p, or 2300 and 5 (4+1) and q.
func calcN(height, width int) (int, int, int) {
	// Find number of digits (start with 3 and try 4)
	spaceUsedPer := 3 + 1                 // + 1 space in between
	perLine := (width + 1) / spaceUsedPer // + 1 because last space is not needed
	total := perLine * height
	if total < 1000 {
		return total, spaceUsedPer, perLine
	}
	// try 4
	spaceUsedPer++
	nPerLine := (width + 1) / spaceUsedPer
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

type State struct {
	ap      *ansipixels.AnsiPixels
	n       int
	padding int
	perLine int
}

func (s *State) NumberAt(n int) (int, int) {
	return ((n - 1) % s.perLine) * s.padding, (n - 1) / s.perLine
}

func (s *State) ShowNumberAt(n int, c tcolor.Color) {
	x, y := s.NumberAt(n)
	s.ap.WriteAt(x, y, "%s%d", s.ap.ColorOutput.Foreground(c), n)
}

func main() {
	fps := flag.Float64("fps", 60.0, "Frames per second")
	flag.Parse()
	ap := ansipixels.NewAnsiPixels(*fps)
	err := ap.Open()
	if err != nil {
		panic(err)
	}
	defer ap.Restore()
	s := &State{ap: ap}
	ap.OnResize = func() error {
		ap.ClearScreen()
		// Calculate and save sieve dimensions:
		s.n, s.padding, s.perLine = calcN(ap.H, ap.W)
		for i := 1; i <= s.n; i++ {
			s.ShowNumberAt(i, tcolor.HSLf(float64(i)/float64(s.n+1), 0.5, 0.5))
		}
		return nil
	}
	_ = ap.OnResize()
	for {
		ap.WriteBoxed(ap.H/2, "%sResize me, Q or ^C to quit", tcolor.Reset)
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
