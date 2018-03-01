package ctxcli

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func WithSignalTrap(parent context.Context, sigs ...os.Signal) context.Context {
	ctx, cancel := context.WithCancel(parent)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, sigs...)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-sigChan:
				cancel()
			}
		}
	}()

	return ctx
}

func WithInterrupt(parent context.Context) context.Context {
	return WithSignalTrap(parent, syscall.SIGINT, syscall.SIGTERM)
}

func ExitIfCancelled(ctx context.Context) {
	select {
	case <-ctx.Done():
		os.Exit(1)
	default:
	}
}

func PanicIfCancelled(ctx context.Context) {
	select {
	case <-ctx.Done():
		panic(ctx.Err())
	default:
	}
}
