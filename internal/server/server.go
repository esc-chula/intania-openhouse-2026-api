package server

import (
	"log"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/handlers"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/middlewares"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/repositories"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/usecases"
	"github.com/esc-chula/intania-openhouse-2026-api/pkg/config"
	"github.com/esc-chula/intania-openhouse-2026-api/pkg/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func InitServer(cfg config.Config) error {
	router := chi.NewMux()
	humaCfg := huma.DefaultConfig("intania-openhouse-2026", "1.0.0")

	router.Use(middleware.Logger)
	if cfg.App().IsProduction {
		humaCfg.DocsPath = ""
		humaCfg.OpenAPIPath = ""
		humaCfg.SchemasPath = ""
		router.Use(middleware.Recoverer)
	}

	// Option is modified from cors.AllowAll()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: cfg.App().AllowedOrigins,
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}))

	api := humachi.New(router, humaCfg)

	// Init Database
	db := database.NewPostgresDB(cfg.Database())

	// Create Repositories
	userRepo := repositories.NewUserRepo(db)

	// Create Usecases
	userUsecase := usecases.NewUserUsecase(userRepo)

	// Create Middleware
	mid := middlewares.NewMiddleware(cfg, api)
	_ = mid

	// Init Handlers
	userGroup := huma.NewGroup(api, "/users")
	handlers.InitUserHandler(userGroup, userUsecase)

	if err := http.ListenAndServe(cfg.App().Address, router); err != nil {
		log.Fatal(err)
	}

	return nil
}
