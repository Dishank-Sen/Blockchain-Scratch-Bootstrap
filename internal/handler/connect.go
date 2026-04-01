package handler

import (
	"context"
	"encoding/json"

	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/internal/peers"
	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/utils/logger"
	"github.com/Dishank-Sen/quicnode/types"
)

type ConnectPayload struct{
	ID string `json:"id"`
}

type PeerInfo struct{
	ID string `json:"id"`
	Addr string `json:"addr"`
}

func (h *Handler) Connect(ctx context.Context, req *types.Request) *types.Response{
	conn := req.Conn
	var rp ConnectPayload
	if err := json.Unmarshal(req.Body, &rp); err != nil{
		return h.handleErrorRes()
	}

	connID := ctx.Value("connID").(types.ConnID)
	h.store.Upsert(rp.ID, req.SourceAddr.String(), conn, connID)
	peerList := h.store.GetAll(rp.ID)

	byteData, err := json.Marshal(peerList)
	if err != nil{
		logger.Debug("sending error response")
		return h.handleErrorRes()
	}

	// Dial to the connected peers
	dialPayload := PeerInfo{
		ID: rp.ID,
		Addr: req.SourceAddr.String(),
	}
	go h.dialPeer(peerList, dialPayload)

	return &types.Response{
		StatusCode: 200,
		Message: "ok",
		Headers: nil,
		Body: byteData,
	}
}

func (h *Handler) dialPeer(peersList []peers.Peer, payload PeerInfo){
	if len(peersList) == 0{
		logger.Info("no peers to dial")
		return
	}
	for _, peer := range peersList{
		byteData, err := json.Marshal(payload)
		if err != nil{
			logger.Error(err.Error())
			continue
		}
		logger.Info("dialing (accept-peers)...")
		resp, err := h.node.Dial(peer.Addr, "accept-peers", nil, byteData)
		if err != nil{
			logger.Error(err.Error())
			continue
		}
		logger.Info(string(resp.Body))
	}
}