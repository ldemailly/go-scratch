package main

import (
	"bufio"
	"flag"
	"io"
	"math"
	"os"

	"github.com/ldemailly/go-scratch/iobench/blockio"
)

func main() {
	var filename string
	var mode string
	var bufferSize int
	flag.StringVar(&filename, "filename", "genfile/random_lines.txt", "Output filename")
	flag.StringVar(&mode, "mode", "blockio", "Mode to run: blockio, scanner, ...")
	flag.IntVar(&bufferSize, "bufferSize", 64*1024, "Buffer size")
	flag.Parse()
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	switch mode {
	case "blockio":
		processBlockio(f, bufferSize)
	case "scanner":
		processScanner(f, bufferSize)
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
	println("BLOCKIO Lines:", lines, "minV:", minV, "Max:", maxV, "Total:", total)
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
	println("BUFIO Scanner Lines:", lines, "minV:", minV, "Max:", maxV, "Total:", total)
}
