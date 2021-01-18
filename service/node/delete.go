package node

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func (s *Service) DeleteNode(nodeID int64) error {
	session := s.db.InitWriteSession()
	defer session.Close()

	_, err := session.WriteTransaction(s.deleteNodeTxFunc(nodeID))
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) deleteNodeTxFunc(nodeID int64) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(
			"MATCH (a:Node) WHERE ID(a) = $ID DETACH DELETE a RETURN a",
			map[string]interface{}{
				"ID": nodeID,
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
