package server

import (
	"time"

	"github.com/quic-go/quic-go"
)

func GetQuicConfig() (*quic.Config, error){
	quicConf := &quic.Config{
		MaxIdleTimeout: 30 * time.Second,
	}

	return quicConf, nil
}