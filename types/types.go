package types

import "net"

type Peer struct{
	ID string `json:"id"`
	Addr string `json:"addr"`
}

type Peers struct{
	Users []Peer `json:"user"`
}


type Request struct{
	Method  string
	Path    string
	Headers map[string]string
	Addr    net.Addr
	Body    []byte
}

type Response struct {
	StatusCode int
	Message    string
	Headers    map[string]string
	Body       []byte
}