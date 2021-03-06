package get

import (
	"github.com/go-martini/martini"

	"github.com/martini-contrib/render"

	"github.com/aed86/amboss-graph-api/model"
	"github.com/aed86/amboss-graph-api/response"
	"github.com/aed86/amboss-graph-api/service/node"
	"github.com/aed86/amboss-graph-api/service/relation"
	"github.com/aed86/amboss-graph-api/utils"
)

type Handler struct {
	ns *node.Service
	rs *relation.Service
}

func NewHandler(ns *node.Service, rs *relation.Service) Handler {
	return Handler{
		ns: ns,
		rs: rs,
	}
}

type Result struct {
	Nodes []model.Node `json:"nodes"`
	Links []model.Link `json:"links"`
}

func (h Handler) GetAllNodes(payload model.ReqIn, r render.Render) {
	result, err := h.ns.GetAllNodes(utils.GetLimitFromReq(payload))

	if err != nil {
		response.Error(r, err.Error(), 400)
		return
	}

	response.Result(r, result)
}

func (h Handler) GetNodeById(params martini.Params, r render.Render) {
	nodeId, err := utils.GetId(params)
	if err != nil {
		response.Error(r, "validation error", 400)
	}

	n, err := h.ns.GetNodeById(nodeId)
	if err != nil {
		response.Error(r, err.Error(), 200)
		return
	}

	if n == nil {
		response.Error(r, "node not found", 400)
	}

	response.Result(r, n)
	return
}

func (h Handler) GetNeighbours(params martini.Params, payload model.ReqIn, r render.Render) {
	nodeId, err := utils.GetId(params)
	if err != nil {
		response.Error(r, "validation error", 400)
	}

	n, err := h.ns.GetNeighboursForNodeById(nodeId, utils.GetLimitFromReq(payload))
	if err != nil {
		response.Error(r, err.Error(), 200)
		return
	}

	if n == nil {
		response.Error(r, "node not found", 404)
	}

	response.Result(r, n)
	return
}

func (h Handler) GetAll(payload model.ReqIn, r render.Render) {
	result, err := h.rs.GetAll(utils.GetLimitFromReq(payload))

	if err != nil {
		response.Error(r, err.Error(), 404)
		return
	}

	response.Result(r, result)
}
