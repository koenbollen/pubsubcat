package utils

import (
	"context"
	"os"
	"testing"
)

func TestCancelOnSignal(t *testing.T) {
	type args struct {
		ctx     context.Context
		cancel  context.CancelFunc
		signals []os.Signal
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CancelOnSignal(tt.args.ctx, tt.args.cancel, tt.args.signals...)
		})
	}
}
