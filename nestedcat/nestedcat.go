/*
 Demonstrates that you can cmd.Exec and pass a file name
 to a sub executable (in this case itself)
*/

package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
)

func inChild(filename string) {
	// open file and cat it
	log.Printf("In child process: Opening file %q", filename)
	buf, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	os.Stdout.Write(buf)
}

func main() {
	fnameSub := flag.String("f", "", "file to cat in the sub process")
	flag.Parse()
	if *fnameSub != "" {
		inChild(*fnameSub)
		return
	}
	log.Printf("In parent process")
	if len(flag.Args()) == 0 {
		log.Fatal("No file name provided")
	}
	// Pass file to child process
	ourName := os.Args[0]
	fname := flag.Arg(1)
	sub := exec.Cmd{
		Path: ourName,
		Args: []string{ourName, "-f", fname},
	}
	log.Printf("Running %v", sub.Args)
	sub.Stdout = os.Stdout
	sub.Stderr = os.Stderr
	err := sub.Run()
	if err != nil {
		log.Fatal(err)
	}
}
