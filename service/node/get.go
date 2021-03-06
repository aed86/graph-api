package node

import (
	"errors"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"

	"github.com/aed86/amboss-graph-api/model"
)

func (s Service) GetNodeById(nodeId int64) (*model.Node, error) {
	session := s.db.InitReadSession([]string{})
	defer session.Close()

	node, err := session.ReadTransaction(s.findNodeByIdTxFunc(nodeId))
	if err != nil {
		return nil, err
	}

	result := model.ParseFromDbTypeToNode(node.(dbtype.Node))
	return &result, nil
}

func (s Service) GetNeighboursForNodeById(baseNodeID, limit int64) (*[]model.Node, error) {
	session := s.db.InitReadSession([]string{})
	defer session.Close()

	nodes, err := session.ReadTransaction(s.findNeighboursByNodeIdTxFunc(baseNodeID, limit))
	if err != nil {
		return nil, err
	}

	var modelResult []model.Node
	for _, node := range nodes.([]interface{}) {
		modelResult = append(modelResult, model.ParseFromDbTypeToNode(node.(dbtype.Node)))
	}

	return &modelResult, nil
}

func (s Service) GetAllNodes(limit int64) (*[]model.Node, error) {
	session := s.db.InitWriteSession()
	defer session.Close()

	nodes, err := session.ReadTransaction(s.findAllNodesTxFunc(limit))
	if err != nil {
		return nil, err
	}

	var modelResult []model.Node
	for _, node := range nodes.([]interface{}) {
		modelResult = append(modelResult, model.ParseFromDbTypeToNode(node.(dbtype.Node)))
	}

	return &modelResult, nil
}

func (s Service) parseNode(record *neo4j.Record) model.Node {
	ID, _ := record.Get("ID")
	name, _ := record.Get("name")
	born, _ := record.Get("born")

	return model.Node{
		ID:   ID.(int64),
		Name: name.(string),
		Born: born.(string),
	}
}

func (s Service) findNodeByIdTxFunc(id int64) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(
			"MATCH (n:Node) WHERE ID(n) = $ID RETURN n",
			map[string]interface{}{
				"ID": id,
			},
		)
		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().Values[0], nil
		}

		return nil, errors.New("node is not found")
	}
}

func (s Service) findAllNodesTxFunc(limit int64) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		rtx, err := tx.Run(
			"MATCH (e:Node) RETURN e LIMIT $limit",
			map[string]interface{}{
				"limit": limit,
			},
		)
		if err != nil {
			return nil, err
		}

		var results []interface{}
		for rtx.Next() {
			results = append(results, rtx.Record().Values[0])
		}

		return results, nil
	}
}

func (s Service) findNeighboursByNodeIdTxFunc(baseNodeID, limit int64) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		rtx, err := tx.Run(
			"MATCH (a)-[:ROAD]->(b) WHERE ID(a) = $ID RETURN b LIMIT $limit",
			map[string]interface{}{
				"ID": baseNodeID,
				"limit": limit,
			},
		)
		if err != nil {
			return nil, err
		}

		var results []interface{}
		for rtx.Next() {
			results = append(results, rtx.Record().Values[0])
		}

		return results, nil
	}
}
