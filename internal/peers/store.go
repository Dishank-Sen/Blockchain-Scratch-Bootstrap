package peers

import "sync"

type Store struct{
	mu sync.Mutex
	peers map[string]Peer
}

func NewStore() *Store{
	return &Store{}
}

func (s *Store) Add(peer Peer) error{

}

func (s *Store) Remove(peer Peer) error{

}

func (s *Store) Update(peer Peer) error{
	
}