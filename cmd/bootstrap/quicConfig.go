package main

import (
	"time"

	"github.com/quic-go/quic-go"
)

func getQuicConfig() *quic.Config{
	quicConf := &quic.Config{
		MaxIdleTimeout: 60 * time.Minute,
	}

	return quicConf
}