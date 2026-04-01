package main

import (
	"github.com/Dishank-Sen/quicnode/node"
	"github.com/Dishank-Sen/transport-config/config"
)

func getConfig(addr string) node.Config{
	return node.Config{
		ListenAddr: addr,
		TlsConfig: getTlsConfig(),
		QuicConfig: config.QuicCfg(),
	}
}