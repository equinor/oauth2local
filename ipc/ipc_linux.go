package ipc

import (
	"context"
	homedir "github.com/mitchellh/go-homedir"
	"net"
	"syscall"

	homedir "github.com/mitchellh/go-homedir"
)

const pipeName = `/oauth2local.socket`

func listener() (net.Listener, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, err

	}

	}
	sock := home + pipeName
	syscall.Unlink(sock)
	l, err := net.Listen("unix", sock)
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
