package reponats

import (
	"context"
	"errors"
	"github.com/bayu-aditya/bayu-aditya-backend/lib/core/model/constant"
	modelmonopoly "github.com/bayu-aditya/bayu-aditya-backend/lib/core/model/monopoly"
	"github.com/bayu-aditya/bayu-aditya-backend/lib/core/util"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/sirupsen/logrus"
)

func (r *repository) MonopolySetState(ctx context.Context, roomID string, state modelmonopoly.State) error {
	prefix := r.prefix + ".MonopolySetState"

	stateDecode, err := state.Encode()
	if err != nil {
		return util.ErrWrap(prefix, err, "encode state")
	}

	kvEntry, err := r.kv.Get(ctx, roomID)

	if err != nil && !errors.Is(err, jetstream.ErrKeyNotFound) {
		return util.ErrWrap(prefix, err, "getting key")
	}

	// if key not found, then create it
	if kvEntry == nil {
		if _, err = r.kv.Create(ctx, roomID, stateDecode); err != nil {
			return util.ErrWrap(prefix, err, "create key")
		}
		return nil
	}

	// updating key
	kvRev := kvEntry.Revision()
	if _, err = r.kv.Update(ctx, roomID, stateDecode, kvRev); err != nil {
		return util.ErrWrap(prefix, err, "update key")
	}

	return nil
}

func (r *repository) MonopolyGetState(ctx context.Context, roomID string) (state modelmonopoly.State, err error) {
	prefix := r.prefix + ".MonopolyGetState"

	kvEntry, err := r.kv.Get(ctx, roomID)
	if err != nil {
		if errors.Is(err, jetstream.ErrKeyNotFound) {
			err = constant.ErrRoomNotFound
			return
		}
		err = util.ErrWrap(prefix, err, "get key")
		return
	}

	if err = state.Decode(kvEntry.Value()); err != nil {
		err = util.ErrWrap(prefix, err, "unmarshal")
		return
	}

	return
}

func (r *repository) MonopolySubscribeState(ctx context.Context, roomID string) (stateChan <-chan modelmonopoly.State, stop func(), err error) {
	prefix := r.prefix + ".MonopolySubscribeState"
	channel := make(chan modelmonopoly.State, 1)

	stateChan = channel

	keyWatcher, err := r.kv.Watch(ctx, roomID)
	if err != nil {
		err = util.ErrWrap(prefix, err, "watch key")
		return
	}

	stop = func() {
		if errStop := keyWatcher.Stop(); errStop != nil {
			errStop = util.ErrWrap(prefix, errStop, "stopping key watcher")
			logrus.Error(errStop)
		}
	}

	go func() {
		for entry := range keyWatcher.Updates() {
			if entry == nil {
				continue
			}

			state := modelmonopoly.State{}
			if err = state.Decode(entry.Value()); err != nil {
				err = util.ErrWrap(prefix, err, "unmarshal")
				return
			}
			channel <- state
		}
	}()

	return
}

func (r *repository) Close() {
	r.conn.Close()
}
