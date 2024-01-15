package usecasemonopolybanking

import (
	"fmt"
	modelmonopoly "github.com/bayu-aditya/bayu-aditya-backend/lib/core/model/monopoly"
	"time"
)

type iLog interface {
	ToStateLog() modelmonopoly.StateLog
}

type logJoinRoom struct {
	playerName string
}

func (l logJoinRoom) ToStateLog() modelmonopoly.StateLog {
	return modelmonopoly.StateLog{
		Datetime: time.Now(),
		Message:  fmt.Sprintf("%s joining room", l.playerName),
	}
}

type logLeaveRoom struct {
	playerName string
}

func (l logLeaveRoom) ToStateLog() modelmonopoly.StateLog {
	return modelmonopoly.StateLog{
		Datetime: time.Now(),
		Message:  fmt.Sprintf("%s leaving room", l.playerName),
	}
}
