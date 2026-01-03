package server

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/internal/peers"
	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/utils/logger"
	"github.com/quic-go/quic-go"
)

func (s *Server) handleSession(ctx context.Context, conn *quic.Conn) {
	go func ()  {
		<-ctx.Done()
		logger.Info("closing connection")
		cleanupPeerEntry(conn.RemoteAddr())
		conn.CloseWithError(0, "server shutdown")
	}()

	addr := conn.RemoteAddr()
	for {
		stream, err := conn.AcceptStream(ctx)
		if err != nil {
			handleConnClose(err, addr)
			return
		}
		go s.handleStream(ctx, addr, stream)
	}
}

func cleanupPeerEntry(addr net.Addr){
	store, err := peers.GetStore()
	if err != nil{
		logger.Error(err.Error())
	}
	id, ok := store.GetPeerIDByAddr(addr.String())
	if ok{
		if err := store.Remove(id); err != nil{
			logger.Error(err.Error())
		}
	}
}

func handleConnClose(err error, addr net.Addr) {
	logger.Warn(fmt.Sprintf("connection closed - addr: %s | error: %v", addr.String(), err))

	// quic-go exposes typed errors
	var (
		appErr    *quic.ApplicationError
		idleErr   *quic.IdleTimeoutError
		resetErr  *quic.StatelessResetError
	)

	switch {
	case errors.As(err, &idleErr):
		logger.Info(fmt.Sprintf("peer idle timeout - addr: %s", addr.String()))
		cleanupPeerEntry(addr)

	case errors.As(err, &resetErr):
		logger.Info(fmt.Sprintf("stateless reset - addr: %s", addr.String()))
		cleanupPeerEntry(addr)

	case errors.As(err, &appErr):
		logger.Info(fmt.Sprintf("application error - addr: %s | code: %v | error: %s", addr.String(), appErr.ErrorCode, appErr.ErrorMessage))
		cleanupPeerEntry(addr)

	case errors.Is(err, context.Canceled):
		// server shutdown — optional cleanup
		logger.Info(fmt.Sprintf("server context canceled - addr: %s", addr.String()))
		cleanupPeerEntry(addr)

	default:
		logger.Info(fmt.Sprintf("unknown connection close - addr: %s", addr.String()))
		cleanupPeerEntry(addr)
	}
}