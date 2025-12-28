package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/types"
	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/utils/logger"
	"github.com/quic-go/quic-go"
)

// alter I have to add mutex also

func loadCert() *tls.Config{
	certFilePath := path.Join("certificate", "server.crt")
	keyFilePath := path.Join("certificate", "server.key")

	cert, err := tls.LoadX509KeyPair(certFilePath, keyFilePath)

	if err != nil{
		panic(err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"quic-example-v1"},
	}
}

func handleSession(sess *quic.Conn) {
	defer func() {
		_ = sess.CloseWithError(0, "closing")
	}()

	remoteAddr := sess.RemoteAddr()
	logger.Info(fmt.Sprintf(
		"New session from %s (IP=%s, Port=%d)",
		remoteAddr.String(),
		remoteAddr.(*net.UDPAddr).IP,
		remoteAddr.(*net.UDPAddr).Port,
	))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Accept streams in a loop
	for {
		stream, err := sess.AcceptStream(ctx)
		if err != nil {
			// AcceptStream returns non-nil err when session is closed
			logger.Info(fmt.Sprintf("AcceptStream error (%s): %v\n", remoteAddr, err))
			return
		}
		go handleStream(stream, remoteAddr)
	}
}

func handleStream(s *quic.Stream, addr net.Addr) {
	defer s.Close()

	// Simple single-read echo example; production should loop and frame messages
	data, err := io.ReadAll(s)
	if err != nil {
		logger.Info(fmt.Sprintf("stream read error: %v", err))
		return
	}

	// break the bytes into different stream section
	var stream types.StreamMessage
	if err := json.Unmarshal(data, &stream); err != nil{
		logger.Info(fmt.Sprintf("error while processing stream bytes: %v", err))
	}

	logger.Info(fmt.Sprintf("version: %d | type: %s | length: %d\n", stream.Version, stream.Type, stream.Length))

	logger.Info("header: \n")
	for k, p := range stream.Header{
		fmt.Printf("%s: %s\n", k, p)
	}

	payload := string(stream.Payload)
	logger.Info(fmt.Sprintf("payload: %s", payload))

	if _, ok := stream.Header["application/json"]; !ok{
		logger.Info("not a json content")
		return
	}

	switch stream.Type{
	case "register":
		if err := handleRegister(stream.Payload, addr); err != nil {
			logger.Info(fmt.Sprintf("register failed: %v", err))
		}
	case "punch":
		handlePunch()
	case "ping":
		handlePing()
	default:
		handleDefault()
	}
	// Echo back with a timestamp
	_, _ = s.Write([]byte("pong: " + time.Now().Format(time.RFC3339)))
}

func handleRegister(payload []byte, addr net.Addr) error{
	logger.Info("register case")

	var r types.RegisterPayload
	if err := json.Unmarshal(payload, &r); err != nil{
		return err
	}

	if err := savePeer(r.ID, addr); err != nil{
		return err
	}
	
	logger.Info("peer added")
	return nil
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

func loadPeers(path string) (types.Peers, error) {
	var peers types.Peers

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return types.Peers{Users: []types.Peer{}}, nil
		}
		return peers, err
	}

	if len(data) == 0 {
		return types.Peers{Users: []types.Peer{}}, nil
	}

	if err := json.Unmarshal(data, &peers); err != nil {
		return peers, err
	}

	return peers, nil
}


func savePeer(id string, addr net.Addr) error {
	path := "peers.json"

	peers, err := loadPeers(path)
	if err != nil {
		return err
	}

	peers.Users = append(peers.Users, types.Peer{
		ID:   id,
		Addr: addr,
	})

	out, err := json.MarshalIndent(peers, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, out, 0644)
}


func main(){
	tlsConf := loadCert()

	quicConf := &quic.Config{
		MaxIdleTimeout: 30 * time.Second,
	}

	addr := ":4242" // UDP port
	listener, err := quic.ListenAddr(addr, tlsConf, quicConf)

	if err != nil {
		logger.Info(fmt.Sprintf("failed to listen: %v", err))
	}
	logger.Info(fmt.Sprintf("QUIC server listening on %s\n", addr))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Accept sessions forever
	for {
		sess, err := listener.Accept(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				logger.Info("server shutdown complete")
				return
			}
			logger.Info(fmt.Sprintf("listener error: %v\n", err))
			return
		}
		go handleSession(sess)
	}
}