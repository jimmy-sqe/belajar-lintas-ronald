package grpc

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	addr string
	srv  *grpc.Server
}

func NewServer(host, port string) *Server {
	return &Server{addr: net.JoinHostPort(host, port), srv: grpc.NewServer()}
}

func (s *Server) Registrar() grpc.ServiceRegistrar { return s.srv }

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("grpc: listen %s: %w", s.addr, err)
	}
	return s.srv.Serve(lis)
}

func (s *Server) Stop(_ context.Context) { s.srv.GracefulStop() }
