package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/term"
)

func main() {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		panic(err)
	}
	defer term.Restore(fd, oldState)
	fmt.Println("Try ^C, ^D, ^Z... Press 'q' to exit\r")
	buf := []byte{0} // 1 byte slice (could be  make([]byte, 1) instead)
	for {
		_, err := os.Stdin.Read(buf)
		if err != nil {
			log.Print(err)
			break
		}
		b := buf[0]
		fmt.Printf("pressed: %d %q\r\n", b, buf) // need \r too in raw mode: CR+LF.
		if b == 'q' {
			break
		}
	}
}
