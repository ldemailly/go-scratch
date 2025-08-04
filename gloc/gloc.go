package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"strings"

	"fortio.org/cli"
	"fortio.org/log"
)

// Count the lines of code in Go files, excluding comments and empty lines.

func main() {
	cli.MinArgs = 1
	cli.ArgsHelp = "*.go files to count effective lines of code"
	debugFlag := flag.Bool("debug", false, "Enable debug lines output")
	cli.Main()
	debug := *debugFlag
	totalLines := 0
	errCount := 0
	// Create a new file set for AST parsing
	fset := token.NewFileSet()
	for _, filename := range flag.Args() {
		lines, err := countEffectiveLinesOfCode(debug, filename, fset)
		if err != nil {
			log.Errf("Error processing %s: %v", filename, err)
			errCount++
			continue
		}
		fmt.Printf("%s: %d\n", filename, lines)
		totalLines += lines
	}
	fmt.Printf("Total: %d lines of code\n", totalLines)
	if errCount > 0 {
		log.Errf("%d files could not be processed due to errors.", errCount)
		os.Exit(1)
	}
}

func countEffectiveLinesOfCode(debug bool, filename string, fset *token.FileSet) (int, error) {
	// Parse the file to get comments (we already read the content)
	// without parser.ParseComments to skip comments
	file, err := parser.ParseFile(fset, filename, nil, parser.SkipObjectResolution)
	if err != nil {
		return 0, fmt.Errorf("failed to parse file: %w", err)
	}
	var buf bytes.Buffer
	originalLines := fset.Position(file.Pos()).Line - 1
	log.Debugf("Processing %s, original lines: %d", filename, originalLines)

	err = printer.Fprint(&buf, fset, file)
	if err != nil {
		panic(err)
	}
	// rescan the file to get the effective lines of code
	scanner := bufio.NewScanner(&buf)
	lines := 0
	for scanner.Scan() {
		originalLines++
		line := scanner.Text()
		if strings.TrimSpace(line) != "" {
			lines++
			if debug {
				fmt.Printf("Line %3d (%3d): \t%s\n", lines, originalLines, line)
			}
		}
	}
	return lines, nil
}
