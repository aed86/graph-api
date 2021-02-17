package get

import (
	"strconv"

	"github.com/go-martini/martini"

	"github.com/martini-contrib/render"

	"github.com/aed86/amboss-graph-api/model"
	"github.com/aed86/amboss-graph-api/response"
	"github.com/aed86/amboss-graph-api/service/node"
	"github.com/aed86/amboss-graph-api/service/relation"
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

func (h Handler) GetAllNodes(r render.Render) {
	result, err := h.ns.GetAllNodes(10)

	if err != nil {
		response.Error(r, err.Error(), 400)
		return
	}

	response.Result(r, result)
}

func (h Handler) GetNodeById(params martini.Params, r render.Render) {
	if v, ok := params["id"]; ok {
		nodeId, _ := strconv.ParseInt(v, 10, 64)
		n, err := h.ns.GetNodeById(nodeId)
		if err != nil {
			response.Error(r, err.Error(), 200)
			return
		}

		response.Result(r, n)
		return
	}

	response.Error(r, "Not found", 404)
}

func (h Handler) GetNeighbours(params martini.Params, r render.Render) {
	if v, ok := params["id"]; ok {
		nodeId, _ := strconv.ParseInt(v, 10, 64)
		n, err := h.ns.GetNeighboursForNodeById(nodeId)
		if err != nil {
			response.Error(r, err.Error(), 200)
			return
		}

		response.Result(r, n)
		return
	}

	response.Error(r, "Not found", 404)
}

func (h Handler) GetAll(r render.Render) {
	result, err := h.rs.GetAll()

	if err != nil {
		response.Error(r, err.Error(), 404)
		return
	}

	response.Result(r, result)
}
