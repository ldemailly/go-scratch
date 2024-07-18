package main

import (
	"context"
	"errors"
	"flag"
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
		Addr:              *port,
		ReadHeaderTimeout: timeout.Get(),
		IdleTimeout:       timeout.Get(),
		Handler:           mux,
		ErrorLog:          log.NewStdLogger("http srv", log.Error),
	}
	mux.HandleFunc("/", log.LogAndCall("http srv", helloHandler))
	dflagEndpoint := endpoint.NewFlagsEndpoint(flag.CommandLine, "/flags/set")
	mux.HandleFunc("/flags", log.LogAndCall("dflags-get", dflagEndpoint.ListFlags))
	mux.HandleFunc("/flags/set", log.LogAndCall("dflags-set", dflagEndpoint.SetFlag))
	go func() {
		log.Infof("Serving on %v", server.Addr)
		err := server.ListenAndServe()
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
