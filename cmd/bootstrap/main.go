package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/internal/router"
	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/utils/logger"
	"github.com/Dishank-Sen/quicnode/node"
)

func main(){
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	addr := ":4242"

	cfg := getConfig(addr)
	n, err := node.NewNode(ctx, cfg)
	if err != nil{
		logger.Error(err.Error())
		cancel()
	}
	if err := n.Start(); err != nil{
		logger.Error(err.Error())
		cancel()
		n.Stop()
	}
	logger.Info("node started")
	router := router.NewRouter(n)
	go router.HandleRoutes()
	<-ctx.Done()
}