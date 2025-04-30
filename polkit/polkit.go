/*
	Demonstrates (self) requesting for root using polkit pkexec: eg:

$ go run github.com/ldemailly/go-scratch/polkit@2f6927ab6d45c8ec2d49572e2d86b1db3491a38e a b c d
go: downloading github.com/ldemailly/go-scratch v0.1.3-0.20250430214433-2f6927ab6d45
go: downloading github.com/ldemailly/go-scratch/polkit v0.0.0-20250430214433-2f6927ab6d45
2025/04/30 21:58:07 Running [/usr/bin/pkexec /tmp/go-build491870754/b001/exe/polkit a b c d]
Running as root 🎉
Hello root world [/tmp/go-build491870754/b001/exe/polkit a b c d]
*/
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

const pkexec = "/usr/bin/pkexec"

func doStuff(args []string) int {
	// do stuff as root
	fmt.Println("Hello root world", args)
	return 0
}

func main() {
	if os.Geteuid() == 0 {
		fmt.Println("Running as root 🎉")
		os.Exit(doStuff(os.Args))
	}
	// get our name and args
	args := []string{pkexec}
	args = append(args, os.Args...)
	sub := exec.Cmd{
		Path: pkexec,
		Args: args,
	}
	log.Printf("Running %v", sub.Args)
	sub.Stdout = os.Stdout
	sub.Stderr = os.Stderr
	err := sub.Run()
	if err != nil {
		log.Fatal(err)
	}
}
