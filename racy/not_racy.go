package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var a = "before"

var s int32

var wg sync.WaitGroup

func SendingGoroutine() {
	var r int32
	for r != 1 {
		r = atomic.LoadInt32(&s)
	}
	// ReceivingGoroutine started and done reading a
	a = "hello, world"
	atomic.StoreInt32(&s, 2)
	wg.Done()
}

func ReceivingGoroutine() {
	var r int32
	r = atomic.LoadInt32(&s)
	if r != 0 {
		panic("unexpected state")
	}
	// Safe to read a
	fmt.Println(a)
	// let SendingGoroutine continue:
	atomic.StoreInt32(&s, 1)
	// Wait that it's done writing
	for r != 2 {
		r = atomic.LoadInt32(&s)
	}
	// a has been written
	fmt.Println(a)
	wg.Done()
}

func main() {
	wg.Add(2)
	go SendingGoroutine()
	go ReceivingGoroutine()
	wg.Wait()
}
