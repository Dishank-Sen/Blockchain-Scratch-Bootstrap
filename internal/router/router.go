package router

import (
	"fmt"

	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/internal/handler"
	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/internal/peers"
	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/utils/logger"
	"github.com/Dishank-Sen/quicnode/node"
	"github.com/Dishank-Sen/quicnode/types"
)

type Router struct{
	node *node.Node
	store *peers.Store
}

func NewRouter(n *node.Node) *Router{
	store, err := peers.GetStore()
	if err != nil{
		logger.Error(err.Error())
	}
	return &Router{
		node: n,
		store: store,
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

func (r *Router) HandleEvents(){
	for event := range r.node.Events() {
		switch event.Type {

		case types.EventConnOpened:
			logger.Info(fmt.Sprintf("connected: %s", event.ConnID))

		case types.EventConnClosed:
			logger.Info(fmt.Sprintf("disconnected: %s", event.ConnID))

			// cleanup mapping
			r.store.RemoveByConnID(event.ConnID)
		}
	}
}