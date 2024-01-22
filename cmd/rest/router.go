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

	// initialize repository
	repoNats, err := reponats.NewNatsRepository(cfg)
	if err != nil {
		logrus.Fatalf("init repo nats: %v", err)
	}

	// initialize usecase
	usecaseMonopoly := usecasemonopolybanking.NewMonopolyBankingUsecase(repoNats)

	// initialize http handler
	handler := rest.NewHttpHandler(usecaseMonopoly)

	// initialize endpoint
	r.Use(rest.CORS())
	r.GET("/", handler.HealthCheck)
	r.POST("/mbs/room-join", handler.MonopolyJoinRoom)
	r.POST("/mbs/room", handler.MonopolyCreateRoom)
	r.GET("/mbs/room/:room_id/state", handler.MonopolyGetState)
	r.POST("/mbs/room/:room_id/transaction", handler.MonopolyCreateTransaction)
	r.GET("/mbs/room/:room_id/sse", handler.MonopolyServerSentEvent)

	// graceful shutdown
	stop = func() {
		shutdownNats(repoNats)
	}

	return
}
