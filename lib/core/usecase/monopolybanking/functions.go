package usecasemonopolybanking

import (
	"context"
	"errors"
	"github.com/bayu-aditya/bayu-aditya-backend/lib/core/model/constant"
	modelmonopoly "github.com/bayu-aditya/bayu-aditya-backend/lib/core/model/monopoly"
	"github.com/bayu-aditya/bayu-aditya-backend/lib/core/util"
)

func (u *usecase) CreateRoom(ctx context.Context, player modelmonopoly.StatePlayer, initialBalance int64) (roomID string, roomPass string, err error) {
	prefix := u.prefix + ".CreateRoom"
	roomIDLength := 5
	passRoomLength := 5

	roomID = util.RandString(roomIDLength)
	roomPass = util.RandString(passRoomLength)

	// TODO check roomID is exist, if exist then recreate roomID

	player.Balance = initialBalance

	state := modelmonopoly.State{
		Version:        modelmonopoly.Version1,
		Pass:           roomPass,
		InitialBalance: initialBalance,
		Players:        []modelmonopoly.StatePlayer{player},
	}

	if err = u.repositoryNats.MonopolySetState(ctx, roomID, state); err != nil {
		err = util.ErrWrap(prefix, err, "set state")
		return
	}

	return
}

func (u *usecase) JoinRoom(ctx context.Context, playerID, playerName, roomID, roomPass string) error {
	prefix := u.prefix + ".JoinRoom"

	state, err := u.repositoryNats.MonopolyGetState(ctx, roomID)
	if err != nil {
		if errors.Is(err, constant.ErrRoomNotFound) {
			return err
		}
		return util.ErrWrap(prefix, err, "get state")
	}

	// check if pass correct
	if state.Pass != roomPass {
		return constant.ErrUnauthenticated
	}

	// insert player
	state.Players = append(state.Players, modelmonopoly.StatePlayer{
		ID:      playerID,
		Name:    playerName,
		Balance: state.InitialBalance,
	})

	// insert log
	state.Logs = append(state.Logs, logJoinRoom{playerName: playerName}.ToStateLog())

	// update to database
	if err = u.repositoryNats.MonopolySetState(ctx, roomID, state); err != nil {
		return util.ErrWrap(prefix, err, "set state")
	}

	return nil
}

func (u *usecase) LeaveRoom(ctx context.Context, playerID, roomID string) error {
	prefix := u.prefix + ".LeaveRoom"

	state, err := u.repositoryNats.MonopolyGetState(ctx, roomID)
	if err != nil {
		if errors.Is(err, constant.ErrRoomNotFound) {
			return err
		}
		return util.ErrWrap(prefix, err, "get state")
	}

	// insert log
	for _, player := range state.Players {
		if player.ID == playerID {
			state.Logs = append(state.Logs, logLeaveRoom{playerName: player.Name}.ToStateLog())
			break
		}
	}

	// removing player
	for i, player := range state.Players {
		if player.ID == playerID {
			state.Players = append(state.Players[:i], state.Players[i+1:]...)
		}
	}

	return nil
}

func (u *usecase) GetState(ctx context.Context, playerID, roomID string) (state modelmonopoly.State, err error) {
	prefix := u.prefix + ".GetState"

	state, err = u.repositoryNats.MonopolyGetState(ctx, roomID)
	if err != nil {
		if errors.Is(err, constant.ErrRoomNotFound) {
			return
		}
		err = util.ErrWrap(prefix, err, "get state")
		return
	}

	return
}

func (u *usecase) SubscribeState(ctx context.Context, roomID string) (stateChan <-chan modelmonopoly.State, stop func(), err error) {
	prefix := u.prefix + ".SubscribeState"

	stateChan, stop, err = u.repositoryNats.MonopolySubscribeState(ctx, roomID)
	if err != nil {
		err = util.ErrWrap(prefix, err, "subscribe state")
		return
	}

	return
}
