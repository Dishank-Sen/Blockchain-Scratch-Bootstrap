package handler

import (
	"encoding/json"
	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/utils/logger"
	"github.com/Dishank-Sen/quicnode/types"
)

type ConnectPayload struct{
	ID string `json:"id"`
}

func (h *Handler) Connect(req *types.Request) *types.Response{
	conn := req.Conn
	var rp ConnectPayload
	if err := json.Unmarshal(req.Body, &rp); err != nil{
		return h.handleErrorRes()
	}

	h.store.Upsert(rp.ID, req.SourceAddr.String(), conn)
	peerList := h.store.GetAll(rp.ID)

	byteData, err := json.Marshal(peerList)
	if err != nil{
		logger.Debug("sending error response")
		return h.handleErrorRes()
	}
	return &types.Response{
		StatusCode: 200,
		Message: "ok",
		Headers: nil,
		Body: byteData,
	}
}