package app

import (
	"fmt"
	"manager/internal/config"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type App struct {
	config *config.Config
	logger *zap.Logger
	router chi.Router
}

func (a *App) Run() error {
	a.logger.Info("Start application",
		zap.String("host", a.config.GetHost()),
		zap.Int64("port", a.config.GetPort()))

	serveAddress := fmt.Sprintf("%s:%d", a.config.GetHost(), a.config.GetPort())

	return http.ListenAndServe(serveAddress, a.router)
}

func NewApp(config *config.Config, logger *zap.Logger, router chi.Router) *App {
	app := &App{
		config: config,
		logger: logger,
		router: router,
	}

	return app
}
