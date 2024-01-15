package reponats

import (
	"context"
	modelmonopoly "github.com/bayu-aditya/bayu-aditya-backend/lib/core/model/monopoly"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type INatsRepository interface {
	MonopolySetState(ctx context.Context, roomID string, state modelmonopoly.State) error
	MonopolyGetState(ctx context.Context, roomID string) (modelmonopoly.State, error)
	MonopolySubscribeState(ctx context.Context, roomID string) (state <-chan modelmonopoly.State, stop func(), err error)
	Close()
}

type repository struct {
	prefix     string
	bucketName string
	conn       *nats.Conn
	kv         jetstream.KeyValue
}
