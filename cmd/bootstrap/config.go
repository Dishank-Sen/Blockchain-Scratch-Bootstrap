package main

import "github.com/Dishank-Sen/quicnode/node"

func getConfig(addr string) node.Config{
	return node.Config{
		ListenAddr: addr,
		TlsConfig: getTlsConfig(),
		QuicConfig: getQuicConfig(),
	}
}