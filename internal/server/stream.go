package server

import (
	"context"
	"encoding/binary"
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

const maxMessageSize = 1 << 20 // 1 MB

type Stream struct {
	store  *peers.Store
	stream *quic.Stream
	addr   net.Addr
	ctx    context.Context
}

func NewStream(sessionCtx context.Context, store *peers.Store, stream *quic.Stream, addr net.Addr) *Stream {
	return &Stream{
		store:  store,
		stream: stream,
		addr:   addr,
		ctx:    sessionCtx,
	}
}

/* -------------------- READ -------------------- */

func (s *Stream) ReadMessage() (*types.StreamMessage, error) {
	// Abort immediately if session is shutting down
	select {
	case <-s.ctx.Done():
		return nil, s.ctx.Err()
	default:
	}

	// 1️⃣ Read fixed-size length header (4 bytes)
	var lenBuf [4]byte
	if _, err := io.ReadFull(s.stream, lenBuf[:]); err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint32(lenBuf[:])
	if length == 0 || length > maxMessageSize {
		return nil, fmt.Errorf("invalid message length: %d", length)
	}

	// 2️⃣ Read exactly `length` bytes
	payload := make([]byte, length)
	if _, err := io.ReadFull(s.stream, payload); err != nil {
		return nil, err
	}

	// 3️⃣ Decode message
	var msg types.StreamMessage
	if err := json.Unmarshal(payload, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

/* -------------------- HANDLE -------------------- */

func (s *Stream) Handle(msg *types.StreamMessage) error {
	// Abort if session is canceled
	select {
	case <-s.ctx.Done():
		return s.ctx.Err()
	default:
	}

	switch msg.Type {

	case "register":
		return fmt.Errorf("already registered")

	case "ping":
		return s.writePong()

	case "punch":
		logger.Info("punch received")
		return nil

	default:
		logger.Info("unknown message type")
		return nil
	}
}

/* -------------------- REGISTER -------------------- */

func (s *Stream) HandleRegister(msg *types.StreamMessage) (string, error) {
	if _, ok := msg.Header["application/json"]; !ok {
		return "", fmt.Errorf("invalid content type")
	}

	register := protocol.NewRegister(s.store)
	id, err := register.Handler(s.ctx, msg.Payload, s.addr.String())
	if err != nil {
		return "", err
	}

	return id, nil
}

/* -------------------- WRITE HELPERS -------------------- */

func (s *Stream) writePong() error {
	resp := []byte("pong")

	var lenBuf [4]byte
	binary.BigEndian.PutUint32(lenBuf[:], uint32(len(resp)))

	if _, err := s.stream.Write(lenBuf[:]); err != nil {
		return err
	}
	if _, err := s.stream.Write(resp); err != nil {
		return err
	}
	return nil
}
