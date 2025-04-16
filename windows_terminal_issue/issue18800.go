// https://github.com/microsoft/terminal/issues/18800

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

func RawPrintln(stuff ...any) {
	fmt.Print(stuff...)
	fmt.Print("\r\n")
}

func main() {
	numIter := 100_000
	flag.IntVar(&numIter, "n", numIter, "number of `iterations`")
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
	for iter := range numIter {
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
	}
	RawPrintln("\r\nDone")
}
