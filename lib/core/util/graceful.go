package util

import (
	"os"
	"os/signal"
	"syscall"
)

func GracefulShutdown(handler func()) chan struct{} {
	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
		<-sigint

		handler()
		close(idleConnsClosed)
	}()

	return idleConnsClosed
}
