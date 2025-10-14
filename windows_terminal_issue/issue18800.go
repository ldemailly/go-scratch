// This was initially created as minimal repro of a hang observed with windows terminal
// with https://github.com/fortio/fps#fps : https://github.com/microsoft/terminal/issues/18800
// It turns out it also reproduces a slow down / issue with Ghostty
// https://github.com/ghostty-org/ghostty/discussions/9187
// Note that this doesn't (yet) reproduce the slow down over time, only that 1.1.3 is faster
// than 1.2.x/tip. It does repro an even bigger slow down if ran _after_ fps (with something
// left on screen, and a bit better after a clear screen).
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/term"
)

func RawPrintln(stuff ...any) {
	fmt.Print(stuff...)
	fmt.Print("\r\n")
}

func main() {
	numIter := 1_000_000
	timeEvery := 100_000
	flag.IntVar(&numIter, "n", numIter, "number of `iterations`")
	flag.IntVar(&timeEvery, "t", timeEvery, "dump timing every `t` iterations")
	flag.Parse()
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Println("Error setting terminal to raw mode:", err)
		return
	}
	defer func() {
		err = term.Restore(fd, oldState)
		fmt.Println("Terminal restored to original , err", err)
	}()
	RawPrintln("Terminal in raw mode - 'q' to exit early")
	RawPrintln("it should end on its own otherwise after ", numIter)
	buf := make([]byte, 16) // fits ansi arrow escape, unicode, etc
	requestCursorPos := []byte("\033[6n")
	expected := len(requestCursorPos)
	now := time.Now()
	for iter := 1; iter <= numIter; iter++ {
		nw, err := os.Stdout.Write(requestCursorPos)
		if err != nil {
			RawPrintln("Error writing to terminal:", err)
			return
		}
		if nw != expected {
			RawPrintln("Short write", nw, "vs", expected)
		}
		n, err := os.Stdin.Read(buf)
		if err != nil {
			RawPrintln("Error reading from terminal:", err)
			return
		}
		bufStr := string(buf[:n])
		// might fail with some ansi echo having a q in them,
		// but this is just a quick repro/test.
		if strings.ContainsRune(bufStr, 'q') {
			break
		}
		fmt.Printf("\r[%05d] Read %d bytes: %q      ", iter, n, bufStr)
		if iter%timeEvery == 0 {
			delta := time.Since(now)
			ops := float64(timeEvery) / delta.Seconds()
			delta = delta.Round(time.Millisecond)
			fmt.Printf("\t-- at %d took %v -- %.1f fps\r\n", iter, delta, ops)
			now = time.Now()
		}
	}
	RawPrintln("\r\nDone")
}
