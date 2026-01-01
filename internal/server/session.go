package server

import (
	"context"

	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/utils/logger"
	"github.com/quic-go/quic-go"
)

func (s *Server) handleSession(ctx context.Context, conn *quic.Conn) {
	go func ()  {
		<-ctx.Done()
		logger.Info("closing connection")
		conn.CloseWithError(0, "server shutdown")
	}()

	for {
		stream, err := conn.AcceptStream(ctx)
		if err != nil {
			return
		}
		go s.handleStream(ctx, stream)
	}
}
