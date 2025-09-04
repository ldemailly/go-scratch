//go:build linux
// +build linux

// Demonstrates nil dereference without panic/segv, works on Linux after
// echo 0 > /proc/sys/vm/mmap_min_addr
// ported from https://github.com/ldemailly/experimental/blob/master/misc-c/happynull.c
// no cgo version of ../nilisfine/nilisfine.go
package main

import (
	"fmt"

	"github.com/ebitengine/purego"
)

func openLibrary(name string) (uintptr, error) {
	return purego.Dlopen(name, purego.RTLD_NOW|purego.RTLD_GLOBAL)
}

var mmap func(addr uintptr, length uintptr, prot int, flags int, fd int, offset int64) uintptr

func init() {
	libc, err := openLibrary("libc.so.6")
	if err != nil {
		panic(err)
	}
	purego.RegisterLibFunc(&mmap, libc, "mmap")
	addr := uintptr(0)
	ret := mmap(addr, 4096, 0x3, 0x32, -1, 0)
	if ret == ^uintptr(0) {
		fmt.Println("mmap failed, you need to run:\n\necho 0 > /proc/sys/vm/mmap_min_addr\n\nas is it will panic.")
	}
}

func main() {
	var p *int
	fmt.Printf("Before nil deref: p=%p\n", p)
	*p = 42
	fmt.Printf("After nil deref: p=%p, *p=%d\n", p, *p)
}
