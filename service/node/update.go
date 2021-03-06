package node

import (
	"errors"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"

	"github.com/aed86/amboss-graph-api/model"
)

func (s *Service) UpdateNode(node model.Node) (*model.Node, error) {
	session := s.db.InitWriteSession()
	defer session.Close()

	record, err := session.WriteTransaction(s.updateNodeTxFunc(node))
	if err != nil {
		return nil, err
	}

	if record == nil {
		return nil, errors.New("node is not found")
	}

	resultNode := model.ParseFromDbTypeToNode(record.(dbtype.Node))
	return &resultNode, nil
}

func (s *Service) updateNodeTxFunc(node model.Node) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(
			"MATCH (a:Node) WHERE ID(a) = $ID SET a.name = $name, a.born = $born RETURN a",
			map[string]interface{}{
				"ID":   node.ID,
				"name": node.Name,
				"born": node.Born,
			},
		)
		if err != nil {
			return nil, err
		}

		for result.Next() {
			return result.Record().Values[0], nil
		}

		return nil, nil
	}
}
