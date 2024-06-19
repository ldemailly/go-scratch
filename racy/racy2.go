package main

import (
	"sync"
	"fmt"
	"flag"
)

func main() {
	n := flag.Int("n", 100000, "number of increments (should be even)")
        flag.Parse()
	shared := 0
	var ready, start, finished sync.WaitGroup
	ready.Add(2)
	start.Add(1)
	finished.Add(2)
	v := *n / 2
	for range 2 {
		go func() {
			ready.Done()
			start.Wait()
			for range v {
				shared++
			}
			finished.Done()
		}()
	}
	ready.Wait()
	start.Done()
	finished.Wait()
	fmt.Printf("expected %d, got %d\n", 2*v, shared)
}
