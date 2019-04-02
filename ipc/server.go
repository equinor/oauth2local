package ipc

import (
	"context"
	"fmt"
	"log"
	"net/url"

	pb "github.com/equinor/oauth2local/ipc/localauth"
	"github.com/equinor/oauth2local/oauth2"
	"google.golang.org/grpc"
)

type Server struct {
	oauthHandler oauth2.Handler
}

func NewServer(oauthHandler oauth2.Handler) (s *Server) {
	return &Server{oauthHandler: oauthHandler}
}

func (s *Server) GetAccessToken(ctx context.Context, _ *pb.Empty) (*pb.ATResponse, error) {

	a, err := s.oauthHandler.GetAccessToken()
	if err != nil {
		log.Println("Error:", err)
		return nil, err
	}
	return &pb.ATResponse{AccessToken: a}, nil
}

func (s *Server) Callback(ctx context.Context, cb *pb.CBRequest) (*pb.Empty, error) {

	rUrl, err := url.Parse(cb.Url)
	if err != nil {
		return nil, err
	}
	err = s.oauthHandler.UpdateFromRedirect(rUrl)
	if err != nil {
		return nil, err
	}
	// s.store.SetToken(storage.AuthorizationCode, c)
	return &pb.Empty{}, nil

}

func (s *Server) Ping(ctx context.Context, _ *pb.Empty) (*pb.PingResponse, error) {
	return &pb.PingResponse{Message: "pong"}, nil
}

func (s *Server) Serve() error {
	lis, err := listener()
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	defer lis.Close()

	gs := grpc.NewServer()
	pb.RegisterLocalAuthServer(gs, s)

	return gs.Serve(lis)
}
