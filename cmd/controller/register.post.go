package controller

import (
	"context"

	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/types"
	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/utils/logger"
)

func RegisterController(ctx context.Context, req *types.Request) (*types.Response, error){
	logger.Debug("register route")
	logger.Debug(string(req.Body))
	res := &types.Response{
		StatusCode: 200,
		Message: "ok",
		Headers: map[string]string{},
		Body: []byte("registed"),
	}
	return res, nil
}