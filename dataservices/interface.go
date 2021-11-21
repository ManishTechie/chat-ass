package dataservices

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type DBClient struct {
	DB *mongo.Client
}

type BackendServiceDBInterface interface {

	//ms-sql
	Connect(connectionString string) (setupError error)
	Close() error
}
