package server

import (
	"context"

	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/types"
)

type HandlerFunc func(ctx context.Context, req *types.Request) (*types.Response, error)

type routeKey struct {
	method string
	path   string
}

type Router struct {
	routes map[routeKey]HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[routeKey]HandlerFunc),
	}
}

func (r *Router) Get(path string, h HandlerFunc) {
	r.routes[routeKey{"GET", path}] = h
}

func (r *Router) Post(path string, h HandlerFunc) {
	r.routes[routeKey{"POST", path}] = h
}

func (r *Router) Dispatch(ctx context.Context, req *types.Request) *types.Response {
	h, ok := r.routes[routeKey{req.Method, req.Path}]
	if !ok {
		return &types.Response{
			StatusCode: 404,
			Message:    "Not Found",
			Body:       []byte("route not found"),
		}
	}

	resp, err := h(ctx, req)
	if err != nil {
		return &types.Response{
			StatusCode: 500,
			Message:    "Internal Error",
			Body:       []byte(err.Error()),
		}
	}
	return resp
}
