/*
  Evaluate the effect on performance of lock protecting Write calls (which log and slog do).
  Test correctness:
    go run . | sort | uniq -c
  Test performance:
    go build .
    hyperfine --warmup 3 --runs 10  "./lognolock > /dev/null"
  vs
    hyperfine --warmup 3 --runs 10  "./lognolock  -lock > /dev/null"

  Profiling:

  ./lognolock -n 16 -w 200_000 -cpuprofile nolock.pprof > /dev/null
  ./lognolock -n 16 -w 200_000 -cpuprofile lock.pprof -lock > /dev/null

  pprof -http :8001 -diff_base=nolock.pprof lock.pprof

  On macos m3 pro: with lock is ~150ms and faster than without lock ~200ms. (!!)
*/

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"strconv"
	"sync"
)

const defaultGoRoutines = 20
const defaultWrites = 10000

func main() {
	doLockFlag := flag.Bool("lock", false, "use a mutex to lock Write calls")
	numRoutinesFlag := flag.Int("n", defaultGoRoutines, "number of goroutines to run")
	numWritesFlag := flag.Int("w", defaultWrites, "number of writes per goroutine")
	pprofFile := flag.String("cpuprofile", "", "write cpu profile to file")
	flag.Parse()
	if *pprofFile != "" {
		f, err := os.Create(*pprofFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not create CPU profile: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Writing CPU profile to %s\n", *pprofFile)
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			fmt.Fprintf(os.Stderr, "could not start CPU profile: %v\n", err)
			os.Exit(1)
		}
		defer pprof.StopCPUProfile()
	}
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
