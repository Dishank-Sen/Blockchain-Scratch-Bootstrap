package router

import (
	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/internal/handler"
	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/utils/logger"
	"github.com/Dishank-Sen/quicnode/node"
)

type Router struct{
	node *node.Node
}

func NewRouter(n *node.Node) *Router{
	return &Router{
		node: n,
	}
}

func (r *Router) HandleRoutes(){
	n := r.node
	h, err := handler.NewHandler(n)
	if err != nil{
		logger.Error(err.Error())
	}
	n.Handle("connect", h.Connect)
	n.Handle("peers", h.Peers)
	n.Handle("heartbeat", h.Heartbeat)
}