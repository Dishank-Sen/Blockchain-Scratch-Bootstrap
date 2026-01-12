package router

import (
	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/internal/handler"
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
	h := handler.NewHandler(n)
	n.Handle("connect", h.Connect)
	n.Handle("peers", h.Peers)
}