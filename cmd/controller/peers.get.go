package controller

import (
	"context"
	"encoding/json"
	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/internal/peers"
	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/types"
)

type peersList struct {
	ID   string `json:"id"`
	Addr string `json:"addr"`
}

func PeersController(ctx context.Context, req *types.Request) (*types.Response, error){
	store, err := peers.GetStore("peers.json")
	if err != nil{
		return handleErrorRes(err)
	}

	peers := store.GetAll()
	var p []peersList
	for _, peer := range peers{
		// logger.Debug(fmt.Sprintf("id: %s | ip: %s", peer.ID, peer.Addr))
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