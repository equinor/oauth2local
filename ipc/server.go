package ipc

import (
	"context"
	"fmt"
	"net/url"

	pb "github.com/equinor/oauth2local/ipc/localauth"
	"github.com/equinor/oauth2local/oauth2"
	jww "github.com/spf13/jwalterweatherman"
	"google.golang.org/grpc"
)

type Server struct {
	oauthHandler oauth2.Handler
}

func NewServer(oauthHandler oauth2.Handler) (s *Server) {
	return &Server{oauthHandler: oauthHandler}
}

func (s *Server) GetAccessToken(ctx context.Context, _ *pb.Empty) (*pb.ATResponse, error) {
	jww.DEBUG.Println("Get access token")
	a, err := s.oauthHandler.GetAccessToken()
	if err != nil {
		jww.ERROR.Println("Error:", err)
		return nil, err
	}
	return &pb.ATResponse{AccessToken: a}, nil
}

func (s *Server) Callback(ctx context.Context, cb *pb.CBRequest) (*pb.Empty, error) {
	jww.DEBUG.Println("Callback from provider: ", cb.Url)
	rURL, err := url.Parse(cb.Url)
	if err != nil {
		jww.ERROR.Println("Callback error:", err)
		return nil, err
	}
	err = s.oauthHandler.UpdateFromRedirect(oauth2.DecodeRedirect(rURL))
	if err != nil {
		jww.ERROR.Println("Callback error:", err)
		return nil, err
	}
	jww.INFO.Println("Auth ok")
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
