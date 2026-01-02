package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/cmd/controller"
	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/internal/server"
	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/utils/logger"
)


func handleRoutes(s *server.Server){
	s.Post("/register", controller.RegisterController)
	s.Get("/peers", controller.PeersController)
}

func main(){
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	port := ":4242"

	server := server.NewServer(ctx)
	logger.Info(fmt.Sprintf("quic server listening on port %s", port))

	tlsConf, err := getTlsConfig()
	if err != nil{
		logger.Error(err.Error())
		cancel()
	}

	quicConf, err := getQuicConfig()
	if err != nil{
		logger.Error(err.Error())
		cancel()
	}

	go handleRoutes(server)

	// blocks on Listen if no error
	if err := server.Listen(port, tlsConf, quicConf); err != nil{
		logger.Error(fmt.Sprintf("server listening error: %v", err))
		cancel()
	}
}

