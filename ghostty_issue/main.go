package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"golang.org/x/term"
)

func main() {
	pFlag := flag.Bool("no-paste-mode", false, "Don't enable bracket paste mode")
	flag.Parse()
	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		fmt.Fprintln(os.Stderr, "Need a terminal.")
		os.Exit(1)
	}
	state, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error making terminal raw: %v\n", err)
		os.Exit(1)
	}
	defer term.Restore(fd, state)
	rw := struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stderr}
	t := term.NewTerminal(rw, "Test> ")
	if !*pFlag {
		t.SetBracketedPasteMode(true)
		fmt.Fprintln(t, "Bracketed paste mode enabled.")
	}
	fmt.Fprintf(t, "Enter/paste (multiline) text (type 'exit' to quit):\n")
	for {
		line, err := t.ReadLine()
		fmt.Fprintf(t, "Received line: %q\n", line)
		if errors.Is(err, io.EOF) {
			fmt.Fprintln(t, "EOF received, exiting.")
			break
		}
		if errors.Is(err, term.ErrPasteIndicator) {
			fmt.Fprintln(t, "Paste indicator received.")
			continue
		}
		if err != nil {
			panic(err)
		}
		if line == "exit" {
			fmt.Fprintln(t, "Exiting.")
			break
		}
	}
}
