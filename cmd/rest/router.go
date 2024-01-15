package main

import (
	"github.com/bayu-aditya/bayu-aditya-backend/lib/core/config"
	reponats "github.com/bayu-aditya/bayu-aditya-backend/lib/core/repository/nats"
	usecasemonopolybanking "github.com/bayu-aditya/bayu-aditya-backend/lib/core/usecase/monopolybanking"
	"github.com/bayu-aditya/bayu-aditya-backend/lib/delivery/rest"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func initializeRouter(r *gin.Engine) (stop func()) {
	cfg, err := config.New("./files/config.yaml")
	if err != nil {
		logrus.Fatalf("init config: %v", err)
	}

	repoNats, err := reponats.NewNatsRepository(cfg)
	if err != nil {
		logrus.Fatalf("init repo nats: %v", err)
	}

	usecaseMonopoly := usecasemonopolybanking.NewMonopolyBankingUsecase(repoNats)

	handler := rest.NewHttpHandler(usecaseMonopoly)

	r.GET("/", handler.HealthCheck)
	r.POST("/mbs/room", handler.MonopolyCreateRoom)
	r.POST("/mbs/room-join", handler.MonopolyJoinRoom)
	r.GET("/mbs/sse", handler.MonopolyServerSentEvent)

	stop = func() {
		shutdownNats(repoNats)
	}

	return
}
