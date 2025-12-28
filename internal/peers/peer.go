package peers

import "time"

type Peer struct {
    ID        string
    Addr      string
    LastSeen  time.Time
}
