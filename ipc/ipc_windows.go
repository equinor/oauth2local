package ipc

import (
	"context"
	
	
	"net"

	winio "github.com/Microsoft/go-winio"
)

// type IPCServer struct {
// 	listener net.Listener
// }

const pipeName = `\\.\pipe\oauth2local`

// func serve(l net.Listener) {
// 	for {
// 		log.Println("Listening on pipe")
// 		c, err := l.Accept()
// 		if err != nil {
// 			log.Fatal(err)

// 		}
// 		// go func(c net.Conn) {
// 		defer c.Close()
// 		rw := bufio.NewReadWriter(bufio.NewReader(c), bufio.NewWriter(c))
// 		s, err := rw.ReadString('\n')
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		log.Println("got " + s)
// 		if s == "<ping>" {
// 			_, err = rw.WriteString("<pong>")
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 		} else {
// 			_, err = rw.WriteString("")
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 		}

// 		err = rw.Flush()
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		// }(c)

// 	}
// }

func listener() (net.Listener, error) {
	l, err := winio.ListenPipe(pipeName, nil)
	if err != nil {
		return nil, err
	}
	// go serve(l)
	return l, nil
}

func localPipeDial(ctx context.Context,addr string) (c net.Conn,err error ){
	c, err = winio.DialPipe(pipeName, nil)
	
		return 
	



	// rw := bufio.NewReadWriter(bufio.NewReader(c), bufio.NewWriter(c))
	// _, err = rw.WriteString("<" + msgType + ">" + msg + "\n")

	// if err != nil {
	// 	return err
	// }
	// err = rw.Flush()
	// if err != nil {
	// 	return err
	// }

	// s, err := rw.ReadString('\n')
	// if err != nil {
	// 	return err
	// }
	// resp <- s
	// return nil
}
