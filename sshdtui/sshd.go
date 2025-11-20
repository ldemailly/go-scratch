package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"encoding/pem"
	"flag"
	"os"
	"time"

	"fortio.org/log"
	"fortio.org/scli"
	"fortio.org/terminal"
	"fortio.org/terminal/ansipixels"
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
)

const KeyFile = "./host_key"

type InputAdapter struct {
	ansipixels.InputReader
}

func (ia InputAdapter) RawMode() error {
	return nil
}

func (ia InputAdapter) NormalMode() error {
	return nil
}

func (ia InputAdapter) StartDirect() {
}

func Handler(s ssh.Session) {
	log.Infof("New SSH session from %v user=%s", s.RemoteAddr(), s.User())
	p, c, ok := s.Pty()
	log.S(log.Info, "Pty:", log.Any("pty", p), log.Any("ok", ok))
	width, height := p.Window.Width, p.Window.Height
	ap := ansipixels.AnsiPixels{
		Out: bufio.NewWriter(s),
		FPS: 60,
		H:   height,
		W:   width,
		C:   make(chan os.Signal, 1),
	}
	fps := 60
	timeout := time.Duration(1000/fps) * time.Millisecond
	ir := terminal.NewTimeoutReader(s, timeout)
	ia := InputAdapter{ir}
	ap.SharedInput = ia
	ap.GetSize = func() error {
		ap.W, ap.H = width, height
		return nil
	}
	ap.ClearScreen()
	ap.WriteBoxed(ap.H/2-1, "Ansipixels sshd demo!\nInitial terminal width: %d, height: %d\nResize me!\nQ to quit", width, height)
	ap.EndSyncMode()
	ap.OnResize = func() error {
		ap.ClearScreen()
		ap.WriteBoxed(ap.H/2-3, "Window size changed:\n%d x %d ", width, height)
		ap.EndSyncMode()
		return nil
	}
	keepGoing := true
	for keepGoing {
		select {
		case w := <-c:
			if w.Width == width && w.Height == height {
				continue
			}
			width, height = w.Width, w.Height
			log.Infof("Window resized to %dx%d", width, height)
			// Only send if it's not already queued
			select {
			case ap.C <- ansipixels.ResizeSignal:
				// signal sent
			default:
				// channel full; nothing to do (will get processed in next ReadOrResizeOrSignalOnce)
			}
		default:
			n, err := ap.ReadOrResizeOrSignalOnce()
			if err != nil {
				log.Errf("Error reading input or resizing or signaling: %v", err)
				keepGoing = false
			}
			if n == 0 {
				continue
			}
			c := ap.Data[0]
			switch c {
			case 3, 'q': // Ctrl-C or 'q'
				ap.WriteAt(0, ap.H-2, "Exit requested, closing session.")
				keepGoing = false
			default:
				// echo back
				ap.WriteAt(0, ap.H-2, "Received %q", ap.Data)
				ap.ClearEndOfLine()
			}
		}
		ap.EndSyncMode()
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
	port := flag.String("port", ":2222", "Port/address to listen on")
	scli.ServerMain()
	CheckKeyFile(KeyFile)
	log.Infof("Starting SSH server on %s", *port)
	log.Fatalf("%v", ssh.ListenAndServe(*port, Handler, ssh.HostKeyFile(KeyFile)))
}
