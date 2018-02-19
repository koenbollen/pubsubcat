package utils

import (
	"context"
	"os"
	"os/signal"
)

// CancelOnSignal will wait for one of the given signals in a goroutine and
// call the cancel function when that happens.
func CancelOnSignal(ctx context.Context, cancel context.CancelFunc, signals ...os.Signal) {
	waitChannel := make(chan os.Signal, 1)
	signal.Notify(waitChannel, signals...)
	go func() {
		select {
		case <-waitChannel:
		case <-ctx.Done():
		}
		signal.Reset(signals...)
		cancel()
	}()
}
