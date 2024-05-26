package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app := NewApp(context.Background())
	go app.Listen()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-signals

	app.Shutdown()
}
