package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"fortio.org/log"
	"fortio.org/scli"
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
)

const KeyFile = "./host_key"

func Handler(s ssh.Session) {
	p, c, ok := s.Pty()
	log.S(log.Info, "Pty:", log.Any("pty", p), log.Any("ok", ok))
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
		log.Infof("Using existing key file %s", keyFile)
		return
	}
	if !os.IsNotExist(err) {
		log.Fatalf("%s: %v", keyFile, err)
	}
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("%v", err)
	}
	privateKeyBlock, err := gossh.MarshalPrivateKey(key, "")
	if err != nil {
		log.Fatalf("%v", err)
	}
	privateKeyBytes := pem.EncodeToMemory(privateKeyBlock)
	err = os.WriteFile(keyFile, privateKeyBytes, 0o600)
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Warnf("Generated new host key at %s", keyFile)
}

func main() {
	scli.ServerMain()
	CheckKeyFile(KeyFile)
	log.Fatalf("%v", ssh.ListenAndServe(":2222", Handler, ssh.HostKeyFile(KeyFile)))
}
