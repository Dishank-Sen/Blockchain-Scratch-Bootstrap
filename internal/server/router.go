package server

import (
	"context"

	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/types"
	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/utils/logger"
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
	// logger.Debug(path)
	r.routes[routeKey{"GET", path}] = h
}

func (r *Router) Post(path string, h HandlerFunc) {
	r.routes[routeKey{"POST", path}] = h
}

func (r *Router) Dispatch(ctx context.Context, req *types.Request) *types.Response {
	// logger.Debug(req.Method)
	// logger.Debug(req.Path)
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
		logger.Error(err.Error())
		return resp
	}
	return resp
}
