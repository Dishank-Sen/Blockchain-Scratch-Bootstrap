package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/utils/logger"
	"github.com/Dishank-Sen/quicnode/types"
)

type heartbeatPayload struct {
	ID string `json:"id"`
}

func (h *Handler) Heartbeat(ctx context.Context, req *types.Request) *types.Response {
	var hb heartbeatPayload

	if err := json.Unmarshal(req.Body, &hb); err != nil {
		logger.Error("invalid heartbeat payload")
		return &types.Response{
			StatusCode: 400,
			Message: "bad request",
		}
	}

	logger.Info(fmt.Sprintf("heartbeat from %s", hb.ID))

	if err := h.store.UpdateLastSeen(hb.ID); err != nil {
		logger.Error(err.Error())
	}

	return &types.Response{
		StatusCode: 200,
		Message: "healthy",
	}
}
