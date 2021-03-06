package update

import (
	"github.com/martini-contrib/render"

	"github.com/aed86/amboss-graph-api/model"
	"github.com/aed86/amboss-graph-api/response"
	node_service "github.com/aed86/amboss-graph-api/service/node"
)

type Handler struct {
	ns *node_service.Service
}

func NewHandler(ns *node_service.Service) Handler {
	return Handler{
		ns: ns,
	}
}

func (h *Handler) UpdateNode(node model.Node, r render.Render) {
	res, err := h.ns.UpdateNode(node)
	if err != nil {
		response.Error(r, err.Error(), 404)
		return
	}

	response.Result(r, res)
}