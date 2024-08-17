package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

// run > /dev/null to avoid flooding stdout
func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Got expected panic:", r)
		}
	}()
	oldLimit := debug.SetMemoryLimit(-1)
	newLimit := int64(20_000_000)
	l := debug.SetMemoryLimit(newLimit)
	log.Printf("Old limit: %d, new %d", oldLimit, l)
	v1 := strings.Repeat("A", 10*int(newLimit))
	log.Printf("1. no panic... len: %d", len(v1))
	v2 := strings.Repeat("B", 10*int(newLimit))
	log.Printf("2. no panic... len: %d", len(v2))
	os.Stdout.Write([]byte(v1)) // using v1 and v2 so it has to exist in memory
	os.Stdout.Write([]byte(v2))
	runtime.GC()
	log.Printf("Sleeping to allow for checking the process using ps etc..")
	time.Sleep(10 * time.Second)
	v3 := v1 + v2
	log.Printf("3. no panic... len: %d", len(v3))
	os.Stdout.Write([]byte(v3))
	log.Printf("Sleeping some more to allow for checking the process using ps etc..")
	time.Sleep(30 * time.Second)
}
