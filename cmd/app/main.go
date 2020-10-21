package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"reddit-clone-example/internal/config"
	"reddit-clone-example/internal/handle"
	"reddit-clone-example/storage"

	"github.com/grafov/kiwi"
	"github.com/grafov/kiwi/where"
)

func main() {
	// The example uses logfmt-styled logger. See docs and
	// examples about the logger at Kiwi project
	// https://github.com/grafov/kiwi
	kiwi.SinkTo(os.Stdout, kiwi.AsLogfmt()).Start()
	// Uncomment the line to add more context info to the log
	// output. Use it for non production environments only.
	kiwi.With(where.What(where.Function | where.FilePos))

	kiwi.Log(
		"build_at", config.BuildAt,
		"commit", config.GitHash,
		"port", config.App.Port,
	)

	storage.Init()

	<-Shutdown(
		[]os.Signal{syscall.SIGINT, syscall.SIGTERM},
		httpServe(handle.Route()).Shutdown,
	)
}

// Shutdown servers
func Shutdown(sig []os.Signal, shutdown ...func(context.Context) error) <-chan struct{} {
	var (
		// shutdown routines counter
		routine = &sync.WaitGroup{}
		// shutdown end notification chanel
		end = make(chan struct{})
		// shutdown signals chanel
		wait = make(chan os.Signal, 1)
	)
	signal.Notify(wait, sig...)
	// run shutdown function
	go func() {
		<-wait // for shutdown signals
		kiwi.With("do", "graceful shutdown")
		kiwi.Log("info", "received interrupt signal")
		// common shutdown deadline time for all services for shutdown
		interrupt := time.Now().Add(config.App.ShutdownTimeout)
		for i := range shutdown {
			// notify wait group about new routine run
			routine.Add(1)
			go func(run func(context.Context) error) {
				defer routine.Done()
				if run == nil {
					return
				}
				ctx, cancel := context.WithDeadline(context.Background(), interrupt)
				defer cancel()
				if err := run(ctx); err != nil {
					kiwi.Log("info", "exec shutdown", "err", err)
				}
			}(shutdown[i])
		}
		routine.Wait()
		kiwi.Log("info", "halt")
		close(end)
		close(wait)
	}()

	return end
}
