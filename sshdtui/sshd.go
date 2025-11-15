package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"encoding/pem"
	"os"
	"time"

	"fortio.org/log"
	"fortio.org/scli"
	"fortio.org/terminal/ansipixels"
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
)

const KeyFile = "./host_key"

func Handler(s ssh.Session) {
	p, c, ok := s.Pty()
	log.S(log.Info, "Pty:", log.Any("pty", p), log.Any("ok", ok))
	width, height := p.Window.Width, p.Window.Height
	ap := ansipixels.AnsiPixels{
		Out: bufio.NewWriter(s),
		FPS: 60,
		H:   height,
		W:   width,
	}
	ap.ClearScreen()
	ap.WriteBoxed(ap.H/2-1, "Hello, SSH!\nInitial terminal width: %d, height: %d\nResize me!", width, height)
	ap.EndSyncMode()

	for {
		select {
		// case <-s.Read():
		case w := <-c:
			if w.Width == width && w.Height == height {
				continue
			}
			width, height = w.Width, w.Height
			ap.H = height
			ap.W = width
			ap.ClearScreen()
			ap.WriteBoxed(ap.H/2-3, "Window size changed:\n%d x %d ", width, height)
			ap.EndSyncMode()
		case <-time.After(10 * time.Second):
			ap.WriteAt(0, ap.H-2, "No window size change for 10 seconds, closing session.")
			ap.MoveCursor(0, ap.H-1)
			ap.EndSyncMode()
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
