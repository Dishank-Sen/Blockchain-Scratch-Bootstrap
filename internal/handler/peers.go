package handler

import (
	"encoding/json"
	"github.com/Dishank-Sen/quicnode/types"
)

type peersList struct {
	ID   string `json:"id"`
	Addr string `json:"addr"`
}

type payload struct{
	ID string `json:"id"`
}

func (h *Handler) Peers(req *types.Request) *types.Response{
	var p payload
	if err := json.Unmarshal(req.Body, &p); err != nil{
		return h.handleErrorRes()
	}
	
	nodeID := p.ID
	peerL := []peersList{}
	peers := h.store.GetAll(nodeID)

	for _, peer := range peers{
		peerL = append(peerL, peersList{
			ID: peer.ID,
			Addr: peer.Addr,
		})
	}

	byteData, err := json.Marshal(peerL)
	if err != nil{
		return h.handleErrorRes()
	}

	return &types.Response{
		StatusCode: 200,
		Message: "ok",
		Headers: map[string]string{},
		Body: byteData,
	}
}