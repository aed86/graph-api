package fixture

import (
	"fmt"
	"io/ioutil"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"

	"github.com/aed86/amboss-graph-api/db"
)

type Fixture struct {
	dbConn db.Db
}

func New(dbConn db.Db) *Fixture {
	return &Fixture{
		dbConn,
	}
}

func init() {
	dbConnection := db.Connect()
	defer dbConnection.Disconnect()
	fixtureRunner := Fixture{dbConnection}
	err := fixtureRunner.Apply()
	if err != nil {
		panic(err.Error())
	}
}

func (f *Fixture) getFixtures() [][]byte {
	return [][]byte{
		f.ReadDump("./fixture/drop_data.dump"),
		f.ReadDump("./fixture/drop_graph.dump"),
		f.ReadDump("./fixture/create_nodes_and_links.dump"),
		f.ReadDump("./fixture/create_graph.dump"),
	}
}

func (f *Fixture) Apply() error {
	var bookmark1, bookmark2 string
	var err error

	fixtures := f.getFixtures()

	if bookmark1, err = f.Drop(); err != nil {
		return err
	}

	if bookmark2, err = f.ApplyDump(fixtures); err != nil {
		return err
	}

	session := f.dbConn.InitReadSession([]string{bookmark1, bookmark2})
	defer session.Close()

	if _, err = session.ReadTransaction(f.printResultTxFunc()); err != nil {
		return err
	}
	return nil
}

func (f *Fixture) Drop() (string, error) {
	session := f.dbConn.InitWriteSession()
	defer session.Close()

	if _, err := session.WriteTransaction(f.dropAllTxFunc()); err != nil {
		return "", err
	}

	return session.LastBookmark(), nil
}

func (f *Fixture) ReadDump(path string) []byte {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err.Error())
	}

	return file
}

func (f *Fixture) ApplyDump(dumps [][]byte) (string, error) {
	session := f.dbConn.InitWriteSession()
	defer session.Close()

	for _, dump := range dumps {
		if _, err := session.WriteTransaction(f.applyDumpTxFunc(dump)); err != nil {
			return "", err
		}
	}

	return session.LastBookmark(), nil
}

func (f *Fixture) printResultTxFunc() neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("MATCH (a)-[:ROAD]->(b) RETURN a.name, b.name", nil)
		if err != nil {
			return nil, err
		}

		for result.Next() {
			fmt.Printf("%s connected with %s\n", result.Record().Values[0], result.Record().Values[1])
		}

		return result.Consume()
	}
}

func (f *Fixture) dropAllTxFunc() neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		return tx.Run(`MATCH (n) DETACH DELETE n`, map[string]interface{}{})
	}
}

func (f *Fixture) applyDumpTxFunc(dump []byte) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		return tx.Run(string(dump), map[string]interface{}{})
	}
}
