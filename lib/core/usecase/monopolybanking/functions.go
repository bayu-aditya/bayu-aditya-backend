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
	state.AppendLog(logJoinRoom{playerName: player.Name})

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

	// if player exist, update the name.
	// else, insert new player
	isPlayerExist := false

	// update the name if player exist
	for i, player := range state.Players {
		if player.ID == playerID {
			state.Players[i].Name = playerName
			isPlayerExist = true
			break
		}
	}

	// insert player if player not exist
	if !isPlayerExist {
		state.Players = append(state.Players, modelmonopoly.StatePlayer{
			ID:      playerID,
			Name:    playerName,
			Balance: state.InitialBalance,
		})
	}

	// insert log
	state.AppendLog(logJoinRoom{playerName: playerName})

	// update to database
	if err = u.repositoryNats.MonopolySetState(ctx, roomID, state); err != nil {
		return util.ErrWrap(prefix, err, "set state")
	}

	return nil
}

func (u *usecase) LeaveRoom(ctx context.Context, playerID, roomID, roomPass string) error {
	prefix := u.prefix + ".LeaveRoom"

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

	// insert log
	for _, player := range state.Players {
		if player.ID == playerID {
			state.AppendLog(logLeaveRoom{playerName: player.Name})
			break
		}
	}

	// removing player
	for i, player := range state.Players {
		if player.ID == playerID {
			state.Players = append(state.Players[:i], state.Players[i+1:]...)
		}
	}

	if err = u.repositoryNats.MonopolySetState(ctx, roomID, state); err != nil {
		return util.ErrWrap(prefix, err, "set state")
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

// CreateTransaction
//   - mode: 'pay' or 'ask'
func (u *usecase) CreateTransaction(ctx context.Context, roomID, playerID, targetPlayerID string, amount int64, mode string) error {
	prefix := u.prefix + ".CreateTransaction"

	state, err := u.repositoryNats.MonopolyGetState(ctx, roomID)
	if err != nil {
		return util.ErrWrap(prefix, err, "get state")
	}

	// create transaction
	sourceName := ""
	targetName := ""

	for i, player := range state.Players {
		// for source player
		if player.ID == playerID {
			sourceName = player.Name

			if mode == "pay" {
				state.Players[i].Balance -= amount
			} else {
				state.Players[i].Balance += amount
			}
		}

		// for target player
		if player.ID == targetPlayerID {
			targetName = player.Name

			if mode == "ask" {
				state.Players[i].Balance -= amount
			} else {
				state.Players[i].Balance += amount
			}
		}
	}

	state.AppendLog(logTransaction{
		sourceName: sourceName,
		targetName: targetName,
		mode:       mode,
		amount:     amount,
	})

	if err = u.repositoryNats.MonopolySetState(ctx, roomID, state); err != nil {
		return util.ErrWrap(prefix, err, "set state")
	}

	return nil
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
