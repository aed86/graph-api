package model

import "github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"

//Node ...
type Node struct {
	ID   int64  `json:"Id"`
	Name string `json:"name"`
	Born string  `json:"born"`
}

func ParseFromDbTypeToNode(node dbtype.Node) Node {
	modelNode := Node{
		ID: node.Id,
	}
	props := node.Props

	if v, ok := props["name"]; ok {
		modelNode.Name = v.(string)
	}

	if v, ok := props["born"]; ok {
		modelNode.Born = v.(string)
	}

	return modelNode
}