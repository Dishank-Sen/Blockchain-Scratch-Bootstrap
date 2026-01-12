package handler

import (
	"encoding/json"
	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/internal/peers"
	"github.com/Dishank-Sen/quicnode/types"
)

type ConnectPayload struct{
	ID string `json:"id"`
}

func (h *Handler) Connect(req *types.Request) *types.Response{
	store, err := peers.GetStore()
	if err != nil{
		return h.handleErrorRes()
	}

	var rp ConnectPayload
	if err := json.Unmarshal(req.Body, &rp); err != nil{
		return h.handleErrorRes()
	}
	store.Upsert(rp.ID, req.SourceAddr.String())
	peerList := store.GetAll(rp.ID)

	byteData, err := json.Marshal(peerList)
	if err != nil{
		return h.handleErrorRes()
	}
	return &types.Response{
		StatusCode: 200,
		Message: "ok",
		Headers: nil,
		Body: byteData,
	}
}

