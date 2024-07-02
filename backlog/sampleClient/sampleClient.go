package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
	"time"
)

// Show a client doing n connections to a server

func main() {
	n := flag.Int("n", 3, "number of simultaneous connections to make")
	d := flag.String("d", "localhost:8118", "`destination` to connect to")
	flag.Parse()
	// connect to destination *d n times in parallel
	for i := 0; i < *n; i++ {
		go connect(*d, i+1)
	}
	select {}
}

func connect(dest string, id int) {
	log.Printf("[%d] Attempting connection to %s", id, dest)
	conn, err := net.DialTimeout("tcp", dest, 10*time.Second)
	if err != nil {
		log.Printf("[%d] Failed to connect to %s: %v", id, dest, err)
		return
	}
	defer conn.Close()
	log.Printf("[%d] Connected to %s", id, dest)
	n, err := io.Copy(os.Stderr, conn)
	if err != nil {
		log.Printf("[%d] Failed to copy data: %v", id, err)
	}
	log.Printf("[%d] Copied %d bytes", id, n)
}
