// Package server provides a thin wrapper around the default net/http Server with Gin Engine as its router.
package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

// Start the HTTP server in a separate Go routine, and listen for OS kill signals.
// Whenever a signal is received, shut it down gracefully.
func Start(router http.Handler, serverName string, cleanupFns ...func() error) {

	// create the HTTP server
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  time.Duration(5) * time.Second,
		WriteTimeout: time.Duration(5) * time.Second,
		IdleTimeout:  time.Duration(5) * time.Second,
	}

	// create the channel for the OS signal, used for shutdown
	stopCh := make(chan os.Signal, 1)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.WithError(err).Fatalf("HTTP Server exited with a fatal error")
		}
	}()

	// notify the stopCh channel about OS interrupt, SIGINT and SIGTERM signals
	signal.Notify(stopCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	// block until a signal is received
	<-stopCh
	logrus.Info("Shut down signal received")

	// create a context with timeout and cancel
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cleanupErrors := []string{}
		for _, cleanupFn := range cleanupFns {
			if err := cleanupFn(); err != nil {
				cleanupErrors = append(cleanupErrors, err.Error())
			}
		}
		if len(cleanupErrors) > 0 {
			logrus.WithError(errors.New(strings.Join(cleanupErrors, ", "))).Error("Failed to clean up properly after shutdown")
		}
		cancel()
	}()

	// initiate the graceful shutdown of the server within the previously created context
	if err := srv.Shutdown(ctx); err != nil {
		logrus.WithError(err).Fatal("Failed to gracefully shut down the Server")
	}
	logrus.Info("Server shut down gracefully")
}
