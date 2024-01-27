package modelmonopoly

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type mockLog struct{}

func (m mockLog) ToStateLog() StateLog {
	return StateLog{
		Datetime: time.Date(2024, 1, 28, 0, 0, 0, 0, time.UTC),
		Message:  "Log mock",
	}
}

func TestState_AppendLog(t *testing.T) {
	type fields struct {
		Version        string
		Pass           string
		InitialBalance int64
		Players        []StatePlayer
		Logs           []StateLog
	}
	type args struct {
		log ILog
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []StateLog
	}{
		{
			name: "append log",
			fields: fields{
				Logs: []StateLog{
					{
						Datetime: time.Date(2024, 1, 27, 0, 0, 0, 0, time.UTC),
						Message:  "Log 1",
					},
				},
			},
			args: args{
				log: mockLog{},
			},
			want: []StateLog{
				{
					Datetime: time.Date(2024, 1, 28, 0, 0, 0, 0, time.UTC),
					Message:  "Log mock",
				},
				{
					Datetime: time.Date(2024, 1, 27, 0, 0, 0, 0, time.UTC),
					Message:  "Log 1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &State{
				Version:        tt.fields.Version,
				Pass:           tt.fields.Pass,
				InitialBalance: tt.fields.InitialBalance,
				Players:        tt.fields.Players,
				Logs:           tt.fields.Logs,
			}
			s.AppendLog(tt.args.log)
			assert.Equal(t, tt.want, s.Logs)
		})
	}
}
