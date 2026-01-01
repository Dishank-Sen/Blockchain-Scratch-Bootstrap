package server

import (
	"context"

	"github.com/quic-go/quic-go"
)

func (s *Server) handleStream(ctx context.Context, stream *quic.Stream) {
	defer stream.Close()

	parser := NewParser(stream)
	req, err := parser.ParseRequest()
	if err != nil {
		return
	}

	resp := s.router.Dispatch(ctx, req)
	_ = writeResponse(stream, resp)
}
