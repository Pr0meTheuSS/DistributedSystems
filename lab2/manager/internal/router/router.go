package router

import (
	"manager/internal/di"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(container *di.AppContainer) chi.Router {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/ping", container.PingHandler.GetPing)

	router.Route("/api/hash-cracks", func(router chi.Router) {
		router.Post("/", container.CrackHashHandler.PostCrackHash)
		router.Get("/{id}/progress", container.CrackHashHandler.GetCrackHashProgress)
		router.Get("/{id}/result", container.CrackHashHandler.GetCrackHashResult)
	})

	router.Get("/api/swagger/*", httpSwagger.WrapHandler)

	return router
}
