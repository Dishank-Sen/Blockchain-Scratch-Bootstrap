package controller

import (
	"context"
	"encoding/json"

	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/internal/peers"
	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/types"
	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/utils/logger"
)

type RegisterPayload struct{
	ID string `json:"id"`
}

func RegisterController(ctx context.Context, req *types.Request) (*types.Response, error){
	logger.Debug("register route")
	logger.Debug(string(req.Body))
	logger.Debug(req.Addr.String())

	store, err := peers.GetStore()
	if err != nil{
		return handleErrorRes(err)
	}

	var rp RegisterPayload
	if err := json.Unmarshal(req.Body, &rp); err != nil{
		return handleErrorRes(err)
	}
	store.Upsert(rp.ID, req.Addr.String())

	res := &types.Response{
		StatusCode: 200,
		Message: "ok",
		Headers: map[string]string{},
		Body: []byte("registed"),
	}
	return res, nil
}

func handleErrorRes(err error) (*types.Response, error){
	errRes := &types.Response{
		StatusCode: 500,
		Message: err.Error(),
		Headers: map[string]string{},
		Body: nil,
	}
	return errRes, err
}