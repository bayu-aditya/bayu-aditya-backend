package usecasemonopolybanking

import (
	"fmt"
	modelmonopoly "github.com/bayu-aditya/bayu-aditya-backend/lib/core/model/monopoly"
	"time"
)

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

type logTransaction struct {
	sourceName string
	targetName string
	mode       string
	amount     int64
}

func (l logTransaction) ToStateLog() modelmonopoly.StateLog {
	if l.targetName == "" {
		l.targetName = "Bank"
	}

	return modelmonopoly.StateLog{
		Datetime: time.Now(),
		Message:  fmt.Sprintf("%s %s to %s with amount %d", l.sourceName, l.mode, l.targetName, l.amount),
	}
}
