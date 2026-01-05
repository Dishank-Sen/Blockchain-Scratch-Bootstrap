package controller

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/internal/peers"
	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/types"
	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/utils/logger"
)

type peersList struct {
	ID   string `json:"id"`
	Addr string `json:"addr"`
}

func PeersController(ctx context.Context, req *types.Request) (*types.Response, error){
	logger.Debug("peer route")
	store, err := peers.GetStore()
	if err != nil{
		return handleErrorRes(err)
	}

	logger.Debug("all peers - peers.get.go - 25")
	store.DebugPrintAll()

	peers := store.GetAll()
	p := []peersList{}

	for _, peer := range peers{
		logger.Debug(fmt.Sprintf("id: %s | ip: %s", peer.ID, peer.Addr))
		p = append(p, peersList{
			ID: peer.ID,
			Addr: peer.Addr,
		})
	}

	byteData, err := json.Marshal(p)
	if err != nil{
		return handleErrorRes(err)
	}

	res := &types.Response{
		StatusCode: 200,
		Message: "ok",
		Headers: map[string]string{},
		Body: byteData,
	}
	return res, nil
}