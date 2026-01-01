package server

import (
	"context"

	"github.com/quic-go/quic-go"
)

func (s *Server) handleSession(ctx context.Context, conn *quic.Conn) {
	for {
		stream, err := conn.AcceptStream(ctx)
		if err != nil {
			return
		}
		go s.handleStream(ctx, stream)
	}
}
