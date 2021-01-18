package node

import (
	"github.com/aed86/amboss-graph-api/db"
)

type Service struct {
	db *db.Db
}

func New(db *db.Db) *Service {
	return &Service{
		db: db,
	}
}