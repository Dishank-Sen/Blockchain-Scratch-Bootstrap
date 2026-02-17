package peers

import "github.com/quic-go/quic-go"

type Status string
const(
	CONNECTED Status = "connected"
	SUSPECT Status = "suspect"
	DEAD Status = "dead"
)

type Peer struct {
	ID       string
	Addr     string
	Conn     *quic.Conn
	LastSeen int64
	Status   Status
}