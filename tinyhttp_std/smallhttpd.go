/*
Normal go version of ../tinyhttp/tinyhttp.go.

17:48:19 main tinyhttp_std/$ go run .
17:48:23.216 r1 [INF] scli.go:125> Starting, command="tinyhttp_std", version="dev  go1.22.5 arm64 darwin", go-max-procs=11
17:48:23.216 r1 [INF] smallhttpd.go:50> Listening on [::]:8000
17:48:23.217 r22 [INF] smallhttpd.go:52> Server entering listen
17:48:40.070 r24 [INF]> http srv, method="GET", url="/", host="localhost:8000", proto="HTTP/1.1", remote_addr="[::1]:61940", user-agent="curl/8.8.0", status=200, size=14, microsec=3
17:48:45.652 r26 [INF]> http srv, method="GET", url="/foo/bar", host="localhost:8000", proto="HTTP/1.1", remote_addr="[::1]:61942", user-agent="curl/8.8.0", status=404, size=19, microsec=54
^C
17:48:54.236 r1 [WRN] scli.go:139> Interrupt received.
17:48:54.237 r22 [INF] smallhttpd.go:57> Server stopped listening
17:48:54.237 r1 [INF] smallhttpd.go:67> All done with graceful shutdown
*/
package main

import (
	"context"
	"errors"
	"flag"
	"net"
	"net/http"
	"time"

	"fortio.org/dflag"
	"fortio.org/dflag/endpoint"
	"fortio.org/log"
	"fortio.org/scli"
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
	l, err := net.Listen("tcp", *port)
	if err != nil {
		log.Fatalf("Unable to listen: %v", err)
	}
	log.Infof("Listening on %v", l.Addr())
	go func() {
		log.Infof("Server entering listen")
		err := server.Serve(l)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
		log.Infof("Server stopped listening")
	}()
	scli.UntilInterrupted()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Shut down the server
	if err := server.Shutdown(ctx); err != nil {
		log.Critf("Server Shutdown: %s", err)
		return
	}
	log.Infof("All done with graceful shutdown")
}
