package main

import (
	"flag"
	"os"
	"strconv"
	"sync"
)

const N = 20

func main() {
	doLockFlag := flag.Bool("lock", false, "use a mutex to lock Write calls")
	flag.Parse()
	doLock := *doLockFlag
	var l sync.Mutex
	var wg sync.WaitGroup
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func(id int) {
			// Write a unique line using a single Write call
			msg := "Goroutine " + strconv.Itoa(id) + "\n"
			for range 10000 {
				if doLock {
					l.Lock()
				}
				os.Stdout.Write([]byte(msg))
				if doLock {
					l.Unlock()
				}
			}
			wg.Done()
		}(i)
	}

	wg.Wait()
}
