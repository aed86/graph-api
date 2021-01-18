package node

import (
	"errors"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"

	"github.com/aed86/amboss-graph-api/model"
)

func (s *Service) AddNode(node model.Node) (*model.Node, error) {
	session := s.db.InitWriteSession()
	defer session.Close()

	record, err := session.WriteTransaction(s.addNodeTxFunc(node))
	if err != nil {
		return nil, err
	}

	resultNode := model.ParseFromDbTypeToNode(record.(dbtype.Node))
	return &resultNode, nil
}

func (s *Service) addNodeTxFunc(node model.Node) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(
			"CREATE (a:Node {name: $name, born: $born}) RETURN a, ID(a)",
			map[string]interface{}{
				"name": node.Name,
				"born": node.Born,
			},
		)

		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().Values[0], nil
		}

		return result.Consume()
	}
}

func (s *Service) matchNodeByNameTxFunc(name string) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(
			"MATCH (a:Node {name: $name}) RETURN id(a)",
			map[string]interface{}{
				"name": name,
			},
		)
		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().Values[0], nil
		}

		return nil, errors.New("one record was expected")
	}
}
