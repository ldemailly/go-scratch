// https://github.com/microsoft/terminal/issues/18800

package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

func RawPritnln(stuff ...any) {
	fmt.Print(stuff...)
	fmt.Print("\r\n")
}

func main() {
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
	RawPritnln("Terminal in raw mode - 'q' to exit")
	buf := make([]byte, 16) // fits ansi arrow escape, unicode, etc
	iter := 1
	for {
		requestCursorPosStr := "\033[6n"
		_, err = os.Stdout.Write([]byte(requestCursorPosStr))
		if err != nil {
			RawPritnln("Error writing to terminal:", err)
			return
		}
		n, err := os.Stdin.Read(buf)
		if err != nil {
			RawPritnln("Error reading from terminal:", err)
			return
		}
		bufStr := string(buf[:n])
		// might fail with some ansi echo having a q in them,
		// but this is just a quick repro/test.
		if strings.ContainsRune(bufStr, 'q') {
			RawPritnln("\r\nExiting...")
			break
		}
		fmt.Printf("\r[%05d] Read %d bytes: %q      ", iter, n, bufStr)
		iter++
	}
}
