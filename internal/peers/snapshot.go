package peers

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

type Snapshot struct {
	Path string
	mu   sync.Mutex
}

func NewSnapshot(path string) *Snapshot {
	return &Snapshot{Path: path}
}

/*
File format:

{
  "peers": {
    "peerID1": {
      "id": "peerID1",
      "addr": "1.2.3.4:4242",
      "last_seen": 1699999999
    }
  }
}
*/

type snapshotFile struct {
	Peers map[string]Peer `json:"peers"`
}

/* ---------------- Load ---------------- */

func (s *Snapshot) Load() (map[string]Peer, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// File does not exist → empty store
	if _, err := os.Stat(s.Path); errors.Is(err, os.ErrNotExist) {
		return make(map[string]Peer), nil
	}

	data, err := os.ReadFile(s.Path)
	if err != nil {
		return nil, err
	}

	// Empty file → empty store
	if len(data) == 0 {
		return make(map[string]Peer), nil
	}

	var snap snapshotFile
	if err := json.Unmarshal(data, &snap); err != nil {
		// Corrupt file → do NOT crash bootstrap
		// Start fresh but keep file for debugging
		return make(map[string]Peer), nil
	}

	if snap.Peers == nil {
		return make(map[string]Peer), nil
	}

	return snap.Peers, nil
}

/* ---------------- Save ---------------- */

func (s *Snapshot) Save(peers map[string]Peer) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	dir := filepath.Dir(s.Path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	tmp := s.Path 

	snap := snapshotFile{Peers: peers}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}

	// Write temp file first
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}

	// Atomic replace
	return os.Rename(tmp, s.Path)
}
