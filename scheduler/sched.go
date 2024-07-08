package main

import (
	"flag"
	"runtime"

	"fortio.org/log"
	"fortio.org/scli"
)

func TightLoop(id, iterations int) {
	log.Infof("TightLoop %d %d started", id, iterations)
	var sum int64
	for i := range iterations {
		for j := range iterations {
			sum += int64(i + j)
		}
	}
	log.Infof("TightLoop %d %d ended", id, iterations)
}

func main() {
	p := flag.Int("n", 2, "GOMAXPROCS to use")
	scli.ServerMain()
	prev := runtime.GOMAXPROCS(*p)
	log.Infof("Previous GOMAXPROCS: %d, changed to %d", prev, *p)
	go TightLoop(1, 100000)
	go TightLoop(2, 100000)
	go TightLoop(3, 100000)
	scli.UntilInterrupted()
}
