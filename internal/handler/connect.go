package handler

import (
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

	// Dial to the connected peers
	go h.dialPeer(peerList)

	return &types.Response{
		StatusCode: 200,
		Message: "ok",
		Headers: nil,
		Body: byteData,
	}
}

func (h *Handler) dialPeer(peersList []peers.Peer){
	if len(peersList) == 0{
		logger.Info("no peers to dial")
		return
	}
	for _, peer := range peersList{
		conn, err := h.store.GetPeerConn(peer.ID)
		if err != nil{
			logger.Error(err.Error())
			continue
		}
		peerInfo := PeerInfo{
			ID: peer.ID,
			Addr: peer.Addr,
		}
		byteData, err := json.Marshal(peerInfo)
		if err != nil{
			logger.Error(err.Error())
			continue
		}
		resp, err := h.node.DialConn(conn, "accept-peers", nil, byteData)
		if err != nil{
			logger.Error(err.Error())
			continue
		}
		logger.Info(string(resp.Body))
	}
}