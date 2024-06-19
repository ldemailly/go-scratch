package main

import (
    "fmt"
    "time"
)

var i int32 = 0

func IncPrint(id int) {
    for {
        v := i
        time.Sleep(time.Millisecond * 1) // trigger i != v+1
        i++
        fmt.Printf("T#%d, i was %d and now %d\n", id, v, i)
        time.Sleep(time.Millisecond * 1)
        if i > 10 {
            return
        }
    }
}

func main() {
    go IncPrint(1)
    IncPrint(2)
}
