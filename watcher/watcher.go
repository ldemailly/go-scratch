package main

import (
	"bytes"
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"fortio.org/cli"
	"fortio.org/log"
	"github.com/google/go-cmp/cmp"
)

func main() {
	pf := flag.Duration("t", 10*time.Second, "Polling interval")
	sf := flag.String("s", "", "what to search for in the reply")
	cli.MinArgs = 1
	cli.MaxArgs = 1
	cli.ArgsHelp = "url"
	cli.Main()
	url := flag.Args()[0]
	pollingInterval := *pf
	etag := ""
	lastModified := ""
	var prevChecksum [32]byte
	client := &http.Client{}
	prevBody := ""
	for {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatalf("Error creating request: %v", err)
		}

		if etag != "" {
			log.Infof("Adding If-None-Match: %q", etag)
			req.Header.Add("If-None-Match", etag)
		}
		if lastModified != "" {
			log.Infof("Adding If-Modified-Since %q", etag)
			req.Header.Add("If-Modified-Since", lastModified)
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Error doing request: %v", err)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading body: %v", err)
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

		if *sf != "" {
			if strings.Contains(bodyStr, *sf) {
				log.Warnf("Found %q in the body", *sf)
			}
		}
		switch resp.StatusCode {
		case http.StatusNotModified:
			log.Infof("Header based content has not changed.")
		case http.StatusOK:
			etag = resp.Header.Get("ETag")
			lastModified = resp.Header.Get("Last-Modified")
			if etag != "" || lastModified != "" {
				log.S(log.NoLevel, "Content has changed", log.Any("ETag", etag), log.Any("Last-Modified", lastModified))
			}
			if bytes.Compare(checksum[:], prevChecksum[:]) != 0 {
				log.Infof("Content has changed %x", checksum)
				if prevBody != "" {
					fmt.Println(cmp.Diff(bodyStr, prevBody))
				}
				prevBody = bodyStr
				copy(prevChecksum[:], checksum[:])
			} else {
				log.Infof("Content has not changed.")
			}
		default:
			log.Fatalf("Unexpected status code: %v", resp.StatusCode)
		}
		log.LogVf("Sleeping for %v", pollingInterval)
		time.Sleep(pollingInterval)
	}
}
