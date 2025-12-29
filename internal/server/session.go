package server

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/internal/peers"
	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/utils/logger"
	"github.com/quic-go/quic-go"
)

type Session struct {
	conn   *quic.Conn
	peerID string
	addr net.Addr
	store  *peers.Store
	ctx context.Context
	cancel context.CancelFunc
	mu sync.Mutex
}

func NewSession(serverCtx context.Context, conn *quic.Conn, addr net.Addr, store *peers.Store) *Session{
	sessionCtx, sessionCancel := context.WithCancel(serverCtx)

	return &Session{
		conn: conn,
		addr: addr,
		store: store,
		ctx: sessionCtx,
		cancel: sessionCancel,
	}
}

// blocking function
func (s *Session) Handle(){
	defer s.cancel()
	defer s.cleanup()
	defer func() {
		_ = s.conn.CloseWithError(0, "closing")
	}()
		
	remoteAddr := s.conn.RemoteAddr()
	logger.Info(fmt.Sprintf(
		"New session from %s (IP=%s, Port=%d)",
		remoteAddr.String(),
		remoteAddr.(*net.UDPAddr).IP,
		remoteAddr.(*net.UDPAddr).Port,
	))

	// Accept streams in a loop
	for {
		stream, err := s.conn.AcceptStream(s.ctx)
		if err != nil {
			// AcceptStream returns non-nil err when session is closed
			logger.Info(fmt.Sprintf("AcceptStream error (%s): %v\n", remoteAddr, err))
			return
		}
		st := stream // capture
		go s.handleStream(st)
	}
}

func (s *Session) cleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.peerID != "" {
		s.store.Remove(s.peerID)
		logger.Info(fmt.Sprintf("peer removed: %s", s.peerID))
	}
}

func (s *Session) handleStream(st *quic.Stream) {
	defer st.Close()

	stream := NewStream(s.ctx, s.store, st, s.addr)

	msg, err := stream.ReadMessage()
	if err != nil {
		logger.Info(fmt.Sprintf("stream read error: %v", err))
		return
	}

	s.mu.Lock()
	peerID := s.peerID
	s.mu.Unlock()

	if peerID != "" && msg.Type == "register" {
		logger.Info("protocol violation: duplicate register")
		s.cancel()
		return
	}

	// --- AUTH GATE ---
	if peerID == "" {
		if msg.Type != "register" {
			logger.Info("protocol violation: first stream must be register")
			s.cancel() // kill entire session
			return
		}

		id, err := stream.HandleRegister(msg)
		if err != nil {
			logger.Info(fmt.Sprintf("register failed: %v", err))
			s.cancel()
			return
		}

		s.mu.Lock()
		s.peerID = id
		s.mu.Unlock()

		return
	}

	// --- AUTHENTICATED ---
	if err := stream.Handle(msg); err != nil {
		logger.Info(fmt.Sprintf("stream error: %v", err))
	}
}
