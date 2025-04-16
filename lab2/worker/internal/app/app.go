package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"worker/internal/di"
)

type App struct {
	container *di.AppContainer
}

func NewApp(container *di.AppContainer) *App {
	return &App{container: container}
}

func (a *App) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	return a.container.Consumer.Consume(ctx)
}
