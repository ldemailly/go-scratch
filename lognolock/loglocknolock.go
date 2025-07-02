/*
  Evaluate the effect on performance of lock protecting Write calls (which log and slog do).
  Test correctness:
    go run . | sort | uniq -c
  Test performance:
    go build .
    hyperfine --warmup 3 --runs 10  "./lognolock > /dev/null"
  vs
    hyperfine --warmup 3 --runs 10  "./lognolock  -lock > /dev/null"

  On macos m3 pro: with lock is ~150ms and faster than without lock ~200ms. (!!)
*/

package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"sync"
)

const defaultGoRoutines = 20
const defaultWrites = 10000

func main() {
	doLockFlag := flag.Bool("lock", false, "use a mutex to lock Write calls")
	numRoutinesFlag := flag.Int("n", defaultGoRoutines, "number of goroutines to run")
	numWritesFlag := flag.Int("w", defaultWrites, "number of writes per goroutine")
	flag.Parse()
	numWrites := *numWritesFlag
	doLock := *doLockFlag
	var l sync.Mutex
	var wg sync.WaitGroup
	fmt.Fprintf(os.Stderr, "Running %d goroutines with %d Write in each; Using lock: %t\n", *numRoutinesFlag, numWrites, doLock)
	for i := 0; i < *numRoutinesFlag; i++ {
		wg.Add(1)
		go func(id int) {
			// Write a unique line using a single Write call
			msg := "Goroutine " + strconv.Itoa(id) + "\n"
			for range numWrites {
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
