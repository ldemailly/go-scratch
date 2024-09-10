package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math"
	"os"

	"github.com/ldemailly/go-scratch/iobench/blockio"
	"github.com/ldemailly/go-scratch/iobench/optio"
)

func main() {
	var filename, mode string
	var bufferSize, readSize, maxLineLength int
	flag.StringVar(&filename, "filename", "genfile/random_lines.txt", "Output filename")
	flag.StringVar(&mode, "mode", "optio", "Mode to run: blockio, scanner, optio")
	flag.IntVar(&bufferSize, "bufferSize", 1024, "Buffer size in kb")
	flag.IntVar(&readSize, "readSize", 128, "Optio unit of read size in kb")
	flag.IntVar(&maxLineLength, "maxline", 256, "Max line length in bytes")
	flag.Parse()
	bufferSize *= 1024
	readSize *= 1024
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	if len(flag.Args()) > 0 {
		mode = flag.Args()[0]
	}
	fmt.Println("Running", mode, "on", filename, "with buffer size", bufferSize)
	switch mode {
	case "blockio":
		processBlockio(f, bufferSize)
	case "scanner":
		processScanner(f, bufferSize)
	case "optio":
		processOptio(f, bufferSize, readSize, maxLineLength)
	default:
		panic("Unknown/unimplemented mode " + mode)
	}
}

func processBlockio(f io.Reader, bufferSize int) {
	rs := blockio.BuildRecordSource(f, bufferSize)
	minV := math.MaxInt
	maxV := 0
	total := 0
	lines := 0
	for {
		line := rs.NextLine()
		length := len(line)
		if length == 0 {
			break
		}
		lines++
		total += length
		if length < minV {
			minV = length
		}
		if length > maxV {
			maxV = length
		}
	}
	println("BLOCKIO Lines:\t", lines, "minV:", minV, "Max:", maxV, "Total:", total)
}

func processScanner(f io.Reader, bufferSize int) {
	s := bufio.NewScanner(f)
	s.Buffer(make([]byte, 0, bufferSize), bufferSize)
	minV := math.MaxInt
	maxV := 0
	total := 0
	lines := 0

	for s.Scan() {
		line := s.Bytes()
		length := len(line)
		lines++
		total += length
		if length < minV {
			minV = length
		}
		if length > maxV {
			maxV = length
		}
	}
	println("Scanner Lines:\t", lines, "minV:", minV, "Max:", maxV, "Total:", total)
}

func processOptio(f io.Reader, bufferSize int, readSize int, maxLineLength int) {
	minV := math.MaxInt
	maxV := 0
	total := 0
	lines := 0
	s := optio.LineScanner(f, bufferSize, readSize, maxLineLength)
	for !s.EOF() {
		line, err := s.Line()
		if err != nil {
			panic(err)
		}
		length := len(line)
		lines++
		total += length
		if length < minV {
			minV = length
		}
		if length > maxV {
			maxV = length
		}
	}
	println("OPTIO Lines:\t", lines, "minV:", minV, "Max:", maxV, "Total:", total)
}
