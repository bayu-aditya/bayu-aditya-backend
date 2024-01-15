package main

import (
	"errors"
	"github.com/bayu-aditya/bayu-aditya-backend/lib/core/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	r := gin.Default()
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	stopRouter := initializeRouter(r)

	waitFinishChan := util.GracefulShutdown(func() {
		shutdownHttp(server)
		stopRouter()
	})

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logrus.Fatalf("Run Server: %v", err)
	}

	<-waitFinishChan
	logrus.Info("Shutdown application")
}
