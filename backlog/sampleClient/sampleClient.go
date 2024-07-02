package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
)

// Show a client doing n connections to a server

func main() {
	n := flag.Int("n", 3, "number of simultaneous connections to make")
	d := flag.String("d", "localhost:8118", "`destination` to connect to")
	flag.Parse()
	// connect to destination *d n times in parallel
	for i := 0; i < *n; i++ {
		go connect(*d)
	}
	select {}
}

func connect(dest string) {
	conn, err := net.Dial("tcp", dest)
	if err != nil {
		log.Printf("Failed to connect to %s: %v", dest, err)
		return
	}
	defer conn.Close()
	log.Printf("Connected to %s", dest)
	n, err := io.Copy(os.Stderr, conn)
	if err != nil {
		log.Printf("Failed to copy data: %v", err)
	}
	log.Printf("Copied %d bytes", n)
	return
}
