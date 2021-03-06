package db

import (
	"fmt"
	"log"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Db struct {
	Driver neo4j.Driver
}

//Neo4jConfiguration holds the configuration for connecting to the DB
type Neo4jConfiguration struct {
	URL      string
	Username string
	Password string
	Database string
}

//newDrive is a method for Neo4jConfiguration to return a connection to the DB
func (nc *Neo4jConfiguration) newDriver() (neo4j.Driver, error) {
	return neo4j.NewDriver(nc.URL, neo4j.BasicAuth(nc.Username, nc.Password, ""))
}

func Connect() Db {
	configuration := parseConfiguration()
	driver, err := configuration.newDriver()
	if err != nil {
		log.Fatal(err)
	}

	return Db{
		Driver: driver,
	}
}

func (db Db) InitReadSession(bookmarks []string) neo4j.Session {
	session := db.Driver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeRead,
		DatabaseName: "neo4j",
		Bookmarks: bookmarks,
	})

	return session
}

func (db Db) InitWriteSession() neo4j.Session {
	session := db.Driver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeWrite,
		DatabaseName: "neo4j",
	})

	return session
}

func parseConfiguration() *Neo4jConfiguration {
	return &Neo4jConfiguration{
		URL:      "neo4j://localhost:7687",
		Username: "neo4j",
		Password: "testing",
	}
}

func (db Db) Disconnect() {
	if err := db.Driver.Close(); err != nil {
		log.Fatal(fmt.Errorf("could not close resource: %w", err))
	}
}
