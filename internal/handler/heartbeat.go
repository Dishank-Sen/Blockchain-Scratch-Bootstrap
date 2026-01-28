package handler

import (
	"fmt"

	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/utils/logger"
	"github.com/Dishank-Sen/quicnode/types"
)

func (h *Handler) Heartbeat(req *types.Request) *types.Response{
	logger.Info(fmt.Sprintf("heartbeat from %s", req.SourceAddr))
	if err := h.store.UpdateLastSeen(req.Conn); err != nil{
		logger.Error(err.Error())
	}
	return &types.Response{
		StatusCode: 200,
		Message: "healthy",
		Headers: nil,
		Body: nil,
	}
}