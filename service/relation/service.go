package relation

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"

	db_connection "github.com/aed86/amboss-graph-api/db"
	"github.com/aed86/amboss-graph-api/model"
	"github.com/aed86/amboss-graph-api/service/node"
)

type Service struct {
	db *db_connection.Db
	ns *node.Service
}

func New(db *db_connection.Db, ns *node.Service) *Service {
	return &Service{
		db: db,
		ns: ns,
	}
}

func (s *Service) GetShortestPath(cond model.PathIn) ([]model.PathOut, error) {
	session := s.db.InitReadSession([]string{})
	defer session.Close()

	res, err := session.ReadTransaction(s.findShortestPath(cond))
	if err != nil {
		return nil, err
	}

	result := s.buildPathInfoFromRecords(res.([]db.Record))

	return result, nil
}

func (s *Service) AddRelation(nodeID1, nodeID2 int64) (*model.Relation, error) {
	bookmark1, err := s.ns.CheckIfNodesExist([]int64{nodeID1, nodeID2})

	session := s.db.InitWriteSession([]string{bookmark1})
	defer session.Close()

	res, err := session.WriteTransaction(s.addRelation(nodeID1, nodeID2))
	if err != nil {
		return nil, err
	}

	relation := s.buildRelationsFromRecordsPair(res.([]db.Record))

	return &relation, nil
}

func (s *Service) DeleteRelation(nodeID1, nodeID2 int64) error {
	bookmark1, err := s.ns.CheckIfNodesExist([]int64{nodeID1, nodeID2})

	session := s.db.InitWriteSession([]string{bookmark1})
	defer session.Close()

	_, err = session.WriteTransaction(s.deleteRelation(nodeID1, nodeID2))
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetAll(limit int64) (*model.Relation, error) {
	session := s.db.InitReadSession([]string{})
	defer session.Close()

	res, err := session.ReadTransaction(s.findAllNodesWithRelationsTxFunc(limit))
	if err != nil {
		return nil, err
	}

	relation := s.buildRelationsForAllCase(res.([]db.Record))

	return &relation, nil
}

func (s *Service) buildRelationsForAllCase(res []db.Record) model.Relation {
	nodes := make(map[int64]model.Node, 0)
	links := make(map[int64]model.Link, 0)
	var nodeRecords []dbtype.Node
	var linkRecords []dbtype.Relationship

	for _, record := range res {
		nodeRecords = append(nodeRecords, record.Values[0].(dbtype.Node))
		for _, lr := range record.Values[1].([]interface{}) {
			linkRecords = append(linkRecords, lr.(dbtype.Relationship))
		}
	}

	if len(nodeRecords) > 0 {
		for _, nodeRecord := range nodeRecords {
			node := model.ParseFromDbTypeToNode(nodeRecord)
			nodes[node.ID] = node
		}
	}

	if len(linkRecords) > 0 {
		for _, linkRecord := range linkRecords {
			link := model.ParseFromDbTypeToLink(linkRecord)
			if _, ok := nodes[link.Target]; ok {
				if _, ok2 := nodes[link.Source]; ok2 {
					links[link.ID] = link
				}
			}
		}
	}

	return model.Relation{
		Nodes: nodes,
		Links: links,
	}
}

func (s *Service) buildRelationsFromRecordsPair(res []db.Record) model.Relation {
	nodes := make(map[int64]model.Node, 0)
	links := make(map[int64]model.Link, 0)
	for _, recordPair := range res {
		sourceNode := model.ParseFromDbTypeToNode(recordPair.Values[0].(dbtype.Node))
		targetNode := model.ParseFromDbTypeToNode(recordPair.Values[1].(dbtype.Node))
		link := model.ParseFromDbTypeToLink(recordPair.Values[2].(dbtype.Relationship))
		if _, ok := nodes[sourceNode.ID]; !ok {
			nodes[sourceNode.ID] = sourceNode
		}
		if _, ok := nodes[targetNode.ID]; !ok {
			nodes[targetNode.ID] = targetNode
		}

		links[link.ID] = link
	}

	return model.Relation{
		Nodes: nodes,
		Links: links,
	}
}

func (s *Service) findAllNodesWithRelationsTxFunc(limit int64) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(
			"MATCH (a:Node) MATCH (a)-[d:ROAD]-(:Node) RETURN a, collect(distinct d) as roads limit $limit",
			map[string]interface{}{
				"limit": limit,
			},
		)
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

func (s *Service) addRelation(nodeID1 int64, nodeID2 int64) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(
			"MATCH (a:Node) WHERE ID(a) = $ID1 "+
				"MATCH (b:Node) WHERE ID(b) = $ID2 "+
				"MERGE (a)-[d:ROAD]->(b) RETURN a, b, d", map[string]interface{}{"ID1": nodeID1, "ID2": nodeID2})

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

func (s *Service) deleteRelation(nodeID1 int64, nodeID2 int64) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(
			"MATCH (a:Node)-[d:ROAD]->(b:Node) WHERE ID(a) = $ID1 and ID(b) = $ID2 DELETE d",
			map[string]interface{}{"ID1": nodeID1, "ID2": nodeID2},
		)

		if err != nil {
			return nil, err
		}

		return result.Consume()
	}
}

func (s *Service) buildPathInfoFromRecords(records []db.Record) []model.PathOut {
	var paths []model.PathOut

	for _, record := range records {
		values := record.Values

		path := model.PathOut{
			Idx:            values[0].(int64),
			SourceNodeName: values[1].(string),
			TargetNodeName: values[2].(string),
			TotalCost:      values[3].(float64),
			Path:           values[4],
			PathCosts:      values[5],
		}

		paths = append(paths, path)
	}

	return paths
}
