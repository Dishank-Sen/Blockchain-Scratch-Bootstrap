package peers

import (
	"fmt"
	"sync"
	"time"

	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/utils/logger"
	"github.com/quic-go/quic-go"
)

type Store struct {
	mu    sync.RWMutex
	peers map[*quic.Conn]Peer
	order []*quic.Conn
}

// ---- global store state ----

var (
	store   *Store
	storeMu sync.Mutex
)

const (
	max = 100
)
// GetStore returns a singleton store instance.
// This MUST always return the same store.
func GetStore() (*Store, error) {
	storeMu.Lock()
	defer storeMu.Unlock()

	if store != nil {
		return store, nil
	}

	store = &Store{
		peers: make(map[*quic.Conn]Peer),
	}
	return store, nil
}

/*
Upsert:
- Called when a peer successfully registers
- Handles first connect and reconnect
*/
func (s *Store) Upsert(id string, addr string, conn *quic.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// If new peer, track insertion order
	if _, exists := s.peers[conn]; !exists {
		s.order = append(s.order, conn)
	}

	s.peers[conn] = Peer{
		ID:       id,
		Addr:     addr,
		LastSeen: time.Now().Unix(),
		Status: "CONNECTED",
	}

	// Enforce max size
	if len(s.order) > max {
		oldest := s.order[0]
		s.order = s.order[1:]
		delete(s.peers, oldest)
	}
}

/*
Remove:
- Called when a session ends
- Peer is considered offline
*/
func (s *Store) Remove(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	var conn *quic.Conn

	for c, p := range s.peers{
		if p.ID == id{
			conn = c
			break
		}
	}
	if conn == nil{
		return fmt.Errorf("no peer exist to remove")
	}
	if _, ok := s.peers[conn]; !ok {
		return fmt.Errorf("no peer with id %s exists", id)
	}

	delete(s.peers, conn)

	// Remove from order slice
	for i, c := range s.order {
		if c == conn {
			s.order = append(s.order[:i], s.order[i+1:]...)
			break
		}
	}

	return nil
}


/*
GetAll:
- Returns a recent peer list excluding the peer requested
*/
func (s *Store) GetAll(peerID string) []Peer {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]Peer, 0, len(s.peers))
	for _, peer := range s.peers {
		if peer.ID != peerID{
			out = append(out, peer)
		}
	}
	return out
}

// DebugPrintAll prints current store state (debug only)
func (s *Store) DebugPrintAll() {
	peers := s.GetAll("")

	if len(peers) == 0 {
		logger.Debug("no peers in store")
		return
	}

	for _, p := range peers {
		logger.Debug(fmt.Sprintf(
			"peer id=%s addr=%s last_seen=%d",
			p.ID,
			p.Addr,
			p.LastSeen,
		))
	}
}

/*
GetPeerIDByAddr:
- Used during connection cleanup
- Note: addr is unstable across reconnects (design limitation)
*/
func (s *Store) GetPeerIDByAddr(addr string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, peer := range s.peers {
		if peer.Addr == addr {
			return peer.ID, true
		}
	}
	return "", false
}

func (s *Store) UpdateLastSeen(conn *quic.Conn) error{
	updatedPeer, ok := s.peers[conn]
	if !ok{
		return fmt.Errorf("no such connection")
	}
	updatedPeer.LastSeen = time.Now().Unix()
	s.peers[conn] = updatedPeer
	return nil
}

func (s *Store) Cleanup(ttl time.Duration) {
    now := time.Now()

    s.mu.Lock()
    defer s.mu.Unlock()

    for conn, peer := range s.peers {
        if now.Sub(time.Unix(peer.LastSeen, 0)) > ttl {
            delete(s.peers, conn)
			logger.Debug("peer deleted")
        }
    }
}
