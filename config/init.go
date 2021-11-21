package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Environment constants
const (
	EnvironmentKey = "ENV"
	ProductionEnv  = "production"
	StagingEnv     = "staging"
	DevEnv         = "dev"
)

// config handles the overall config of entire ums application
func SetConfigs() {

	// For compatibility reasons setting env to localhost
	if os.Getenv("ENV") == "" {
		if err := os.Setenv("ENV", "localhost"); err != nil {
			logrus.Error("Error occured setting ENV to localhost")
		}
		logrus.Info("For compatibility reasons ENV set to localhost")
	}
	SetEnvFromDotEnv()
	ensureEnvironment()
}
