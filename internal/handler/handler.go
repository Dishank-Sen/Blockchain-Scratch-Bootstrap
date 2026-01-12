package handler

import (
	"github.com/Dishank-Sen/quicnode/node"
	"github.com/Dishank-Sen/quicnode/types"
)

type Handler struct{
	node *node.Node
}

func NewHandler(n *node.Node) *Handler{
	return &Handler{
		node: n,
	}
}

func (h *Handler) handleErrorRes() *types.Response{
	errRes := &types.Response{
		StatusCode: 500,
		Message: "Error",
		Headers: map[string]string{},
		Body: []byte("Internal Server Error"),
	}
	return errRes
}