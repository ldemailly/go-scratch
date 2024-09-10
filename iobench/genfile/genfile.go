package main

import (
	"flag"
	"fmt"
	"math/rand/v2"
	"os"
)

var (
	minLen    int
	maxLen    int
	lineCount int
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func init() {
	flag.IntVar(&minLen, "minLen", 7, "Minimum length of a line")
	flag.IntVar(&maxLen, "maxLen", 250, "Maximum length of a line")
	flag.IntVar(&lineCount, "lineCount", 2_000_000, "Number of lines")
}

func main() {
	flag.Parse()

	file, err := os.Create("random_lines.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	for i := 0; i < lineCount; i++ {
		length := rand.IntN(maxLen-minLen+1) + minLen
		line := randomString(length)
		file.WriteString(line + "\n")
	}

	fmt.Println("File generated: random_lines.txt")
}

func randomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.IntN(len(charset))]
	}
	return string(b)
}
