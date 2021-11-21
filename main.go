// XXImo Mobility Management - Customer Settings Service

package main

import (
	"os"

	"github.com/backend-service/chat"
	"github.com/backend-service/config"
	"github.com/backend-service/constants"
	"github.com/backend-service/dataservices"
	"github.com/backend-service/engine"
	"github.com/backend-service/server"
	"github.com/sirupsen/logrus"
)

// VERSION will be overwritten by the Go toolchain during release
var VERSION string = "v0.0.1"

func main() {
	config.SetConfigs()
	// set up the database connection
	var db dataservices.BackendServiceDBInterface = &dataservices.DBClient{}
	if dbConnectErr := db.Connect(os.Getenv(constants.DB_CONNECTION_STRING)); dbConnectErr != nil {
		logrus.WithError(dbConnectErr).Fatal("Failed to set up dataservices")
	}
	db = dataservices.DB()
	chat.ChatServer()

	// start the server
	server.Start(
		engine.BuildGinEngine(db, VERSION),
		"backend-service",
		db.Close,
	)
}
