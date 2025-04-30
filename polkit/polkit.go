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
