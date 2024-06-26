package xgin

import (
	"context"
	"fmt"
	"github.com/wolfelee/gocomm/pkg/constant"
	"github.com/wolfelee/gocomm/pkg/jlog"
	"github.com/wolfelee/gocomm/pkg/server"
	"net/http"

	"net"

	//"github.com/wolfelee/gocomm/pkg/constant"
	//"github.com/wolfelee/gocomm/pkg/server"
	"github.com/gin-gonic/gin"
)

type Context = gin.Context

// Server ...
type Server struct {
	*gin.Engine
	Server   *http.Server
	config   *Config
	listener net.Listener
}

func newServer(config *Config) *Server {
	listener, err := net.Listen("tcp", config.Address())
	if err != nil {
		jlog.Panic("new gin server err" + err.Error())
	}
	config.Port = listener.Addr().(*net.TCPAddr).Port
	gin.SetMode(config.Mode)

	addr := fmt.Sprintf("0.0.0.0:%d  mode:%s", config.Port, config.Mode)
	jlog.Info("http server listen on:" + addr)

	return &Server{
		Engine:   gin.New(),
		config:   config,
		listener: listener,
	}
}

// Upgrade protocol to WebSocket
func (s *Server) Upgrade(ws *WebSocket) gin.IRoutes {
	return s.GET(ws.Pattern, func(c *gin.Context) {
		ws.Upgrade(c.Writer, c.Request)
	})
}

// Serve implements server.Server interface.
func (s *Server) Serve() error {
	//for _, route := range s.Engine.Routes() {
	//	s.config.logger.Info("add route", xlog.FieldMethod(route.Method), xlog.String("path", route.Path))
	//}
	s.Server = &http.Server{
		Addr:    s.config.Address(),
		Handler: s,
	}
	err := s.Server.Serve(s.listener)
	if err == http.ErrServerClosed {
		jlog.Info("close gin" + s.config.Address())
		return nil
	}

	return err
}

// Stop implements server.Server interface
// it will terminate gin server immediately
func (s *Server) Stop() error {
	return s.Server.Close()
}

// GracefulStop implements server.Server interface
// it will stop gin server gracefully
func (s *Server) GracefulStop(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}

// Info returns server info, used by governor and consumer balancer
func (s *Server) Info() *server.ServiceInfo {
	serviceAddr := s.listener.Addr().String()
	if s.config.ServiceAddress != "" {
		serviceAddr = s.config.ServiceAddress
	}

	info := server.ApplyOptions(
		server.WithScheme("http"),
		server.WithAddress(serviceAddr),
		server.WithKind(constant.ServiceProvider),
	)
	// info.Name = info.Name + "." + ModName
	return &info
}
