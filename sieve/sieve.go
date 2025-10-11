// Super minimal demo of ansipixels with Eratosthenes' sieve for prime numbers
package main

import (
	"bytes"
	"flag"
	"fmt"
	"math"

	"fortio.org/terminal/ansipixels"
	"fortio.org/terminal/ansipixels/tcolor"
)

// returns the N and width of each cell and how many cell per line - eg 700 and 4 (3+1) and p, or 2300 and 5 (4+1) and q.
func calcN(height, width int) (int, int, int) {
	// Find number of digits (start with 3 and try 4)
	spaceUsedPer := 3 + 1 // + 1 space in between
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

type State struct {
	ap        *ansipixels.AnsiPixels
	n         int
	padding   int
	perLine   int
	current   int
	multiple  int
	state     []bool // index is n-2
	numPrimes int
	angle     float64
	chroma    float64
	lightness float64
}

func (s *State) NumberAt(n int) (int, int) {
	return ((n - 1) % s.perLine) * s.padding, (n - 1) / s.perLine
}

// InitialState optimizes the initial state where we can write it
// all in 1 string instead of cursor moving to each position.
func (s *State) InitialState() error {
	s.ap.WriteString(tcolor.Reset)
	s.ap.ClearScreen()
	// Calculate and save sieve dimensions:
	// leave the last line empty so on exit the whole sieve is visible.
	s.n, s.padding, s.perLine = calcN(s.ap.H-1, s.ap.W)
	// Reset state
	s.current = 1
	s.multiple = 0
	s.state = make([]bool, s.n-1) // all false == all possible primes. (starting at 2)
	var buf bytes.Buffer
	for i := 1; i <= s.n; i++ {
		fmt.Fprintf(&buf, "%*d", (s.padding - 1), i)
		if i == s.n {
			break // so no extra newline scrolling it all at the end/last n
		}
		if i%s.perLine == 0 {
			buf.WriteString("\r\n")
		} else {
			buf.WriteByte(' ')
		}
	}
	_, _ = s.ap.Out.Write(buf.Bytes())
	_ = s.ap.Out.Flush()
	return nil
}

func (s *State) Color() tcolor.Color {
	hue := float64(s.numPrimes) * s.angle // trying to get nice steps.
	hue -= float64(int(hue))
	return tcolor.Oklchf(s.lightness, s.chroma, hue)
}

func (s *State) ShowNumberAt(n int, prefix string) {
	x, y := s.NumberAt(n)
	s.ap.WriteAt(x, y, "%s%*d ", prefix, s.padding-1, n)
}

func (s *State) DemoColor(n int) {
	s.ShowNumberAt(n, s.ap.ColorOutput.Foreground(tcolor.HSLf(float64(n)/float64(s.n+1), 0.5, 0.5)))
}

func (s *State) IsFlagged(n int) bool {
	return s.state[n-2]
}

func (s *State) Flag(n int) {
	s.state[n-2] = true
}

func main() { //nolint: gocognit // a bit big one function demo indeed.
	fps := flag.Float64("fps", 120.0, "Frames per second")
	palette := flag.Bool("palette", false, "Just show the palette")
	flag.Parse()
	ap := ansipixels.NewAnsiPixels(*fps)
	err := ap.Open()
	if err != nil {
		panic(err)
	}
	defer ap.Restore()
	// Golden angle https://en.wikipedia.org/wiki/Golden_angle
	// 2-phi == 1 - 1/phi == 1/phi^2 (!)
	s := &State{ap: ap, angle: 2 - math.Phi, chroma: 0.5, lightness: 0.7}
	if *palette {
		for _, chroma := range []float64{0.4, 0.45, 0.5, 0.55, 0.6} {
			s.chroma = chroma
			for _, lightness := range []float64{0.4, 0.5, 0.6, 0.7} {
				s.lightness = lightness
				for _, angle := range []struct {
					name  string
					angle float64
				}{{"Golden angle", 2 - math.Phi}} {
					s.angle = angle.angle
					fmt.Printf("Palette with chroma %.1f lightness %.1f angle %s %f\r\n", s.chroma, s.lightness, angle.name, angle.angle)
					for i := range 39 {
						n := i + 1
						s.numPrimes = i
						fmt.Printf("%s%3d ", s.ap.ColorOutput.Background(s.Color()), n)
						if n%13 == 0 {
							fmt.Printf("%s\r\n", tcolor.Reset)
						}
					}
					fmt.Printf("%s\r\n", tcolor.Reset)
				}
			}
		}
		return
	}
	ap.OnResize = s.InitialState
	_ = ap.OnResize()
	ap.WriteBoxed(ap.H/2, " Resize me, Q or ^C to quit \n any key to start, fps: %.0f ", *fps)
	_ = ap.ReadOrResizeOrSignal()
	_ = ap.OnResize() // redraw without the box
	frame := 0
	stoppedEarly := true
	err = ap.FPSTicks(func() bool {
		if len(s.ap.Data) > 0 {
			c := ap.Data[0]
			switch c {
			case 'q', 3, 'Q':
				return false // exit requested
			default:
				// nothing else to do for now
			}
		}
		if s.current*s.current >= s.n {
			stoppedEarly = false
			return false // all done
		}
		// Either we're marking multiples of a prime or we find the next one:
		if s.multiple == 0 { // find next one mode
			candidate := s.current + 1
			for candidate <= s.n && s.IsFlagged(candidate) {
				candidate++
			}
			s.numPrimes++
			s.lightness = 0.75
			color := s.Color()
			s.ap.WriteAtStr(0, ap.H-1, tcolor.Reset)
			s.ap.ClearEndOfLine()
			s.ap.WriteCentered(ap.H-1, "Next Prime found: %d", candidate)
			s.ShowNumberAt(candidate, s.ap.ColorOutput.Background(color)+tcolor.Bold+tcolor.Underlined)
			s.current = candidate
			s.multiple = candidate * candidate
			return true
		}
		// next multiple marking:
		s.lightness = 0.7
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
			s.ap.WriteCentered(ap.H-1, "%sMarking multiples of %d - marking %d%s",
				tcolor.Reset, s.current, s.multiple, tcolor.Black.Foreground())
			s.ShowNumberAt(s.multiple, s.ap.ColorOutput.Background(color))
			return true // one at a time
		}
		s.multiple = 0 // done with this prime, back to find next one mode
		return true
	})
	if err != nil {
		panic(err)
	}
	if stoppedEarly {
		ap.WriteCentered(ap.H, "%sStopped while marking %d, press a key to list primes found and exit ...", tcolor.Reset, s.current)
	} else {
		ap.WriteCentered(ap.H, "%sAll done (after %d), press a key to list primes found and exit ...", tcolor.Reset, s.current)
	}
	_ = s.ap.ReadOrResizeOrSignal() // pause at the end
	ap.MoveCursor(0, ap.H)
	ap.SaveCursorPos()
	ap.ClearEndOfLine()
	fmt.Fprintln(ap.Logger, "Primes found:")
	for i := 2; i <= min(s.n, s.current*s.current); i++ {
		if !s.IsFlagged(i) {
			fmt.Fprintf(ap.Logger, "%d ", i)
		}
	}
	fmt.Fprintln(ap.Logger)
}
