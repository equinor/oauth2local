package ipc

import (
	"context"
	"fmt"
	"log"

	pb "github.com/equinor/oauth2local/ipc/localauth"
	"github.com/equinor/oauth2local/oauth2"
	"github.com/equinor/oauth2local/storage"
	"google.golang.org/grpc"
)

type Server struct {
	oauthCli oauth2.Client
	store    storage.Storage
}

func NewServer(cli oauth2.Client, store storage.Storage) (s *Server) {
	return &Server{oauthCli: cli, store: store}
}

func (s *Server) GetAccessToken(ctx context.Context, _ *pb.Empty) (*pb.ATResponse, error) {

	c, err := s.store.GetToken(storage.AuthorizationCode)
	if err != nil {
		log.Println("Error:", err)
		return nil, err
	}
	a, err := s.oauthCli.GetToken(c)
	if err != nil {
		log.Println("Error:", err)
		return nil, err
	}
	return &pb.ATResponse{AccessToken: a}, nil
}

func (s *Server) Callback(ctx context.Context, cb *pb.CBRequest) (*pb.Empty, error) {

	c, err := s.oauthCli.CodeFromURL(cb.Url)
	if err != nil {
		return nil, err
	}
	s.store.SetToken(storage.AuthorizationCode, c)
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
	gs := grpc.NewServer()
	pb.RegisterLocalAuthServer(gs, s)

	return gs.Serve(lis)
}
