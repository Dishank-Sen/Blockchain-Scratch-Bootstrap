package server

import (
	"context"
	"crypto/tls"

	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/utils/logger"
	"github.com/quic-go/quic-go"
)

type Server struct {
	router *Router
	ctx context.Context
	cancel context.CancelFunc
}

func NewServer(ctx context.Context) *Server {
	svrCtx, svrCancel := context.WithCancel(ctx)

	return &Server{
		router: NewRouter(),
		ctx: svrCtx,
		cancel: svrCancel,
	}
}

func (s *Server) Get(path string, h HandlerFunc) {
	s.router.Get(path, h)
}

func (s *Server) Post(path string, h HandlerFunc) {
	s.router.Post(path, h)
}

func (s *Server) Listen(addr string, tlsCfg *tls.Config, quicConfig *quic.Config) error {
	listener, err := quic.ListenAddr(addr, tlsCfg, quicConfig)
	if err != nil {
		return err
	}

	go func ()  {
		<-s.ctx.Done()
		logger.Info("server shutting down")
		listener.Close()
	}()

	for {
		logger.Debug("waiting for session...")
		conn, err := listener.Accept(s.ctx)
		if err != nil {
			logger.Debug("error at server.go - 50")
			logger.Error(err.Error())
			if s.ctx.Err() != nil{
				return nil
			}
			return err
		}
		logger.Debug("session received")
		go s.handleSession(s.ctx, conn)
	}
}
