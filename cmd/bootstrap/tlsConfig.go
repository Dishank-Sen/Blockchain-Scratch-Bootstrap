package main

import (
	"crypto/tls"
	"path"

	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/utils/logger"
)

func getTlsConfig() *tls.Config{
	certFilePath := path.Join("certificate", "server.crt")
	keyFilePath := path.Join("certificate", "server.key")

	cert, err := tls.LoadX509KeyPair(certFilePath, keyFilePath)
	if err != nil{
		logger.Debug("error in tls")
		logger.Error(err.Error())
		return nil
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"quicnode"},
	}
	return tlsConfig
}