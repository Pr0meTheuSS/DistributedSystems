package router

import (
	"manager/internal/handler"
	"manager/internal/service"

	_ "manager/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
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

	// ============== Ping layouts initialization ==============
	crackHashService := service.NewCrackHashService(logger)
	crackHashHandler := handler.NewCrackHashHandler(crackHashService)

	router.Route("/api/hash-cracks", func(router chi.Router) {
		router.Post("/", crackHashHandler.PostCrackHash)
		router.Get("/{id}/progress", crackHashHandler.GetCrackHashProgress)
		router.Get("/{id}/result", crackHashHandler.GetCrackHashResult)
	})

	// =========================================================

	// ===================== Swagger UI ========================
	router.Get("/api/swagger/*", httpSwagger.WrapHandler)
	// =========================================================

	return router
}
