package ipc

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/zaker/oauth2local/ipc/localauth"
	"google.golang.org/grpc"
)

type server struct{}
type client struct{}

func (s *server) GetAccessToken(ctx context.Context, _ *pb.Empty) (*pb.ATResponse, error) {

	return nil, nil
}

func (s *server) UpdateCode(ctx context.Context, cr *pb.UCRequest) (*pb.Empty, error) {
	r := new(pb.Empty)
	fmt.Println("Received:", cr.Code)
	return r, nil

}

func (s *server) Ping(ctx context.Context, _ *pb.Empty) (*pb.PingResponse, error) {
	r := new(pb.PingResponse)
	r.Message = "pong"
	return r, nil
}

// func localPipe(ctx context.Context,addr string)(c net.Conn, err error){

// }

func HasSovereign() (bool, error) {
	conn, err := grpc.Dial("pipe", grpc.WithInsecure(), grpc.WithContextDialer(localPipeDial))
	if err != nil {
		return false, err
	}
	defer conn.Close()
	c := pb.NewLocalAuthClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Ping(ctx, new(pb.Empty))
	if err != nil {
		return false, err
	}
	return r.Message == "pong", nil
}

func SendCode(code string) error {
	conn, err := grpc.Dial("pipe", grpc.WithInsecure(), grpc.WithContextDialer(localPipeDial))
	if err != nil {
		return err
	}
	defer conn.Close()
	c := pb.NewLocalAuthClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	_, err = c.UpdateCode(ctx, &pb.UCRequest{Code: code})
	if err != nil {
		return err
	}
	return nil
}

func StartServer() {
	lis, err := listener()
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterLocalAuthServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
