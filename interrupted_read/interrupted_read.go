package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fortio.org/log"
	"fortio.org/scli"
)

func main() {
	c, err := net.Dial("tcp", "0.0.0.0:8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer c.Close()
	fmt.Println("Connection successful")
	scli.UntilInterrupted()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Dial error:", err)
		return
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		select {
		case <-ctx.Done():
			log.Printf("Read operation cancelled: %v", ctx.Err())
			return
		case sig := <-sigCh:
			log.Printf("Received signal: %v", sig)
			return
		default:
			conn.SetReadDeadline(time.Now().Add(3 * time.Second)) // Set read deadline
			log.Printf("Entering read...")
			n, err := conn.Read(buf)
			log.Printf("Exiting read... %d %v", n, err)
			if err != nil {
				if ne, ok := err.(net.Error); ok && ne.Timeout() {
					log.Printf("Read timeout, retrying...")
					continue
				}
				log.Printf("Read error:", err)
				return
			}
			log.Printf("Read %d bytes: %s", n, string(buf[:n]))
		}
	}
}
