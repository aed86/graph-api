package add

import (
	"github.com/martini-contrib/render"

	"github.com/aed86/amboss-graph-api/model"
	"github.com/aed86/amboss-graph-api/response"
	node_service "github.com/aed86/amboss-graph-api/service/node"
	relation_service "github.com/aed86/amboss-graph-api/service/relation"
)

type Handler struct {
	nodeService *node_service.Service
	linkService *relation_service.Service
}

func NewHandler(nodeService *node_service.Service, linkService *relation_service.Service) Handler {
	return Handler{
		nodeService: nodeService,
		linkService: linkService,
	}
}

func (h Handler) AddNode(node model.Node, r render.Render) {
	result, err := h.nodeService.AddNode(node)
	if err != nil {
		response.Error(r, err.Error(), 400)
		return
	}

	response.Result(r, result)
}

func (h Handler) AddRelation(link model.Link, r render.Render) {
	result, err := h.linkService.AddRelation(link)
	if err != nil {
		response.Error(r, err.Error(), 400)
		return
	}

	response.Result(r, result)
}
