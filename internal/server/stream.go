package server

import (
	"context"
	"net"

	"github.com/quic-go/quic-go"
)

func (s *Server) handleStream(ctx context.Context, addr net.Addr, stream *quic.Stream) {
	defer stream.Close()

	parser := NewParser(stream)
	req, err := parser.ParseRequest()
	req.Addr = addr  // attach peer address to req
	if err != nil {
		return
	}

	resp := s.router.Dispatch(ctx, req)
	if err = writeResponse(stream, resp); err != nil{
		stream.Close()
	}
}
