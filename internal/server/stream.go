package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"

	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/internal/peers"
	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/internal/protocol"
	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/types"
	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/utils/logger"
	"github.com/quic-go/quic-go"
)

type Stream struct{
	store *peers.Store
	stream *quic.Stream
	addr net.Addr
	ctx context.Context
	cancel context.CancelFunc
}

func NewStream(sessionCtx context.Context, store *peers.Store, stream *quic.Stream, addr net.Addr) *Stream{
	streamCtx, streamCancel := context.WithCancel(sessionCtx)

	return &Stream{
		store: store,
		stream: stream,
		addr: addr,
		ctx: streamCtx,
		cancel: streamCancel,
	}
}

func (s *Stream) Handle(msg *types.StreamMessage) error{
	defer s.cancel()

	switch msg.Type{
	case "register":
		return fmt.Errorf("already registered")
	case "punch":
		handlePunch()
	case "ping":
		handlePing()
	default:
		handleDefault()
	}
	
	return nil
}

func (s *Stream) HandleRegister(msg *types.StreamMessage) (string, error){
	if _, ok := msg.Header["application/json"]; !ok {
		return "", fmt.Errorf("not json")
	}
	
	register := protocol.NewRegister(s.store)
	id, err := register.Handler(s.ctx, msg.Payload, s.addr.String())
	if err != nil{
		return "", fmt.Errorf("error in registering: %v", err)
	}
	return id, nil
}

func (s *Stream) ReadMessage() (*types.StreamMessage, error) {
	data, err := io.ReadAll(s.stream)
	if err != nil {
		return nil, err
	}

	var msg types.StreamMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}


func handlePunch(){
	logger.Info("punch case")
}

func handlePing(){
	logger.Info("ping case")
}

func handleDefault(){
	logger.Info("default case")
} 