package ipc

import (
	"context"
	"net"
)

const pipeName = `~/oauth2local.socket`

func listener() (net.Listener, error) {
	l, err := net.Listen("unix", pipeName)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func localPipeDial(ctx context.Context, addr string) (c net.Conn, err error) {
	c, err = net.Dial("unix", pipeName)
	return
}
