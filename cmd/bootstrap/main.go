package main

import (
	"context"
	"fmt"

	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/internal/peers"
	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/internal/server"
	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/utils/logger"
)


func main(){
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	port := ":4242"
	store, err := peers.NewStore("peers.json")
	if err != nil{
		logger.Error(fmt.Sprintf("store error: %v", err))
		cancel()
	}

	quicServer, err := server.NewServer(ctx, store)
	if err != nil{
		logger.Error(fmt.Sprintf("quic server error: %v", err))
		cancel()
	}

	logger.Info(fmt.Sprintf("quic server listening on port %s", port))

	// blocks on Listen if no error
	if err := quicServer.Start(port); err != nil{
		logger.Error(fmt.Sprintf("server listening error: %v", err))
		cancel()
	}
}