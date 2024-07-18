/* Run this with:

go install github.com/stealthrocket/wasi-go/cmd/wasirun@latest
GOOS=wasip1 GOARCH=wasm go run -exec wasirun .

23:26:19.565 r1 [INF] scli.go:125> Starting, command="tinyhttp", version="dev  go1.22.5 wasm wasip1", go-max-procs=1
23:26:19.566 r1 [INF] tinyhttp.go:60> Listening on 0.0.0.0:8000
23:26:19.566 r1 [INF] tinyhttp.go:61> Server entering listen
23:26:22.659 r6 [INF]> http srv, method="GET", url="/", host="localhost:8000", proto="HTTP/1.1", remote_addr="127.0.0.1:59290", user-agent="curl/8.8.0", status=200, size=14, microsec=35

*/

package main

import (
	"errors"
	"flag"
	"net/http"
	"time"

	// To test without wasi, uncomment and cooment out the last 2 imports
	// wasip1 "net"
	"fortio.org/dflag"
	"fortio.org/dflag/endpoint"
	"fortio.org/log"
	"fortio.org/scli"

	// these:
	// _ "github.com/stealthrocket/net/http" // not needed, that's for http client
	"github.com/stealthrocket/net/wasip1" // needed to use that Listen instead of net.Listen
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	_, err := w.Write([]byte("Hello, TinyGo!"))
	if err != nil {
		log.Fatalf("Error writing response: %v", err)
	}
}

func main() {
	timeout := dflag.New(3*time.Second, "timeout for the http server")
	dflag.Flag("timeout", timeout)
	port := flag.String("port", ":8000", "port to listen on")
	scli.ServerMain()
	mux := http.NewServeMux()
	server := &http.Server{
		ReadHeaderTimeout: timeout.Get(),
		IdleTimeout:       timeout.Get(),
		Handler:           mux,
		ErrorLog:          log.NewStdLogger("http srv", log.Error),
	}
	mux.HandleFunc("/", log.LogAndCall("http srv", helloHandler))
	dflagEndpoint := endpoint.NewFlagsEndpoint(flag.CommandLine, "/flags/set")
	mux.HandleFunc("/flags", log.LogAndCall("dflags-get", dflagEndpoint.ListFlags))
	mux.HandleFunc("/flags/set", log.LogAndCall("dflags-set", dflagEndpoint.SetFlag))
	// Serve (accept) must be on main thread somehow - waspi restriction. BUG?
	l, err := wasip1.Listen("tcp4", *port)
	if err != nil {
		log.Fatalf("Unable to listen: %v", err)
	}
	log.Infof("Listening on %v", l.Addr())
	log.Infof("Server entering listen")
	err = server.Serve(l)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
	log.Infof("Server stopped listening") // not likely given we don't have Shutdown() anymore
}
