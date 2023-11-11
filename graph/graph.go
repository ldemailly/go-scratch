package main

import (
	"embed"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strconv"

	"fortio.org/fortio/fhttp"
	"fortio.org/log"
	"fortio.org/scli"
)

// Embed the index.html file
var (
	//go:embed static/*
	staticFS embed.FS
)

type SVG struct {
	w io.Writer
}

func NewSVG(out io.Writer) *SVG {
	svg := &SVG{w: out}
	fmt.Fprint(out, "<svg>")
	return svg
}

func (s *SVG) Circle(x, y, r int, style string) {
	fmt.Fprintf(s.w, `<circle cx="%d" cy="%d" r="%d" style="%s"/>`, x, y, r, style)
}

func (s *SVG) Rect(x, y, w, h int, style string) {
	fmt.Fprintf(s.w, `<rect x="%d" y="%d" width="%d" height="%d" style="%s"/>`, x, y, w, h, style)
}

func (s *SVG) End() {
	fmt.Fprintln(s.w, "</svg>") // extra newline to separate multiple svg
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
	svg := NewSVG(w)
	//canvas.Start(200, 200)
	if i%2 == 0 {
		svg.Circle(50, 50, 20, "fill:green")
	} else {
		svg.Rect(40, 40, 20, 20, "fill:red")
	}
	svg.Circle(150, 50, 20, "fill:green")
	svg.Circle(100, 150, 20, "fill:green")
	svg.End()
}

func main() {
	port := flag.String("port", ":8080", "port to listen on")
	dev := flag.Bool("dev", false, "dev mode (no embed)")
	scli.ServerMain()
	mux, addr := fhttp.HTTPServer("svg", *port)
	if addr == nil {
		os.Exit(1) // already logged
	}
	mux.HandleFunc("/svg", fhttp.LogAndCall("drawSVG", drawSVG))
	var staticHandler http.Handler
	if *dev {
		// So edits to static files are picked up without restart.
		// go run . -dev # must be in this directory, parent of static/
		staticHandler = http.FileServer(http.Dir("./static/"))
	} else {
		subFS, err := fs.Sub(staticFS, "static")
		if err != nil {
			log.Fatalf("Unable to get sub FS: %v", err)
		}
		staticHandler = http.FileServer(http.FS(subFS))
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.LogRequest(r, "static")
		staticHandler.ServeHTTP(w, r)
	})
	scli.UntilInterrupted()
}
