package ctxcli

import (
	"context"
	"os"
	"os/exec"
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
	testCtx := WithInterrupt(context.Background())
	cli, ok := FromContext(testCtx)
	if !ok {
		t.Fatalf("failed to get cli FromContext: %v", cli)
	}
	defer signal.Stop(cli.sigChan)

	cli.sigChan <- syscall.SIGINT

	waitDone(t, testCtx)

	if testCtx.Err() != context.Canceled {
		t.Fatalf("testCtx should be cancelled: %+v", testCtx)
	}
}

func TestExitIfCancelled(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		testCtx := WithInterrupt(context.Background())
		cli, ok := FromContext(testCtx)
		if !ok {
			t.Fatalf("failed to get cli FromContext: %v", cli)
		}
		defer signal.Stop(cli.sigChan)

		cli.sigChan <- syscall.SIGINT

		waitDone(t, testCtx)
		ExitIfCancelled(testCtx, 1)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestExitIfCancelled")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestPanicIfCancelled(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("the code did not panic")
		}
	}()

	testCtx := WithInterrupt(context.Background())
	cli, ok := FromContext(testCtx)
	if !ok {
		t.Fatalf("failed to get cli FromContext: %v", cli)
	}
	defer signal.Stop(cli.sigChan)

	cli.sigChan <- syscall.SIGINT
	waitDone(t, testCtx)
	PanicIfCancelled(testCtx)
}

func waitDone(t *testing.T, ctx context.Context) {
	select {
	case <-time.After(1 * time.Second):
		t.Fatalf("timeout waiting for Done() from %v", ctx)
	case <-ctx.Done():
	}
}
