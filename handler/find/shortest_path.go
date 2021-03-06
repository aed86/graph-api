package find

import (
	"github.com/martini-contrib/render"

	"github.com/aed86/amboss-graph-api/model"
	"github.com/aed86/amboss-graph-api/response"
	relation_service "github.com/aed86/amboss-graph-api/service/relation"
)

type Handler struct {
	rs *relation_service.Service
}

func NewHandler(relationService *relation_service.Service) Handler {
	return Handler{
		rs: relationService,
	}
}

func (h *Handler) ShortestPath(cond model.PathIn, r render.Render) {
	result, err := h.rs.GetShortestPath(cond)
	if err != nil {
		response.Error(r, err.Error(), 400)
		return
	}

	response.Result(r, result)
}
