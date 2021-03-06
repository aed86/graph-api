package model

import "github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"

type Link struct {
	Source int64 `json:"source"`
	Target int64 `json:"target"`
	Cost   int64 `json:"cost"`
}

func ParseFromDbTypeToLink(relationship dbtype.Relationship) Link {
	modelNode := Link{
		Source: relationship.StartId,
		Target: relationship.EndId,
	}

	if v, ok := relationship.Props["cost"]; ok {
		modelNode.Cost = v.(int64)
	}

	return modelNode
}
