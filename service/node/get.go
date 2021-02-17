package node

import (
	"errors"
	"log"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"

	"github.com/aed86/amboss-graph-api/model"
)

func (s Service) GetNodeById(nodeId int64) (*model.Node, error) {
	session := s.db.InitReadSession()
	defer session.Close()

	node, err := session.ReadTransaction(s.findNodeByIdTxFunc(nodeId))
	if err != nil {
		return nil, err
	}

	result := model.ParseFromDbTypeToNode(node.(dbtype.Node))
	return &result, nil
}

func (s Service) GetNeighboursForNodeById(baseNodeID int64) (*[]model.Node, error) {
	session := s.db.InitReadSession()
	defer session.Close()

	nodes, err := session.ReadTransaction(s.findNeighboursByNodeIdTxFunc(baseNodeID))
	if err != nil {
		return nil, err
	}

	var modelResult []model.Node
	for _, node := range nodes.([]interface{}) {
		modelResult = append(modelResult, model.ParseFromDbTypeToNode(node.(dbtype.Node)))
	}

	return &modelResult, nil
}

func (s Service) GetAllNodes(limit int) (*[]model.Node, error) {
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

func (s Service) GetAll() *[]model.Node {
	session := s.db.InitWriteSession()
	defer session.Close()

	query := `MATCH (e:Node) RETURN e.ID as ID, e.name as name, e.born as born LIMIT $limit`
	result, err := session.Run(query, map[string]interface{}{"limit": 10})
	if err != nil {
		log.Println("Error querying Neo4j", err)
		return nil
	}
	var nodes []model.Node
	for result.Next() {
		record := result.Record()
		node := s.parseNode(record)
		nodes = append(nodes, node)
	}

	return &nodes
}

func (s Service) parseNode(record *neo4j.Record) model.Node {
	ID, _ := record.Get("ID")
	name, _ := record.Get("name")
	born, _ := record.Get("born")

	return model.Node{
		ID:   ID.(int64),
		Name: name.(string),
		Born: born.(int64),
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

		return nil, errors.New("one record was expected")
	}
}

func (s Service) findAllNodesTxFunc(limit int) neo4j.TransactionWork {
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

func (s Service) findNeighboursByNodeIdTxFunc(baseNodeID int64) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		rtx, err := tx.Run(
			"MATCH (a)-[:DIRECTED]->(b) WHERE ID(a) = $ID RETURN b",
			map[string]interface{}{
				"ID": baseNodeID,
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
