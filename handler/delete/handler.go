package delete

import (
	"strconv"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"

	"github.com/aed86/amboss-graph-api/model"
	"github.com/aed86/amboss-graph-api/response"
	node_service "github.com/aed86/amboss-graph-api/service/node"
	relation_service "github.com/aed86/amboss-graph-api/service/relation"
)

type Handler struct {
	nodeService *node_service.Service
	relationService *relation_service.Service
}

func NewHandler(nodeService *node_service.Service, relationService *relation_service.Service) Handler {
	return Handler{
		nodeService: nodeService,
		relationService: relationService,
	}
}

func (h *Handler) DeleteNode(params martini.Params, r render.Render) {
	if id, ok := params["id"]; ok {
		nodeId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.Error(r, err.Error(), 400)
			return
		}

		err = h.nodeService.DeleteNode(nodeId)
		if err != nil {
			response.Error(r, err.Error(), 400)
			return
		}

		response.Result(r, "Node is removed")
		return
	}

	response.Error(r, "NodeID must be provided", 400)
}

func (h *Handler) DeleteRelation(r render.Render, link model.Link) {
	err := h.relationService.DeleteRelation(link.Source, link.Target)
	if err != nil {
		response.Error(r, err.Error(), 400)
		return
	}

	response.Result(r, "Relation is removed")
}