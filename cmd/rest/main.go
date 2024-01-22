package main

import (
	"errors"
	"github.com/bayu-aditya/bayu-aditya-backend/lib/core/util"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	router, stopRouter := initializeRouter()

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	finishChan := util.GracefulShutdown(func() {
		shutdownHttp(server)
		stopRouter()
	})

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logrus.Fatalf("Run Server: %v", err)
	}

	<-finishChan
	logrus.Info("Shutdown application")
}
