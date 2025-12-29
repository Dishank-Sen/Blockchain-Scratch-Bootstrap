package peers

import (
	"fmt"
	"sync"
	"time"
)

type Store struct {
	mu    sync.Mutex
	peers map[string]Peer
	snap  *Snapshot
}

/*
NewStore:
- Loads peers from snapshot (crash recovery)
- Never fails if snapshot is missing / empty / corrupt
*/
func NewStore(snapshotPath string) (*Store, error) {
	snap := NewSnapshot(snapshotPath)

	peers, err := snap.Load()
	if err != nil {
		return nil, err
	}

	return &Store{
		peers: peers,
		snap:  snap,
	}, nil
}

/*
Upsert:
- Called when a peer successfully registers
- Handles first connect and reconnect
*/
func (s *Store) Upsert(id string, addr string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.peers[id] = Peer{
		ID:       id,
		Addr:     addr,
		LastSeen: time.Now().Unix(),
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
	return nil
}

/*
GetAll:
- Returns a snapshot copy (safe for readers)
- Useful for bootstrap peer list responses
*/
func (s *Store) GetAll() []Peer {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]Peer, 0, len(s.peers))
	for _, p := range s.peers {
		out = append(out, p)
	}
	return out
}

/*
Cleanup:
- Called on graceful shutdown
- Persists in-memory state to peers.json
*/
func (s *Store) Cleanup() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.snap.Save(s.peers)
}
