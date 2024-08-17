package main

import (
	"fmt"
	"runtime/debug"
	"strings"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Got expected panic:", r)
		}
	}()
	oldLimit := debug.SetMemoryLimit(-1)
	newLimit := int64(20_000_000)
	l := debug.SetMemoryLimit(newLimit)
	fmt.Println("Old limit:", oldLimit, "New limit:", l)
	v1 := strings.Repeat(".", 10*int(newLimit))
	fmt.Println("1. no panic... len:", len(v1))
	v2 := strings.Repeat(".", 10*int(newLimit))
	fmt.Println("2. no panic... len:", len(v2))
}
