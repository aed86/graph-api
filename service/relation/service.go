package relation

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"

	db_connection "github.com/aed86/amboss-graph-api/db"
	"github.com/aed86/amboss-graph-api/model"
)

type Service struct {
	db *db_connection.Db
}

func New(db *db_connection.Db) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) GetShortestPath(cond model.PathIn) ([]model.PathOut, error) {
	session := s.db.InitReadSession()
	defer session.Close()

	res, err := session.ReadTransaction(s.findShortestPath(cond))
	if err != nil {
		return nil, err
	}

	result := s.buildPathInfoFromRecords(res.([]db.Record))

	return result, nil
}

func (s *Service) buildPathInfoFromRecords(records []db.Record) []model.PathOut {
	var paths []model.PathOut

	for _, record := range records {
		values := record.Values

		path := model.PathOut{
			Idx: values[0].(int64),
			SourceNodeName: values[1].(string),
			TargetNodeName: values[2].(string),
			TotalCost: values[3].(float64),
			Path: values[4],
			PathCosts: values[5],
		}

		paths = append(paths, path)
	}

	return paths
}

func (s Service) AddRelation(personID1, personID2 int64) (*model.Relation, error) {
	session := s.db.InitWriteSession()
	defer session.Close()

	res, err := session.WriteTransaction(s.addRelation(personID1, personID2))
	if err != nil {
		return nil, err
	}

	relation := s.buildRelationsFromRecordsPair(res.([]db.Record))

	return &relation, nil
}

func (s Service) DeleteRelation(personID1, personID2 int64) error {
	session := s.db.InitWriteSession()
	defer session.Close()

	_, err := session.WriteTransaction(s.deleteRelation(personID1, personID2))
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetAll() (*model.Relation, error) {
	session := s.db.InitReadSession()
	defer session.Close()

	res, err := session.ReadTransaction(s.findAllNodesWithRelationsTxFunc())
	if err != nil {
		return nil, err
	}

	relation := s.buildRelationsFromRecordsPair(res.([]db.Record))

	return &relation, nil
}

func (s *Service) buildRelationsFromRecordsPair(res []db.Record) model.Relation {
	nodes := make(map[int64]model.Node, 0)
	var links []model.Link
	for _, recordPair := range res {
		sourceNode := model.ParseFromDbTypeToNode(recordPair.Values[0].(dbtype.Node))
		targetNode := model.ParseFromDbTypeToNode(recordPair.Values[1].(dbtype.Node))
		if _, ok := nodes[sourceNode.ID]; !ok {
			nodes[sourceNode.ID] = sourceNode
		}
		if _, ok := nodes[targetNode.ID]; !ok {
			nodes[targetNode.ID] = targetNode
		}

		links = append(links, model.Link{
			Source: sourceNode.ID,
			Target: targetNode.ID,
		})
	}

	return model.Relation{
		Nodes: nodes,
		Links: links,
	}
}

func (s *Service) findAllNodesWithRelationsTxFunc() neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("MATCH (a)-[:ROAD]->(b) RETURN a, b", nil)
		if err != nil {
			return nil, err
		}

		var records []db.Record
		for result.Next() {
			if result.Record() != nil {
				records = append(records, *result.Record())
			}
		}

		return records, nil
	}
}

func (s *Service) findShortestPath(cond model.PathIn) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		query := `MATCH (source:Node {name: $sname}), (target:Node {name: $tname}) CALL gds.beta.shortestPath.dijkstra.stream('myGraph4', {sourceNode: id(source),targetNode: id(target),relationshipWeightProperty: 'cost'}) YIELD index, sourceNode, targetNode, totalCost, nodeIds, costs RETURN index, gds.util.asNode(sourceNode).name AS sourceNodeName, gds.util.asNode(targetNode).name AS targetNodeName, totalCost, [nodeId IN nodeIds |gds.util.asNode(nodeId).name] AS nodeNames, costs ORDER BY index`
		result, err := tx.Run(query, map[string]interface{}{
			"sname": cond.SourceName,
			"tname": cond.TargetName,
		})

		if err != nil {
			return nil, err
		}
		var records []db.Record
		for result.Next() {
			if result.Record() != nil {
				records = append(records, *result.Record())
			}
		}

		return records, nil
	}
}

func (s Service) addRelation(personID1 int64, personID2 int64) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(
			"MATCH (a:Node) WHERE ID(a) = $ID1 "+
				"MATCH (b:Node) WHERE ID(b) = $ID2 "+
				"MERGE (a)-[:ROAD]->(b) RETURN a, b", map[string]interface{}{"ID1": personID1, "ID2": personID2})

		if err != nil {
			return nil, err
		}

		var records []db.Record
		for result.Next() {
			if result.Record() != nil {
				records = append(records, *result.Record())
			}
		}

		return records, nil
	}
}

func (s Service) deleteRelation(personID1 int64, personID2 int64) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(
			"MATCH (a:Node)-[d:ROAD]->(b:Node) WHERE ID(a) = $ID1 and ID(b) = $ID2 DELETE d",
			map[string]interface{}{"ID1": personID1, "ID2": personID2},
		)

		if err != nil {
			return nil, err
		}

		return result.Consume()
	}
}
