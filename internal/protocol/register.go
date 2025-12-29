package protocol

import (
	"context"
	"encoding/json"
	"fmt"

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
	if err := json.Unmarshal(payload, &rp); err != nil{
		return "", err
	}

	store := r.store
	store.Upsert(rp.ID, addr)

	
	store.DebugPrintAll()
	
	logger.Info(fmt.Sprintf("peer registered: %s", rp.ID))
	return rp.ID, nil
}

