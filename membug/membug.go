package main

import (
	"fmt"
	"log"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

func FreeMemory() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	currentAlloc := memStats.HeapAlloc
	gomemlimit := debug.SetMemoryLimit(-1)

	fmt.Printf("*** Current HeapAlloc: %d bytes\n", currentAlloc)
	fmt.Printf("*** Usage percentage: %.2f%%\n", (float64(currentAlloc)/float64(gomemlimit))*100)
}

// go build . && GOMEMLIMIT=30MiB GODEBUG=gcstoptheworld=1,gctrace=1 ./membug
func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Got expected panic:", r)
		}
	}()
	oldLimit := debug.SetMemoryLimit(-1)
	newLimit := int64(20_000_000)
	l := debug.SetMemoryLimit(newLimit)
	l = debug.SetMemoryLimit(newLimit) // twice to confirm we did change something
	runtime.GC()
	time.Sleep(1 * time.Second)
	FreeMemory()
	log.Printf("### Post initial GC, Old limit: %d, new %d", oldLimit, l)
	v1 := strings.Repeat("ABC", 4*int(newLimit)) // ie 12x the limit
	runtime.GC()
	time.Sleep(1 * time.Second)
	FreeMemory()
	log.Printf("### 1. post 2nd GC no panic... len: %d", len(v1))
	a1 := make([]byte, 10*int(newLimit)) // now at 22x the limit
	runtime.GC()
	time.Sleep(1 * time.Second)
	FreeMemory()
	log.Printf("### 2. post 3rd GC no panic... a1 cap: %d", cap(a1))
	copy(a1, []byte(v1))
	runtime.GC()
	time.Sleep(1 * time.Second)
	log.Printf("### 3. post 4th GC, no panic... len v1: %d, len a1:%d", len(v1), len(a1))
	v2 := strings.Repeat("BAC", 4*int(newLimit)) // now 34x the limit and still going...
	log.Printf("### 4.  no panic... len: %d", len(v2))
	/*	//	os.Stdout.Write([]byte(v1)) // using v1 and v2 so it has to exist in memory
		//	os.Stdout.Write([]byte(v2))
		runtime.GC()
		log.Printf("Sleeping to allow for checking the process using ps etc..")
		time.Sleep(10 * time.Second)
		v3 := v1 + v2
		log.Printf("3. no panic... len: %d", len(v3))
		// os.Stdout.Write([]byte(v3))
	*/
	runtime.GC()
	time.Sleep(1 * time.Second)
	log.Printf("#### after last GC, sleeping to allow for checking the process using ps etc..")
	time.Sleep(30 * time.Second)
	log.Printf("%d %d %d", len(v1), len(a1), len(v2))
}
