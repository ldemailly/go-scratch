//go:build linux
// +build linux

// Demonstrates nil dereference without panic/segv, works on Linux after
// echo 0 > /proc/sys/vm/mmap_min_addr
// ported from https://github.com/ldemailly/experimental/blob/master/misc-c/happynull.c
// no cgo version of ../nilisfine/nilisfine.go
package main

import (
	"fmt"
	"syscall"
)

func init() {
	addr := uintptr(0)
	_, err := syscall.Mmap(-1, 0, 4096, 0x3, 0x32)
	if err != nil {
		fmt.Printf("mmap failed (%v), you need to run:\n\necho 0 > /proc/sys/vm/mmap_min_addr\n\nas is it will panic.", err)
	}
}

func main() {
	var p *int
	fmt.Printf("Before nil deref: p=%p\n", p)
	*p = 42
	fmt.Printf("After nil deref: p=%p, *p=%d\n", p, *p)
}
