package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
)

func main() {
	a, err := net.ResolveUDPAddr("udp4", "239.0.10.10:3000")
	if err != nil {
		slog.Error(err.Error())
		return
	}
	d, err := net.ResolveUDPAddr("udp4", "239.0.10.11:3000")
	if err != nil {
		slog.Error(err.Error())
		return
	}
	l, err := net.ListenMulticastUDP("udp4", nil, a)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	c, err := net.DialUDP("udp4", nil, d)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	buf := make([]byte, 1500)
	go func() {
		i := 1
		for {
			n, remote, err := l.ReadFromUDP(buf)
			if err != nil {
				slog.Error(err.Error())
				continue
			}
			msg := fmt.Sprintf("proxy #%d %s:%d %s", i, remote.IP.String(), remote.Port, string(buf[:n]))
			slog.Info(msg)
			c.Write([]byte(msg))
			i++
		}
	}()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	l.Close()
	c.Close()
}
