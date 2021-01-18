package main

import (
	"github.com/codegangsta/martini-contrib/binding"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"

	"github.com/aed86/amboss-graph-api/db"
	"github.com/aed86/amboss-graph-api/handler/add"
	delete2 "github.com/aed86/amboss-graph-api/handler/delete"
	"github.com/aed86/amboss-graph-api/handler/get"
	"github.com/aed86/amboss-graph-api/handler/update"
	"github.com/aed86/amboss-graph-api/model"
	node_service "github.com/aed86/amboss-graph-api/service/node"
	relation_service "github.com/aed86/amboss-graph-api/service/relation"
)

func main() {
	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
		Charset: "UTF-8",
	}))

	dbConnection := db.Connect()
	defer dbConnection.Disconnect()

	nodeService := node_service.New(&dbConnection)
	relationService := relation_service.New(&dbConnection)
	getHandler := get.NewHandler(nodeService, relationService)
	addHandler := add.NewHandler(nodeService, relationService)
	deleteHandler := delete2.NewHandler(nodeService, relationService)
	updateHandler := update.NewHandler(nodeService)

	m.Get("/", getHandler.GetAll)

	m.Group("/node", func (r martini.Router) {
		r.Get("", getHandler.GetAllNodes)
		r.Get("/:id", getHandler.GetNodeById)
		r.Post("", binding.Bind(model.Node{}), addHandler.AddNode)
		r.Put("/:id", binding.Bind(model.Node{}), updateHandler.UpdateNode)
		r.Delete("/:id", deleteHandler.DeleteNode)
	})

	m.Group("/relation", func (r martini.Router) {
		r.Post("", binding.Bind(model.Link{}), addHandler.AddRelation)
		r.Delete("", binding.Bind(model.Link{}), deleteHandler.DeleteRelation)
	})

	m.Run()
}
