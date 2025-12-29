package peers

type Peer struct {
	ID       string `json:"id"`
	Addr     string `json:"addr"`
	LastSeen int64  `json:"last_seen"`
}
