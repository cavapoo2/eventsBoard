package dblayer

import (
	"log"

	"github.com/cavapoo2/eventsBoard/lib/persistence"
	"github.com/cavapoo2/eventsBoard/lib/persistence/mongolayer"
)

type DBTYPE string

const (
	MONGODB    DBTYPE = "mongodb"
	DOCUMENTDB DBTYPE = "documentdb"
	DYNAMODB   DBTYPE = "dynamodb"
)

func NewPersistenceLayer(options DBTYPE, connection string) (persistence.DatabaseHandler, error) {
	log.Println("NewPersistenceLayer")
	switch options {
	case MONGODB:
		return mongolayer.NewMongoDBLayer(connection)
	}
	return nil, nil
}
