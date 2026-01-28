package handler

import (
	"time"

	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/internal/peers"
	"github.com/Dishank-Sen/quicnode/node"
	"github.com/Dishank-Sen/quicnode/types"
)

type Handler struct{
	node *node.Node
	store *peers.Store
}

func NewHandler(n *node.Node) (*Handler, error){
	store, err := peers.GetStore()
	if err != nil{
		return nil, err
	}
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			store.Cleanup(15 * time.Second)
		}
	}()
	return &Handler{
		node: n,
		store: store,
	}, nil
}

func (h *Handler) handleErrorRes() *types.Response{
	errRes := &types.Response{
		StatusCode: 500,
		Message: "Error",
		Headers: map[string]string{},
		Body: []byte("Internal Server Error"),
	}
	return errRes
}