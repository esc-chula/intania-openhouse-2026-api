package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/handlers"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/middlewares"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/repositories"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/usecases"
	"github.com/esc-chula/intania-openhouse-2026-api/pkg/baserepo"
	"github.com/esc-chula/intania-openhouse-2026-api/pkg/config"
	"github.com/esc-chula/intania-openhouse-2026-api/pkg/database"
	"github.com/esc-chula/intania-openhouse-2026-api/pkg/firebaseadapter"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/uptrace/bun/extra/bundebug"
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Init Database
	db := database.NewPostgresDB(cfg.Database())
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(!cfg.App().IsProduction),
	))

	// Initialize Middleware
	firebaseAdapter := firebaseadapter.InitFirebaseAuthAdapter(ctx, cfg)
	mid := middlewares.NewMiddleware(cfg, api, firebaseAdapter)

	// Create Repositories
	userRepo := repositories.NewUserRepo(db)
	workshopRepo := repositories.NewWorkshopRepo(db)
	bookingRepo := repositories.NewBookingRepo(db)
	boothRepo := repositories.NewBoothRepo(db)

	// Create Transactioner
	transactioner := baserepo.NewTransactioner(db)

	// Create Usecases
	userUsecase := usecases.NewUserUsecase(userRepo)
	workshopUsecase := usecases.NewWorkshopUsecase(workshopRepo)
	bookingUsecase := usecases.NewBookingUsecase(bookingRepo, workshopRepo, userRepo, transactioner)
	checkInUsecase := usecases.NewCheckInUsecase(bookingRepo, boothRepo, userRepo)

	// Register Handler
	userGroup := huma.NewGroup(api, "/users")
	workshopGroup := huma.NewGroup(api, "/workshops")
	checkInGroup := huma.NewGroup(api, "/check-in")

	handlers.InitUserHandler(userGroup, userUsecase, mid)
	handlers.InitWorkshopHandler(workshopGroup, workshopUsecase, mid)
	handlers.InitBookingHandler(workshopGroup, userGroup, bookingUsecase, userUsecase, mid)
	handlers.InitCheckInHandler(checkInGroup, checkInUsecase, mid)

	if err := http.ListenAndServe(cfg.App().Address, router); err != nil {
		log.Fatal(err)
	}

	return nil
}
