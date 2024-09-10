package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand/v2"
	"os"
)

// Spaces 4 times more frequent... just to pretend to make it more readable.
// also have it have 64 characters, (could) make the randomString() faster.
const charset = " abcdefghijkmnopqrstuvwxyz ABCDEFGHJKLMNOPQRSTUVWXYZ 0123456789 "

func main() {
	var minLen int
	var maxLen int
	var lineCount int
	var filename string
	flag.IntVar(&minLen, "minLen", 7, "Minimum length of a line")
	flag.IntVar(&maxLen, "maxLen", 250, "Maximum length of a line")
	flag.IntVar(&lineCount, "lineCount", 10_000_000, "Number of lines")
	flag.StringVar(&filename, "filename", "random_lines.txt", "Output filename")
	flag.Parse()
	fmt.Println("charset len:", len(charset))
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()
	f := bufio.NewWriter(file)
	buf := make([]byte, maxLen+1)
	totalSize := int64(0)
	for range lineCount {
		length := rand.IntN(maxLen-minLen+1) + minLen //nolint: gosec // not crypto rand.
		randomString(buf, length)
		buf[length] = '\n'
		n, _ := f.Write(buf[:length+1])
		totalSize += int64(n)
	}
	f.Flush()
	fmt.Printf("File generated: %q, %d lines, total size %d\n", filename, lineCount, totalSize)
}

func randomString(buf []byte, length int) {
	for i := range length {
		buf[i] = charset[rand.IntN(len(charset))] //nolint: gosec // not crypto rand.
	}
}
