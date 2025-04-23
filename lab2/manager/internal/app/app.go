package app

import (
	"context"
	"fmt"
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
			a.container.Logger.Fatal("WorkerResponseConsumer failed", zap.Error(err))
		}
	}()

	// go func() {
	// 	stop := make(chan os.Signal, 1)
	// 	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	// 	<-stop
	// 	a.container.Logger.Info("Received shutdown signal, stopping application...")
	// 	cancel()
	// 	time.Sleep(2 * time.Second)
	// }()

	r := router.NewRouter(a.container)
	address := fmt.Sprintf("%s:%d", a.container.Config.GetHost(), a.container.Config.GetPort())

	a.container.Logger.Info("Starting server", zap.String("address", address))
	return http.ListenAndServe(address, r)
}
