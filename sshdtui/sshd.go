package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
)

const KeyFile = "./host_key"

func Handler(s ssh.Session) {
	p, c, ok := s.Pty()
	log.Printf("Pty: %+v, winch chan: %v, ok: %v", p, c, ok)
	width, height := p.Window.Width, p.Window.Height
	fmt.Fprintf(s, "Hello, SSH!\nInitial terminal width: %d, height: %d\nResize me!\n", width, height)
	for {
		select {
		case w := <-c:
			if w.Width == width && w.Height == height {
				continue
			}
			width, height = w.Width, w.Height
			fmt.Fprintf(s, "Window size changed: width=%d, height=%d\n", width, height)
		case <-time.After(10 * time.Second):
			fmt.Fprintf(s, "No window size change for 10 seconds, closing session.\n")
			return
		}
	}
}

func CheckKeyFile(keyFile string) {
	_, err := os.Stat(keyFile)
	if err == nil {
		log.Printf("Using existing key file %s", keyFile)
		return
	}
	if !os.IsNotExist(err) {
		log.Fatal(err)
	}
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}
	privateKeyBlock, err := gossh.MarshalPrivateKey(key, "")
	if err != nil {
		log.Fatal(err)
	}
	privateKeyBytes := pem.EncodeToMemory(privateKeyBlock)
	err = os.WriteFile(keyFile, privateKeyBytes, 0o600)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Generated new host key at %s", keyFile)
}

func main() {
	CheckKeyFile(KeyFile)
	log.Fatal(ssh.ListenAndServe(":2222", Handler, ssh.HostKeyFile(KeyFile)))
}
