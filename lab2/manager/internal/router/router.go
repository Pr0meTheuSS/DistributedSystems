package router

import (
	"manager/internal/handler"
	"manager/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func NewRouter(logger *zap.Logger) chi.Router {
	router := chi.NewRouter()

	// ============== Set router middlewares ==============
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	// ====================================================

	// ============== Ping layouts initialization ==============
	pingService := service.NewPingService(logger)
	pingHandler := handler.NewPingHandler(pingService)
	router.Get("/ping", pingHandler.GetPing)
	// =========================================================

	return router
}
