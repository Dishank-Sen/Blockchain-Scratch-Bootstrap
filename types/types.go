package types

import "net"

type Peer struct{
	ID string `json:"id"`
	Addr net.Addr `json:"addr"`
}

type Peers struct{
	Users []Peer `json:"user"`
}

type MessageType string

const (
    Register MessageType = "register"
    Punch    MessageType = "punch"
    Ping     MessageType = "ping"
)


type StreamMessage struct{
	Version uint16 `json:"version"`
	Header  map[string]string `json:"header"`
    Type    MessageType `json:"type"`
    Length  uint32 `json:"length"`
    Payload []byte `json:"payload"`
}

type RegisterPayload struct{
	ID string `json:"id"`
}