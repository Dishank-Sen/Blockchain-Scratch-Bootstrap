package peers

import (
	"fmt"
	"sync"
	"time"

	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/utils/logger"
)

type Store struct {
	mu    sync.Mutex
	peers map[string]Peer
	snap  *Snapshot
}

// ---- global store state ----

var (
	store   *Store
	storeMu sync.Mutex
)

/*
GetStore:
- Returns existing store if already created
- Otherwise creates a new store
- Safe for concurrent callers
- Retries creation if previous attempt failed
*/
func GetStore(snapshotPath string) (*Store, error) {
	storeMu.Lock()
	defer storeMu.Unlock()

	if store != nil {
		return store, nil
	}

	snap := NewSnapshot(snapshotPath)

	peers, err := snap.Load()
	if err != nil {
		return nil, err
	}

	store = &Store{
		peers: peers,
		snap:  snap,
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

func (s *Store) DebugPrintAll() {
	peers := s.GetAll()

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

func (s *Store) GetPeerIDByAddr(addr string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, peer := range s.peers {
		if peer.Addr == addr {
			return id, true
		}
	}
	return "", false
}
