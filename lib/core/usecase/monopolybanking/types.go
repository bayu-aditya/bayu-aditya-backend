package usecasemonopolybanking

import (
	"context"
	modelmonopoly "github.com/bayu-aditya/bayu-aditya-backend/lib/core/model/monopoly"
	reponats "github.com/bayu-aditya/bayu-aditya-backend/lib/core/repository/nats"
)

type IMonopolyBankingUsecase interface {
	CreateRoom(ctx context.Context, player modelmonopoly.StatePlayer, initialBalance int64) (roomID string, roomPass string, err error)
	JoinRoom(ctx context.Context, playerID, playerName, roomID, roomPass string) error
	LeaveRoom(ctx context.Context, playerID, roomID, roomPass string) error
	GetState(ctx context.Context, playerID, roomID string) (modelmonopoly.State, error)
	CreateTransaction(ctx context.Context, roomID, playerID, targetPlayerID string, amount int64, mode string) error
	SubscribeState(ctx context.Context, roomID string) (stateChan <-chan modelmonopoly.State, stop func(), err error)
}

type usecase struct {
	prefix         string
	repositoryNats reponats.INatsRepository
}
