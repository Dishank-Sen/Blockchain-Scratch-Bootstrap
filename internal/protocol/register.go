package protocol

import (
	"context"
	"encoding/json"

	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/internal/peers"
	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/types"
	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/utils/logger"
)

type Register struct{
	store *peers.Store
}

func NewRegister(store *peers.Store) *Register{
	return &Register{
		store: store,
	}
}

func (r *Register) Handler(ctx context.Context, payload []byte, addr string) (string, error){
	var rp types.RegisterPayload
	if err := json.Unmarshal(payload, &r); err != nil{
		return "", err
	}

	store := r.store
	store.Upsert(rp.ID, addr)
	
	logger.Info("peer registered")
	return rp.ID, nil
}