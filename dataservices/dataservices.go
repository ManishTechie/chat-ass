package dataservices

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	// appDB is the application's database connection
	appDB *DBClient
	// setupOnce ensures that the connection can be setup only once
	setupOnce sync.Once
	client    *mongo.Client
)

// Connect sets up the global database connection with sensible defaults.
func (ms *DBClient) Connect(connectionString string) (setupError error) {
	setupOnce.Do(func() {

		clientOptions := options.Client().
			ApplyURI(connectionString)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		client, setupError = mongo.Connect(ctx, clientOptions)
		if setupError != nil {
			log.Fatal(setupError)
		}
		setupError = client.Ping(ctx, readpref.Primary())
		if setupError != nil {
			log.Fatal("ping error", setupError)
		}
		log.Info("DB connected successfully")
		appDB = &DBClient{DB: client}

	})
	return
}

// Close the connection to the database.
func (ms *DBClient) Close() error {
	err := ms.DB.Disconnect(context.Background())
	if err != nil {
		return err
	}
	return nil
}

// DB returns the global database connection.
func DB() *DBClient {
	return appDB
}
