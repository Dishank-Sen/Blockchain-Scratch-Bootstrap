package peers

import (
	"fmt"
	"sync"
	"time"

	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/utils/logger"
)

type Store struct {
	mu    sync.RWMutex
	peers map[string]Peer
	order []string
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
		peers: make(map[string]Peer),
	}
	return store, nil
}

/*
Upsert:
- Called when a peer successfully registers
- Handles first connect and reconnect
*/
func (s *Store) Upsert(id string, addr string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// If new peer, track insertion order
	if _, exists := s.peers[id]; !exists {
		s.order = append(s.order, id)
	}

	s.peers[id] = Peer{
		ID:       id,
		Addr:     addr,
		LastSeen: time.Now().Unix(),
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

	if _, ok := s.peers[id]; !ok {
		return fmt.Errorf("no peer with id %s exists", id)
	}

	delete(s.peers, id)

	// Remove from order slice
	for i, pid := range s.order {
		if pid == id {
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
	for id, p := range s.peers {
		if id != peerID{
			out = append(out, p)
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

	for id, peer := range s.peers {
		if peer.Addr == addr {
			return id, true
		}
	}
	return "", false
}
