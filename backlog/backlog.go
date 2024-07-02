package backlog

import (
	"fmt"
	"net"
	"syscall"
)

// backlog.Set() will attempt to reset the TCP connection backlog queue to the given value.
// Works on MacOS and Linux. On linux seems to allow n+1 connections to be queued (one in
// userland 1 in kernel maybe? despite no accept).
// To test
//   - run the server
//     go run sampleServer.go -b 3
//   - run the client with -n 5 only first 3 (or 4) will connect
//     go run sampleClient.go -n 5
//
// On linux `ss -ltn6` will show the backlog as SendQ column.
//
// PS: none of this meant to work according to POSIX, it just happens to seem to do so.\
// pending https://github.com/golang/go/issues/39000
func Set(l net.Listener, backlog int) error {
	tl, ok := l.(*net.TCPListener)
	if !ok {
		return fmt.Errorf("only tcp listener supported, called with %#v", l)
	}
	file, err := tl.File()
	if err != nil {
		return err
	}
	fd := int(file.Fd())
	err = syscall.Listen(fd, backlog)
	if err != nil {
		return err
	}
	return nil
}
