package main

import (
	"errors"
	"io"
	"log"
	"os"
)

const bufSize = 1024

func main() {
	buf := make([]byte, bufSize)
	var err error
	var r, w int
	for {
		r, err = os.Stdin.Read(buf)
		if r > 0 {
			w, err = os.Stdout.Write(buf[:r])
			if err != nil {
				log.Fatalf("write error %d %v", w, err)
			}
			if w != r {
				log.Fatalf("write mismatch %d %d", w, r)
			}
		}
		if err != nil {
			break
		}
	}
	if !errors.Is(err, io.EOF) {
		log.Fatalf("read error: %v", err)
	}
}
