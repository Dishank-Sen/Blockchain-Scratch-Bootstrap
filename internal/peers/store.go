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
	peers map[string]*Peer // key = peer ID
	order []string         // insertion order by peer ID
}

// ---- global store state ----

var (
	store   *Store
	storeMu sync.Mutex
)

const max = 100

// GetStore returns a singleton store instance.
func GetStore() (*Store, error) {
	storeMu.Lock()
	defer storeMu.Unlock()

	if store != nil {
		return store, nil
	}

	store = &Store{
		peers: make(map[string]*Peer),
	}
	return store, nil
}

//
// Upsert
//
func (s *Store) Upsert(id string, addr string, conn *quic.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// If new peer, track insertion order
	if _, exists := s.peers[id]; !exists {
		s.order = append(s.order, id)
	}

	s.peers[id] = &Peer{
		ID:       id,
		Addr:     addr,
		Conn:     conn,
		LastSeen: time.Now().Unix(),
		Status:   "CONNECTED",
	}

	// Enforce max size
	if len(s.order) > max {
		oldestID := s.order[0]
		s.order = s.order[1:]
		delete(s.peers, oldestID)
	}
}

//
// Remove
//
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

//
// GetAll
//
func (s *Store) GetAll(excludeID string) []Peer {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]Peer, 0, len(s.peers))
	for id, peer := range s.peers {
		if id != excludeID {
			out = append(out, *peer)
		}
	}
	return out
}

//
// GetPeerConn
//
func (s *Store) GetPeerConn(id string) (*quic.Conn, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	peer, ok := s.peers[id]
	if !ok {
		return nil, fmt.Errorf("no peer id exists")
	}
	return peer.Conn, nil
}

//
// UpdateLastSeen
//
func (s *Store) UpdateLastSeen(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	peer, ok := s.peers[id]
	if !ok {
		return fmt.Errorf("no such peer")
	}

	peer.LastSeen = time.Now().Unix()
	return nil
}

//
// Cleanup
//
func (s *Store) Cleanup(ttl time.Duration) {
	now := time.Now()

	s.mu.Lock()
	defer s.mu.Unlock()

	for id, peer := range s.peers {
		if now.Sub(time.Unix(peer.LastSeen, 0)) > ttl {
			delete(s.peers, id)
			logger.Debug("peer deleted: " + id)
		}
	}
}

//
// DebugPrintAll
//
func (s *Store) DebugPrintAll() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.peers) == 0 {
		logger.Debug("no peers in store")
		return
	}

	for _, p := range s.peers {
		logger.Debug(fmt.Sprintf(
			"peer id=%s addr=%s last_seen=%d",
			p.ID,
			p.Addr,
			p.LastSeen,
		))
	}
}
