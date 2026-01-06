package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"syscall"
)

const bufSize = 128 // small on purpose to get lots of syscalls

func fileCat() (int64, int) {
	var buf [bufSize]byte
	var err error
	var r, w int
	var total int64
	var numCall int
	for {
		numCall++
		r, err = os.Stdin.Read(buf[:])
		if r > 0 {
			numCall++
			w, err = os.Stdout.Write(buf[:r])
			if err != nil {
				log.Fatalf("write error %d %v", w, err)
			}
			if w != r {
				log.Fatalf("write mismatch %d %d", w, r)
			}
			total += int64(w)
		}
		if err != nil {
			break
		}
	}
	if !errors.Is(err, io.EOF) {
		log.Fatalf("read error: %v", err)
	}
	return total, numCall
}

func syscallCat() (int64, int) {
	var buf [bufSize]byte
	var err error
	var r int
	// var w int
	var total int64
	var numCall int
	for {
		numCall++
		r, err = syscall.Read(0, buf[:])
		if err != nil || r < 0 {
			break
		}
		if r == 0 {
			return total, numCall
		}
		/*
			numCall++
			w, err = syscall.Write(1, buf[:r])
			if err != nil {
				log.Fatalf("write error %d %v", w, err)
			}
			if w != r {
				log.Fatalf("write mismatch %d %d", w, r)
			}
		*/
		total += int64(r)
	}
	if !errors.Is(err, io.EOF) {
		log.Fatalf("read error: %v", err)
	}
	return total, numCall
}

func main() {
	syscall := flag.Bool("syscall", false, "use syscall directly for I/O (file otherwise)")
	flag.Parse()
	var total int64
	var numCall int
	if *syscall {
		total, numCall = syscallCat()
	} else {
		total, numCall = fileCat()
	}
	fmt.Fprintf(os.Stderr, "%d (%d)\n", total, numCall)
}
