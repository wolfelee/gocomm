package xgrpc

import (
	"context"
	"github.com/wolfelee/gocomm/pkg/jlog"
	"net"
	"strconv"

	"github.com/wolfelee/gocomm/pkg/constant"
	"github.com/wolfelee/gocomm/pkg/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server ...
type Server struct {
	*grpc.Server
	listener net.Listener
	*Config
}

func newServer(config *Config) *Server {
	var streamInterceptors = append(
		[]grpc.StreamServerInterceptor{defaultStreamServerInterceptor()},
		config.streamInterceptors...,
	)

	var unaryInterceptors = append(
		[]grpc.UnaryServerInterceptor{defaultUnaryServerInterceptor()},
		config.unaryInterceptors...,
	)

	config.serverOptions = append(config.serverOptions,
		grpc.StreamInterceptor(StreamInterceptorChain(streamInterceptors...)),
		grpc.UnaryInterceptor(UnaryInterceptorChain(unaryInterceptors...)),
	)

	newServer := grpc.NewServer(config.serverOptions...)
	reflection.Register(newServer)

	listener, err := net.Listen(config.Network, config.Address())
	if err != nil {
		jlog.Error("new grpc server err" + err.Error())
	}
	config.Port = listener.Addr().(*net.TCPAddr).Port

	jlog.Info("grpc server listen on: 0.0.0.0:" + strconv.Itoa(config.Port))

	return &Server{
		Server:   newServer,
		listener: listener,
		Config:   config,
	}
}

// Server implements server.Server interface.
func (s *Server) Serve() error {
	err := s.Server.Serve(s.listener)
	return err
}

// Stop implements server.Server interface
// it will terminate echo server immediately
func (s *Server) Stop() error {
	s.Server.Stop()
	return nil
}

// GracefulStop implements server.Server interface
// it will stop echo server gracefully
func (s *Server) GracefulStop(ctx context.Context) error {
	s.Server.GracefulStop()
	return nil
}

// Info returns server info, used by governor and consumer balancer
func (s *Server) Info() *server.ServiceInfo {
	serviceAddress := s.listener.Addr().String()
	if s.Config.ServiceAddress != "" {
		serviceAddress = s.Config.ServiceAddress
	}

	info := server.ApplyOptions(
		server.WithScheme("grpc"),
		server.WithAddress(serviceAddress),
		server.WithKind(constant.ServiceProvider),
	)
	return &info
}
