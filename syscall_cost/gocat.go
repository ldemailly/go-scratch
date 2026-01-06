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

const bufSize = 1024

func fileCat() int64 {
	var buf [bufSize]byte
	var err error
	var r, w int
	var total int64
	for {
		r, err = os.Stdin.Read(buf[:])
		if r > 0 {
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
	return total
}

func syscallCat() int64 {
	var buf [bufSize]byte
	var err error
	var r, w int
	var total int64
	for {
		r, err = syscall.Read(0, buf[:])
		if err != nil || r < 0 {
			break
		}
		if r == 0 {
			return total
		}
		w, err = syscall.Write(1, buf[:r])
		if err != nil {
			log.Fatalf("write error %d %v", w, err)
		}
		if w != r {
			log.Fatalf("write mismatch %d %d", w, r)
		}
		total += int64(w)
	}
	if !errors.Is(err, io.EOF) {
		log.Fatalf("read error: %v", err)
	}
	return total
}

func main() {
	syscall := flag.Bool("syscall", false, "use syscall directly for I/O (file otherwise)")
	flag.Parse()
	var total int64
	if *syscall {
		total = syscallCat()
	} else {
		total = fileCat()
	}
	fmt.Fprintf(os.Stderr, "%d\n", total)
}
