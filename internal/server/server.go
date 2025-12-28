package server

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/internal/peers"
	"github.com/quic-go/quic-go"
)

type Server struct{
	tlsConfig *tls.Config
	quicConfig *quic.Config
	addr string
	listener *quic.Listener
	ctx context.Context
	cancel context.CancelFunc
	peers *peers.Store
}

func NewServer() (*Server, error){
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
	}

	return server, nil
}

func (s *Server) Listen(ctx context.Context, cancel context.CancelFunc, addr string) error{
	listener, err := quic.ListenAddr(addr, s.tlsConfig, s.quicConfig)
	if err != nil{
		return fmt.Errorf("failed to listen: %v", err)
	}

	s.addr = addr
	s.ctx = ctx
	s.cancel = cancel
	s.listener = listener

	for {
		sess, err := listener.Accept(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return fmt.Errorf("context cancelled error: %v", err)
			}
			return fmt.Errorf("listener error: %v", err)
		}
		go handleSession(sess)
	}
}