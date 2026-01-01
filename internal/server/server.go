package server

import (
	"context"
	"crypto/tls"

	"github.com/quic-go/quic-go"
)

type Server struct {
	router *Router
}

func NewServer() *Server {
	return &Server{
		router: NewRouter(),
	}
}

func (s *Server) Get(path string, h HandlerFunc) {
	s.router.Get(path, h)
}

func (s *Server) Post(path string, h HandlerFunc) {
	s.router.Post(path, h)
}

func (s *Server) Listen(ctx context.Context, addr string, tlsCfg *tls.Config, quicConfig *quic.Config) error {
	listener, err := quic.ListenAddr(addr, tlsCfg, quicConfig)
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept(ctx)
		if err != nil {
			return err
		}
		go s.handleSession(ctx, conn)
	}
}
