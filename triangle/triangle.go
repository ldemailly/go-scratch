package main

import (
	"bytes"
	"flag"
	"os"
)

func triangle(buf []byte, n int) []byte {
	row := bytes.Repeat([]byte{' '}, n*2+2)
	row[0] = '|'
	row[n*2] = '|'
	row[n*2+1] = '\n'
	row[n] = '#'
	for i := range n {
		buf = append(buf, row...)
		row[n-i-1] = '#'
		row[n+i+1] = '#'
	}
	return buf
}

func main() {
	var n int
	flag.IntVar(&n, "n", 5, "number of rows")
	flag.Parse()
	buf := make([]byte, 0, (n*2+2)*n)
	buf = triangle(buf, n)
	os.Stdout.Write(buf)
}
