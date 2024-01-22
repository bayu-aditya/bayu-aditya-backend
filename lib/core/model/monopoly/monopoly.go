package modelmonopoly

import (
	"github.com/vmihailenco/msgpack/v5"
	"time"
)

var (
	Version1 = "v1"
)

type ILog interface {
	ToStateLog() StateLog
}

type (
	State struct {
		Version        string        `json:"version"`
		Pass           string        `json:"pass"`
		InitialBalance int64         `json:"initial_balance"`
		Players        []StatePlayer `json:"players"`
		Logs           []StateLog    `json:"logs"`
	}

	StatePlayer struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Balance int64  `json:"balance"`
	}

	StateLog struct {
		Datetime time.Time `json:"datetime"`
		Message  string    `json:"message"`
	}
)

func (s *State) Encode() ([]byte, error) {
	return msgpack.Marshal(s)
}

func (s *State) Decode(data []byte) error {
	return msgpack.Unmarshal(data, s)
}

func (s *State) AppendLog(log ILog) {
	s.Logs = append(s.Logs, log.ToStateLog())
}
