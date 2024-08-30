package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/term"
)

func main() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	fmt.Println("Try ^C, ^D, ^Z... Press 'q' to exit\r")
	buf := make([]byte, 1)
	for {
		_, err := os.Stdin.Read(buf)
		if err != nil {
			log.Print(err)
			break
		}
		b := buf[0]
		fmt.Printf("pressed: %d %q\r\n", b, string(b)) // need \r to in raw mode CR+LF.
		if b == 'q' {
			break
		}
	}
}
