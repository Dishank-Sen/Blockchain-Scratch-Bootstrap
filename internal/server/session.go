package server

import (
	"context"
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
			return
		}
		go s.handleStream(ctx, addr, stream)
	}
}

func cleanupPeerEntry(addr net.Addr){
	store, err := peers.GetStore("peers.json")
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