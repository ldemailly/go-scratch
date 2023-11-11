package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"fortio.org/fortio/fhttp"
	"fortio.org/log"
	"fortio.org/scli"
	svg "github.com/ajstarks/svgo"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins
	},
}

// handler for serving files from the static directory
func staticHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func drawSVG(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	//w.Header().Set("Access-Control-Allow-Origin", "*") // For prod, replace * with specific allowed domain
	canvas := svg.New(w)
	canvas.Start(200, 200)
	canvas.Circle(50, 50, 20, "fill:green")
	canvas.Circle(150, 50, 20, "fill:green")
	canvas.Circle(100, 150, 20, "fill:green")
	canvas.End()
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade failed:", err)
		return
	}

	// Defer the closure of the connection
	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Println("Failed to close connection:", err)
		}
	}()
	log.Infof("wsHandler: connection from %v - sleeping", conn.RemoteAddr())
	time.Sleep(10 * time.Second)
	// Create a buffer to write SVG
	var b bytes.Buffer
	canvas := svg.New(&b)

	// Draw the red square
	canvas.Start(200, 200)
	canvas.Rect(40, 40, 20, 20, "fill:red")
	canvas.End()

	// Send the SVG data as a text WebSocket message
	err = conn.WriteMessage(websocket.TextMessage, b.Bytes())
	if err != nil {
		log.Errf("Failed to send message: %v", err)
	}
}

func main() {
	port := flag.String("port", ":8080", "port to listen on")
	scli.ServerMain()
	mux, addr := fhttp.HTTPServer("svg", *port)
	if addr == nil {
		os.Exit(1) // already logged
	}
	mux.HandleFunc("/svg", fhttp.LogAndCall("drawSVG", drawSVG))
	mux.HandleFunc("/ws", fhttp.LogAndCall("ws", wsHandler))
	mux.HandleFunc("/", fhttp.LogAndCall("static", staticHandler))
	scli.UntilInterrupted()
}
