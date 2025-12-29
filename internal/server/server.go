package server

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/internal/peers"
	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/utils/logger"
	"github.com/quic-go/quic-go"
)

type Server struct{
	tlsConfig *tls.Config
	quicConfig *quic.Config
	addr string
	listener *quic.Listener
	ctx context.Context
	cancel context.CancelFunc
	store *peers.Store
}

func NewServer(ctx context.Context, store *peers.Store) (*Server, error){
	tlsConfig, err := GetTlsConfig()
	if err != nil{
		return nil, err
	}

	quicConfig, err := GetQuicConfig()
	if err != nil{
		return nil, err
	}

	server := &Server{
		tlsConfig: tlsConfig,
		quicConfig: quicConfig,
		store: store,
		ctx: ctx,
	}

	return server, nil
}

func (s *Server) Start(addr string) error {
	s.addr = addr

	// OS signal handling INSIDE server
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		logger.Warn("received signal: " + sig.String())
		s.cancel() // 🔑 trigger graceful shutdown
	}()

	return s.listen()
}

func (s *Server) listen() error{
	defer func() {
		if err := s.cleanup(); err != nil {
			logger.Warn(fmt.Sprintf("cleanup failed: %v", err))
		}
	}()

	serverCtx, serverCancel := context.WithCancel(s.ctx)
	defer serverCancel()

	listener, err := quic.ListenAddr(s.addr, s.tlsConfig, s.quicConfig)
	if err != nil{
		return fmt.Errorf("failed to listen: %v", err)
	}

	
	s.listener = listener
	s.cancel = serverCancel
	s.ctx = serverCtx

	for {
		sess, err := listener.Accept(s.ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				logger.Info("server shutdown complete")
				return nil
			}
			return fmt.Errorf("listener error: %v", err)
		}
		conn := sess
		go func (conn *quic.Conn)  {
			addr := sess.RemoteAddr()
			session := NewSession(s.ctx, conn, addr, s.store)
			session.Handle() // blocking function
		}(conn)
	}
}

// pushes all the peers to json file
func (s *Server) cleanup() error{
	store := s.store
	if err := store.Cleanup(); err != nil{
		return err
	}
	return nil
}