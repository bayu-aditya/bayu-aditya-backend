package main

import (
	"context"
	reponats "github.com/bayu-aditya/bayu-aditya-backend/lib/core/repository/nats"
	"github.com/sirupsen/logrus"
	"net/http"
)

func shutdownHttp(server *http.Server) {
	ctx := context.Background()

	logrus.Info("Shutdown http: start")
	if err := server.Shutdown(ctx); err != nil {
		logrus.Errorf("Shutdown http: error %v", err)
	}
	logrus.Info("Shutdown http: success")
}

func shutdownNats(repoNats reponats.INatsRepository) {
	logrus.Info("Shutdown nats: start")
	repoNats.Close()
	logrus.Info("Shutdown nats: success")
}
