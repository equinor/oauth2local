package ipc

import (
	"context"
	"net"

	winio "github.com/Microsoft/go-winio"
)

const pipeName = `\\.\pipe\oauth2local`

func listener() (net.Listener, error) {
	l, err := winio.ListenPipe(pipeName, nil)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func localPipeDial(ctx context.Context, addr string) (c net.Conn, err error) {
	c, err = winio.DialPipe(pipeName, nil)
	return
}
