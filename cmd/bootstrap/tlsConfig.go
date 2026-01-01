package main

import (
	"crypto/tls"
	"path"
)

func getTlsConfig() (*tls.Config, error){
	certFilePath := path.Join("certificate", "server.crt")
	keyFilePath := path.Join("certificate", "server.key")

	cert, err := tls.LoadX509KeyPair(certFilePath, keyFilePath)
	if err != nil{
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"quic-example-v1"},
	}
	return tlsConfig, nil
}