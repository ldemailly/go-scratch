package main

import (
	"flag"
	"fmt"
	"math/rand/v2"
	"os"
)

// Spaces 4 times more frequent... just to pretend to make it more readable.
const charset = " abcdefghijklmnopqrstuvwxyz ABCDEFGHIJKLMNOPQRSTUVWXYZ 0123456789_. "

func main() {
	var minLen int
	var maxLen int
	var lineCount int
	var filename string
	flag.IntVar(&minLen, "minLen", 7, "Minimum length of a line")
	flag.IntVar(&maxLen, "maxLen", 250, "Maximum length of a line")
	flag.IntVar(&lineCount, "lineCount", 2_000_000, "Number of lines")
	flag.StringVar(&filename, "filename", "random_lines.txt", "Output filename")
	flag.Parse()

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()
	buf := make([]byte, maxLen+1)
	totalSize := int64(0)
	for range lineCount {
		length := rand.IntN(maxLen-minLen+1) + minLen //nolint: gosec // not crypto rand.
		randomString(buf, length)
		buf[length] = '\n'
		n, _ := file.Write(buf[:length+1])
		totalSize += int64(n)
	}
	fmt.Printf("File generated: %q, %d lines, total size %d", filename, lineCount, totalSize)
}

func randomString(buf []byte, length int) {
	b := buf[:length]
	for i := range b {
		b[i] = charset[rand.IntN(len(charset))] //nolint: gosec // not crypto rand.
	}
}
