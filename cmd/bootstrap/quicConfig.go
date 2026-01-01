package main

import (
	"time"

	"github.com/quic-go/quic-go"
)

func getQuicConfig() (*quic.Config, error){
	quicConf := &quic.Config{
		MaxIdleTimeout: 30 * time.Second,
	}

	return quicConf, nil
}