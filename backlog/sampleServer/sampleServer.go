package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/ldemailly/go-scratch/backlog"
)

// Show a server using a set backlog

var connCount int

func main() {
	b := flag.Int("b", 1, "`backlog` to set")
	p := flag.String("p", ":8118", "`port` to listen on")
	flag.Parse()
	// listen on port *p with backlog *b
	listener, err := net.Listen("tcp", *p)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Server listening on %s", listener.Addr().String())

	// Set the backlog
	if err := backlog.Set(listener, *b); err != nil {
		log.Fatalf("Failed to set backlog: %v", err)
	}
	log.Printf("Backlog set to %d - not accepting any...", *b)
	// don't actually accept connections so we check the backlog
	select {}
	/*
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("Failed to accept connection: %v", err)
				continue
			}
			connCount++
			go handleConnection(conn, connCount)
		}
	*/
}

func handleConnection(conn net.Conn, id int) {
	defer conn.Close()
	log.Printf("New connection #%d from %s", id, conn.RemoteAddr())
	conn.Write([]byte(fmt.Sprintf("Hello, connection #%d\n", id)))
	n, err := io.Copy(os.Stderr, conn)
	if err != nil {
		log.Printf("Failed to copy data: %v", err)
	}
	log.Printf("Copied %d bytes", n)
	return
}
