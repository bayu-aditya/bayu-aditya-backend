package reponats

import (
	"context"
	"fmt"
	"github.com/bayu-aditya/bayu-aditya-backend/lib/core/config"
	"github.com/bayu-aditya/bayu-aditya-backend/lib/core/util"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func NewNatsRepository(cfg *config.Config) (INatsRepository, error) {
	prefix := "reponats"
	bucketName := fmt.Sprintf("ba-mbs-%s", cfg.Project.Environment)

	natsConn, err := nats.Connect(fmt.Sprintf("%s:%s@%s", cfg.Nats.Username, cfg.Nats.Password, cfg.Nats.Url))
	if err != nil {
		return nil, util.ErrWrap(prefix, err, "new connection")
	}

	js, err := jetstream.New(natsConn)
	if err != nil {
		return nil, util.ErrWrap(prefix, err, "new JetStream")
	}

	kv, err := js.KeyValue(context.Background(), bucketName)
	if err != nil {
		return nil, util.ErrWrap(prefix, err, "new key value")
	}

	return &repository{
		prefix:     prefix,
		bucketName: bucketName,
		conn:       natsConn,
		kv:         kv,
	}, nil
}
