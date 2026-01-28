package main

import (
	"time"

	"github.com/quic-go/quic-go"
)

func getQuicConfig() *quic.Config{
	quicConf := &quic.Config{
		MaxIdleTimeout: 30*time.Second,
		KeepAlivePeriod: 10*time.Second,
	}

	return quicConf
}