package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"fortio.org/cli"
	"fortio.org/log"
	"github.com/google/go-cmp/cmp"
)

// openBrowser opens the specified URL in the default browser of the user.
func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = append(args, "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = "open"
		args = append(args, url)
	case "linux":
		cmd = "xdg-open"
		args = append(args, url)
	default:
		return errors.New("unsupported platform")
	}
	return exec.Command(cmd, args...).Start()
}

type State struct {
	url              string
	etag             string
	lastModified     string
	prevChecksum     [32]byte
	prevBody         string
	client           *http.Client
	disableEtags     bool
	disableKeepAlive bool
	search           string
	doOpen           bool
}

func (s *State) checkOne() error {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, s.url, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	if s.etag != "" && !s.disableEtags {
		log.Infof("Adding If-None-Match: %q", s.etag)
		req.Header.Add("If-None-Match", s.etag)
	}
	if s.lastModified != "" {
		log.Infof("Adding If-Modified-Since %q", s.etag)
		req.Header.Add("If-Modified-Since", s.lastModified)
	}
	if s.disableKeepAlive {
		req.Close = true
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("calling Do: %w", err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("calling ReadAll: %w", err)
	}
	resp.Body.Close()
	log.Infof("Got %d code and %d bytes", resp.StatusCode, len(body))
	log.LogVf("Headers: %+v", resp.Header)
	// checksum body
	checksum := sha256.Sum256(body)
	if len(body) > 0 {
		log.Infof("Checksum: %x", checksum)
	}
	bodyStr := string(body)
	log.Debugf("Body: %v", bodyStr)

	if s.search != "" {
		if strings.Contains(bodyStr, s.search) {
			log.Warnf("Found %q in the body", s.search)
		}
	}
	if resp.StatusCode == http.StatusNotModified {
		log.Infof("Header based content has not changed.")
		// We don't (re)set the body or checksum in this case
		return nil
	}
	if resp.StatusCode == http.StatusNotFound {
		log.Warnf("Got 404 Not found for %q: %q", s.url, bodyStr)
		// bodyStr could be empty and we want to trigger open when status becomes 200
		s.prevBody = "-404 not found-"
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}
	s.etag = resp.Header.Get("ETag")
	s.lastModified = resp.Header.Get("Last-Modified")
	if s.etag != "" || s.lastModified != "" {
		log.S(log.NoLevel, "Got headers (and no 304 so likely change)", log.Any("ETag", s.etag), log.Any("Last-Modified", s.lastModified))
	}
	if bytes.Equal(checksum[:], s.prevChecksum[:]) {
		log.Infof("Content has not changed.")
		return nil
	}
	log.Infof("Content has changed %x", checksum)
	copy(s.prevChecksum[:], checksum[:])
	prevBody := s.prevBody
	s.prevBody = bodyStr
	if prevBody == "" {
		log.Debugf("First time body, not comparing.")
		return nil
	}
	fmt.Println(cmp.Diff(bodyStr, prevBody))
	if !s.doOpen {
		return nil
	}
	err = openBrowser(s.url)
	if err != nil {
		return fmt.Errorf("opening browser: %w", err)
	}
	log.Infof("Opened browser for %q", s.url)
	return nil
}

func main() {
	pf := flag.Duration("t", 10*time.Second, "Polling interval")
	sf := flag.String("s", "", "what to search for in the reply")
	disableKeepAlive := flag.Bool("k", false, "Disable KeepAlive")
	disableEtags := flag.Bool("e", false, "Disable ETags")
	noOpen := flag.Bool("noopen", false, "Do not open the browser upon changes")

	cli.MinArgs = 1
	cli.MaxArgs = 1
	cli.ArgsHelp = "url"
	cli.Main()
	state := &State{
		url:              flag.Args()[0],
		disableEtags:     *disableEtags,
		disableKeepAlive: *disableKeepAlive,
		search:           *sf,
		doOpen:           !*noOpen,
		client:           http.DefaultClient,
	}
	pollingInterval := *pf
	for {
		err := state.checkOne()
		if err != nil {
			log.Fatalf("Error checking: %v", err)
		}
		log.LogVf("Sleeping for %v", pollingInterval)
		time.Sleep(pollingInterval)
	}
}
