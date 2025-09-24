// Super minimal demo of ansipixels with Eratosthenes' sieve for prime numbers
package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"

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
	ap        *ansipixels.AnsiPixels
	n         int
	padding   int
	perLine   int
	current   int
	multiple  int
	state     []bool // index is n-1
	numPrimes int
}

func (s *State) NumberAt(n int) (int, int) {
	return ((n - 1) % s.perLine) * s.padding, (n - 1) / s.perLine
}

// Optimize the initial state where we can write it all in 1 string instead of cursor moving to each position
func (s *State) InitialState() error {
	s.ap.ClearScreen()
	// Calculate and save sieve dimensions:
	s.n, s.padding, s.perLine = calcN(s.ap.H, s.ap.W)
	// Reset state
	s.current = 1
	s.multiple = 0
	s.state = make([]bool, s.n) // all false == all possible primes.
	var buf bytes.Buffer
	for i := 1; i <= s.n; i++ {
		fmt.Fprintf(&buf, "%*d", -(s.padding - 1), i)
		if i == s.n {
			break // so no extra newline scrolling it all at the end/last n
		}
		if i%s.perLine == 0 {
			buf.WriteString("\r\n")
		} else {
			buf.WriteByte(' ')
		}
	}
	s.ap.Out.Write(buf.Bytes())
	s.ap.Out.Flush()
	return nil
}

func (s *State) Color() tcolor.Color {
	hue := float64(s.numPrimes) * 37.3 / 100. // just some random number of hue steps so colors are far enough apart
	// get the decimal part if greater than 1 / wrap
	hue -= float64(int(hue))
	return tcolor.HSLf(hue, 0.7, 0.4)
}

func (s *State) ShowNumberAt(n int, prefix string, suffix string) {
	x, y := s.NumberAt(n)
	s.ap.WriteAt(x, y, "%s%d%s", prefix, n, suffix)
}

func (s *State) DemoColor(n int) {
	s.ShowNumberAt(n, s.ap.ColorOutput.Foreground(tcolor.HSLf(float64(n)/float64(s.n+1), 0.5, 0.5)), tcolor.Reset)
}

func (s *State) IsFlagged(n int) bool {
	return s.state[n-1]
}
func (s *State) Flag(n int) {
	s.state[n-1] = true
}

func main() {
	fps := flag.Float64("fps", 120.0, "Frames per second")
	flag.Parse()
	ap := ansipixels.NewAnsiPixels(*fps)
	err := ap.Open()
	if err != nil {
		panic(err)
	}
	defer ap.Restore()
	s := &State{ap: ap}
	ap.OnResize = s.InitialState
	_ = ap.OnResize()
	ap.WriteBoxed(ap.H/2, " Resize me, Q or ^C to quit \n any key to start, fps: %.0f ", *fps)
	_ = ap.ReadOrResizeOrSignal()
	_ = ap.OnResize() // redraw without the box
	frame := 0
	err = ap.FPSTicks(context.Background(), func(_ context.Context) bool {
		if len(s.ap.Data) > 0 {
			c := ap.Data[0]
			switch c {
			case 'q', 3, 'Q':
				ap.MoveCursor(0, ap.H)
				return false
			default:
				// nothing else to do for now
			}
		}
		if s.current*s.current >= s.n {
			return true // all done just wait for resize or Quit
		}
		// Either we're marking multiples of a prime or we find the next one:
		if s.multiple == 0 { // find next one mode
			candidate := s.current + 1
			for candidate <= s.n && s.IsFlagged(candidate) {
				candidate++
			}
			s.numPrimes++
			color := s.Color()
			s.ShowNumberAt(candidate, s.ap.ColorOutput.Background(color), tcolor.Reset)
			s.current = candidate
			s.multiple = candidate * candidate
			return true
		}
		// next multiple marking:
		color := s.Color()
		frame++
		// slowdown based on current number (2 fastest, updates every frame, higher primer slower)
		if frame%(s.current-1) != 0 {
			return true
		}
		for ; s.multiple <= s.n; s.multiple += s.current {
			if s.IsFlagged(s.multiple) {
				continue // already marked
			}
			s.Flag(s.multiple)
			s.ShowNumberAt(s.multiple, s.ap.ColorOutput.Foreground(color), tcolor.Reset)
			return true // one at a time
		}
		s.multiple = 0 // done with this prime, back to find next one mode
		return true
	})
	if err != nil {
		panic(err)
	}
}
