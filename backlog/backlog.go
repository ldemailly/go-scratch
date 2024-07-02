package backlog

import (
	"fmt"
	"net"
	"syscall"
)

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
