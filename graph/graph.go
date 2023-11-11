package main

import (
	"flag"
	"net/http"
	"os"
	"strconv"

	"fortio.org/fortio/fhttp"
	"fortio.org/log"
	"fortio.org/scli"
	svg "github.com/ajstarks/svgo"
)

// handler for serving index file
func staticHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func drawSVG(w http.ResponseWriter, r *http.Request) {
	// Read the i query parameter
	iStr := r.URL.Query().Get("i")
	i, err := strconv.Atoi(iStr)
	if err != nil {
		log.Errf("Failed to parse i parameter: %v", err)
	}
	log.Infof("drawSVG: i=%d", i)
	// Set the content type to SVG
	w.Header().Set("Content-Type", "image/svg+xml")
	//w.Header().Set("Access-Control-Allow-Origin", "*") // For prod, replace * with specific allowed domain
	canvas := svg.New(w)
	canvas.Start(200, 200)
	if i%2 == 0 {
		canvas.Circle(50, 50, 20, "fill:green")
	} else {
		canvas.Rect(40, 40, 20, 20, "fill:red")
	}
	canvas.Circle(150, 50, 20, "fill:green")
	canvas.Circle(100, 150, 20, "fill:green")
	canvas.End()
}

func main() {
	port := flag.String("port", ":8080", "port to listen on")
	scli.ServerMain()
	mux, addr := fhttp.HTTPServer("svg", *port)
	if addr == nil {
		os.Exit(1) // already logged
	}
	mux.HandleFunc("/svg", fhttp.LogAndCall("drawSVG", drawSVG))
	mux.HandleFunc("/", fhttp.LogAndCall("static", staticHandler))
	scli.UntilInterrupted()
}
