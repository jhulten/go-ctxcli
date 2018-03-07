package ctxcli

import (
	"context"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func TestWithSignalTrap(t *testing.T) {
	testCtx := WithSignalTrap(context.Background(), syscall.SIGUSR1)
	cli, ok := FromContext(testCtx)
	if !ok {
		t.Fatalf("failed to get cli FromContext: %v", cli)
	}
	defer signal.Stop(cli.sigChan)

	cli.sigChan <- syscall.SIGUSR1

	waitDone(t, testCtx)

	if testCtx.Err() != context.Canceled {
		t.Fatalf("testCtx should be cancelled: %+v", testCtx)
	}
}

func TestWithInterrupt(t *testing.T) {

}

func TestExitIfCancelled(t *testing.T) {
}

func TestPanicIfCancelled(t *testing.T) {
}

func waitDone(t *testing.T, ctx context.Context) {
	select {
	case <-time.After(1 * time.Second):
		t.Fatalf("timeout waiting for Done() from %v", ctx)
	case <-ctx.Done():
	}
}
