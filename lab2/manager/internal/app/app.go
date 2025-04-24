package app

import (
	"context"
	"fmt"
	"log"
	"manager/internal/di"
	"manager/internal/router"
	"net/http"

	"go.uber.org/zap"
)

type App struct {
	container *di.AppContainer
}

func NewApp(container *di.AppContainer) *App {
	return &App{container: container}
}

func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := a.container.WorkerResponseConsumer.Consume(ctx); err != nil {
			// a.container.Logger.Fatal("WorkerResponseConsumer failed", zap.Error(err))
		}
	}()
	if err := a.container.CrackHashService.HandleDeadTasks(context.Background()); err != nil {
		log.Fatal(err.Error())
	}

	r := router.NewRouter(a.container)
	address := fmt.Sprintf("%s:%d", a.container.Config.GetHost(), a.container.Config.GetPort())

	a.container.Logger.Info("Starting server", zap.String("address", address))
	return http.ListenAndServe(address, r)
}
