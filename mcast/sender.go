package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	a, err := net.ResolveUDPAddr("udp4", "239.0.10.10:3000")
	if err != nil {
		slog.Error(err.Error())
		return
	}
	c, err := net.DialUDP("udp4", nil, a)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	go func() {
		i := 1
		t := time.NewTicker(time.Second)
		for {
			select {
			case <-t.C:
				c.Write([]byte(fmt.Sprintf("hello world #%d", i)))
				i++
			case <-ctx.Done():

				return
			}
		}

	}()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	cancel()
	c.Close()
}
