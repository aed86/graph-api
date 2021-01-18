package config

import (
	"log"
	"os"
	"github.com/go-pg/pg/v9"
)
// Connecting to db
func Connect() *pg.DB {
	opts := &pg.Options{
		User: "db_username",
		Password: "db_password",
		Addr: "localhost:5432",
		Database: "db_dbname",
	}
	var db *pg.DB = pg.Connect(opts)
	if db == nil {
		log.Printf("Failed to connect")
		os.Exit(100)
	}
	log.Printf("Connected to db")
	return db
}