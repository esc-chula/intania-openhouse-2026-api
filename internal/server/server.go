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
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/extra/bundebug"
)

func InitServer(cfg config.Config, db *bun.DB) error {
	router := chi.NewMux()
	humaCfg := huma.DefaultConfig("intania-openhouse-2026", "1.0.0")

	// Setup request error logger
	humaCfg.Transformers = append(humaCfg.Transformers, ErrorCaptureTransformer)

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

	// Setup request error logger
	api.UseMiddleware(ErrorRecorderMiddleware)
	api.UseMiddleware(ErrorLoggerMiddleware)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Init Database
	if db == nil {
		db = database.NewPostgresDB(cfg.Database())
	}
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
	activityRepo := repositories.NewActivityRepo(db)
	stampRepo := repositories.NewStampRepo(db)

	// Create Transactioner
	transactioner := baserepo.NewTransactioner(db)

	// Create Usecases
	userUsecase := usecases.NewUserUsecase(userRepo, stampRepo, transactioner)
	workshopUsecase := usecases.NewWorkshopUsecase(workshopRepo)
	bookingUsecase := usecases.NewBookingUsecase(bookingRepo, workshopRepo, userRepo, transactioner)
	checkInUsecase := usecases.NewCheckInUsecase(bookingRepo, boothRepo, userRepo)
	stampUsecase := usecases.NewStampUsecase(stampRepo, bookingRepo, boothRepo)
	activityUsecase := usecases.NewActivityUsecase(activityRepo)

	// Register Handler
	userGroup := huma.NewGroup(api, "/users")
	workshopGroup := huma.NewGroup(api, "/workshops")
	checkInGroup := huma.NewGroup(api, "/check-in")
	activityGroup := huma.NewGroup(api, "/activities")
	stampGroup := huma.NewGroup(api, "/stamps")

	userGroup.UseMiddleware(mid.WithAuthContext)
	workshopGroup.UseMiddleware(mid.WithAuthContext)
	checkInGroup.UseMiddleware(mid.WithAuthContext)
	activityGroup.UseMiddleware(mid.WithAuthContext)
	stampGroup.UseMiddleware(mid.WithAuthContext)

	handlers.InitUserHandler(userGroup, userUsecase, stampUsecase, mid)
	handlers.InitWorkshopHandler(workshopGroup, workshopUsecase, mid)
	handlers.InitBookingHandler(workshopGroup, userGroup, bookingUsecase, userUsecase, mid)
	handlers.InitCheckInHandler(checkInGroup, checkInUsecase, mid)
	handlers.InitActivityHandler(activityGroup, activityUsecase, mid)
	handlers.InitStampHandler(stampGroup, userGroup, stampUsecase, userUsecase, mid)

	if err := http.ListenAndServe(cfg.App().Address, router); err != nil {
		log.Fatal(err)
	}

	return nil
}
