package utils_test

import (
	"context"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/koenbollen/pubsubcat/utils"
)

func TestCancelOnSignal(t *testing.T) {
	ctx := context.Background()
	called := false
	cancelFunc := func() {
		called = true
	}

	utils.CancelOnSignal(ctx, cancelFunc, syscall.SIGUSR1)

	syscall.Kill(os.Getpid(), syscall.SIGUSR1)

	time.Sleep(1 * time.Millisecond)
	if !called {
		t.Error("cancelFunc not called on signal")
	}
}
