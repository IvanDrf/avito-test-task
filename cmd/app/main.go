package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/IvanDrf/avito-test-task/internal/app"
	"github.com/IvanDrf/avito-test-task/internal/config"
)

func main() {
	cfg := config.MustLoad()

	app := app.New(cfg)
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	app.Stop()
}
