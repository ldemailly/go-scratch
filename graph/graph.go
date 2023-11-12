package main

import (
	"embed"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"math"
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
	w       io.Writer
	started bool
}

func NewSVG(out io.Writer) *SVG {
	return &SVG{w: out}
}

// For fully compliant svg, needed for serving svg images directly (vs as inner html)
func (s *SVG) StartLong(x, y int) {
	fmt.Fprintf(s.w, "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 %d %d\">", x, y)
	s.started = true
}

func (s *SVG) Start() {
	fmt.Fprintf(s.w, "<svg>")
	s.started = true
}

func (s *SVG) Misc(svg string) {
	if !s.started {
		s.Start()
	}
	fmt.Fprint(s.w, svg)
}

func (s *SVG) Circle(x, y, r int, style string) {
	if !s.started {
		s.Start()
	}
	fmt.Fprintf(s.w, `<circle cx="%d" cy="%d" r="%d" style="%s"/>`, x, y, r, style)
}

func (s *SVG) Rect(x, y, w, h int, style string) {
	if !s.started {
		s.Start()
	}
	fmt.Fprintf(s.w, `<rect x="%d" y="%d" width="%d" height="%d" style="%s"/>`, x, y, w, h, style)
}

func (s *SVG) Line(x1, y1, x2, y2 int, style string) {
	if !s.started {
		s.Start()
	}
	fmt.Fprintf(s.w, `<line x1="%d" y1="%d" x2="%d" y2="%d" style="%s"/>`, x1, y1, x2, y2, style)
}

func (s *SVG) End() {
	fmt.Fprintln(s.w, "</svg>") // extra newline to separate multiple svg
	s.started = false
}

func (s *SVG) Text(x, y int, text string, style string) {
	if !s.started {
		s.Start()
	}
	fmt.Fprintf(s.w, `<text x="%d" y="%d" style="%s">%s</text>`, x, y, style, text)
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	svg := NewSVG(w)
	svg.StartLong(42, 42)
	svg.Misc(`<defs><linearGradient id="orangeGradient" x1="0%" y1="0%" x2="100%" y2="100%"><stop offset="0%" style="stop-color:#FFA500; stop-opacity:1" /><stop offset="100%" style="stop-color:#FF4500; stop-opacity:1" /></linearGradient></defs>`)
	svg.Circle(21, 21, 20, "fill:url(#orangeGradient)")
	svg.Text(21, 21, "SVG", "font-size:18px;fill:white;text-anchor:middle;dominant-baseline:central;")
	svg.End()
}

type Point struct {
	x int
	y int
}

func CompleteCircle(svg *SVG, points []Point) {
	mn := len(points)
	if mn%2 != 0 {
		mn--
	}
	for j := 0; j < mn; j++ {
		current := points[j]
		nextIndex := (j + 1) % mn
		next := points[nextIndex]
		log.Debugf("Completing circle from (%d, %d) to (%d, %d)", current.x, current.y, next.x, next.y)
		svg.Line(current.x, current.y, next.x, next.y, "stroke:rgb(200,0,0);stroke-width:2")
	}
}

func DrawOpposites(svg *SVG, points []Point) {
	mn := len(points)
	if mn%2 != 0 {
		mn--
	}
	for j := 0; j < mn/2; j++ {
		current := points[j]
		nextIndex := (j + mn/2) % mn
		next := points[nextIndex]
		log.Debugf("Drawing opposite to center: line from (%d, %d) to (%d, %d)", current.x, current.y, next.x, next.y)
		svg.Line(current.x, current.y, next.x, next.y, "stroke:rgb(100,0,200);stroke-width:2")
	}
}

func drawSVG(w http.ResponseWriter, r *http.Request) {
	// Read the i query parameter
	iStr := r.URL.Query().Get("i")
	i, err := strconv.Atoi(iStr)
	if err != nil {
		log.Errf("Failed to parse i parameter: %v", err)
	}
	nStr := r.URL.Query().Get("n")
	n, err := strconv.Atoi(nStr)
	if err != nil {
		log.Errf("Failed to parse n parameter: %v", err)
	}
	if n < 1 || n > 100 {
		log.Warnf("Fixed invalid n parameter: %d", n)
		n = 1
	}
	log.S(log.Info, "drawSVG", log.Int("i", i), log.Int("n", n))
	// Set the content type to SVG
	w.Header().Set("Content-Type", "image/svg+xml")
	//w.Header().Set("Access-Control-Allow-Origin", "*") // For prod, replace * with specific allowed domain
	svg := NewSVG(w)
	svg.StartLong(512, 512)
	centerX := 256
	centerY := 256
	radius := 256.0 - 16.0
	// point list
	points := make([]Point, n)
	odd := (n%2 != 0)
	mn := n
	if odd {
		points[n-1] = Point{centerX, centerY}
		svg.Circle(centerX, centerY, 8, "fill:blue")
		mn--
	}
	for j := 0; j < mn; j++ {
		angle := 2 * math.Pi * float64(j) / float64(mn)
		x := centerX + int(radius*math.Cos(angle))
		y := centerY + int(radius*math.Sin(angle))
		points[j] = Point{x, y}
		log.Infof("x %d y %d", x, y)
		svg.Circle(x, y, 8, "fill:green")
	}
	// First connection is either to the next point or to the opposite side of the circle
	// then the other for 2nd connection
	if odd {
		DrawOpposites(svg, points)
		if i > 1 {
			CompleteCircle(svg, points)
		}
	} else {
		CompleteCircle(svg, points)
		if i > 1 {
			DrawOpposites(svg, points)
		}
	}
	// For n > 2 : Connect each point to opposite + 1
	for k := 3; k <= i; k++ {
		for j := 0; j <= mn/2; j++ {
			current := points[j]
			delta := (k - 1) / 2
			if (j+k)%2 == 0 {
				delta = -delta
			}
			nextIndex := (j + mn/2 + delta) % mn
			next := points[nextIndex]
			color := fmt.Sprintf("%d,%d,%d", 150+delta*20, k*255/i, 150-delta*20)
			log.Debugf("Drawing opposite + %d to center: line from (%d, %d) to (%d, %d)", delta, current.x, current.y, next.x, next.y)
			svg.Line(current.x, current.y, next.x, next.y, "stroke:rgb("+color+")")
		}
	}
	/*
		svg.Circle(50, 50, 20, "fill:green")
		if i%2 == 1 {
			l := int(math.RoundToEven(20 * math.Sqrt2)) // radius * sqrt(2) for inscribed square
			svg.Rect(50-l/2, 50-l/2, l, l, "fill:red")
		}
		svg.Circle(150, 50, 20, "fill:green")
		svg.Circle(100, 150, 20, "fill:green")
	*/
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
	mux.HandleFunc("/favicon.svg", fhttp.LogAndCall("favinco", faviconHandler))
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
