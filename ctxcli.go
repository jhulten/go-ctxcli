package ctxcli

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

type CLIContext struct {
	sigChan chan os.Signal
}

// From lambdacontext, an unexported type to be used as key for types in
// this package. This prevents collisions with keys defined in other packages.
type key struct{}

// the key for CLIContext in Contexts.
var contextKey = &key{}

func NewContext(parent context.Context, clictx *CLIContext) context.Context {
	return context.WithValue(parent, contextKey, clictx)
}

func FromContext(ctx context.Context) (*CLIContext, bool) {
	clictx, ok := ctx.Value(contextKey).(*CLIContext)
	return clictx, ok
}

func WithSignalTrap(parent context.Context, sigs ...os.Signal) context.Context {
	cli, ok := FromContext(parent)
	if !ok {
		cli = &CLIContext{
			sigChan: make(chan os.Signal, 2),
		}
	}

	cancelCtx, cancel := context.WithCancel(parent)
	signal.Notify(cli.sigChan, sigs...)

	go func(cCtx context.Context, cliCtx *CLIContext, cancel context.CancelFunc) {
		for {
			select {
			case <-cliCtx.sigChan:
				cancel()
			case <-cCtx.Done():
			}
		}
	}(cancelCtx, cli, cancel)

	return NewContext(cancelCtx, cli)
}

func WithInterrupt(parent context.Context) context.Context {
	return WithSignalTrap(parent, syscall.SIGINT, syscall.SIGTERM)
}

func ExitIfCancelled(ctx context.Context, exitCode int) {
	defer func() {
		if r := recover(); r != nil {
			os.Exit(exitCode)
		}
	}()
	PanicIfCancelled(ctx)
}

func PanicIfCancelled(ctx context.Context) {
	select {
	case <-ctx.Done():
		panic(ctx.Err())
	default:
	}
}
