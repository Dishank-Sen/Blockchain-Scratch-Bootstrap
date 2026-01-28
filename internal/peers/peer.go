package peers

type Status string
const(
	CONNECTED Status = "connected"
	SUSPECT Status = "suspect"
	DEAD Status = "dead"
)

type Peer struct {
	ID       string `json:"id"`
	Addr     string `json:"addr"`
	LastSeen int64  `json:"last_seen"`
	Status   Status `json:"status"`
}