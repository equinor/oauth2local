package ipc

import (
	"context"
	homedir "github.com/mitchellh/go-homedir"
	"net"
)

const pipeName = `/oauth2local.socket`

func listener() (net.Listener, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, err

	}

	l, err := net.Listen("unix", home+pipeName)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func localPipeDial(ctx context.Context, addr string) (c net.Conn, err error) {
	home, err := homedir.Dir()
	if err != nil {
		return

	}
	c, err = net.Dial("unix", home+pipeName)
	return
}
