package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/golang-pipeline-test/cmd/pipeline/stage"
)

func main() {

	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)
	quit := make(chan struct{})
	var done sync.WaitGroup

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	first := stage.First(ctx, quit, &done)
	second := stage.Second(ctx, 20, first, &done)
	stage.Final(ctx, 20, second, &done)

	<-signals
	close(quit)
	done.Wait()
	ctx.Done()
}
