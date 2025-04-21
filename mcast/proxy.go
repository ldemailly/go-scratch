package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
)

func getIF() (*net.Interface, *net.IP) {
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	for _, iface := range ifaces {
		// Skip down or loopback interfaces
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}

			ip = ip.To4()
			if ip == nil {
				continue // Skip non-IPv4
			}

			fmt.Printf("Using interface: %s (%s)\n", iface.Name, ip)
			return &iface, &ip
		}
	}
	return nil, nil
}

func main() {
	iface, localIP := getIF()
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
	l, err := net.ListenMulticastUDP("udp4", iface, a)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	localAddr := &net.UDPAddr{
		IP:   *localIP,
		Port: 0, // let OS choose a free port
	}
	c, err := net.DialUDP("udp4", localAddr, d)
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
