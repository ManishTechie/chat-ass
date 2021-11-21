// Package config  will be basis of environment configs in future
package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// SetEnvFromDotEnv sets the environment defined in dotenv
func SetEnvFromDotEnv() {
	err := godotenv.Load()
	if err != nil {
		logrus.WithError(err).Fatal(`Error loading .env file. Copy '.env.sample' to '.env' and add your environment data.
		DO NOT, I repeat DO NOT DELETE .env.sample file. Its part of version control.
		That being said here is the error`)
	}
	logrus.Info("Using .env environment")
}

func ensureEnvironment() {
	missingKeys := []string{}
	mapEnv, err := godotenv.Read(".env.sample")
	if err != nil {
		logrus.Error("Env sample: ", err)
		return
	}

	for key := range mapEnv {
		fmt.Println("key", key)
		val, ok := os.LookupEnv(key)
		if !ok {
			logrus.Errorf("Env variable %v not found", key)
			missingKeys = append(missingKeys, key)
		}

		logrus.Infof("%v : %q", key, val)
	}
	if len(missingKeys) > 0 {
		logrus.Panicf("Proper environment not set for: %v", missingKeys)
	}
}
