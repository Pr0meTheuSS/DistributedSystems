package app

import (
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
	r := router.NewRouter(a.container)

	address := fmt.Sprintf("%s:%d",
		a.container.Config.GetHost(),
		a.container.Config.GetPort(),
	)

	a.container.Logger.Info("Starting server",
		zap.String("address", address),
	)

	return http.ListenAndServe(address, r)
}
