package main

import (
	"bdv-avito-merch/libs/1_domain_methods/handlers"
	"bdv-avito-merch/libs/3_infrastructure/db_manager"
	"bdv-avito-merch/libs/4_common/env_vars"
	"bdv-avito-merch/libs/4_common/middleware"
	"bdv-avito-merch/libs/4_common/safe_go"
	"bdv-avito-merch/libs/4_common/smart_context"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	env_vars.LoadEnvVars() // load env vars from .env file if ENV_PATH is specified
	BACKEND_PORT, ok := os.LookupEnv("BACKEND_PORT")
	if !ok {
		BACKEND_PORT = "8080"
	}

	logger := smart_context.NewSmartContext()

	dbm, err := db_manager.NewDbManager(logger)
	if err != nil {
		logger.Fatalf("Error connecting to database: %v", err)
	}
	logger = logger.WithDbManager(dbm)
	logger = logger.WithDB(dbm.GetGORM())

	r := chi.NewRouter()

	r.Use(chi_middleware.Logger)
	r.Use(chi_middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Requested-With", "X-Request-Id", "X-Session-Id", "Apikey", "X-Api-Key"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Post("/api/auth", handlers.LoginHandler(logger))

	r.Get("/api/info", middleware.WithSmartContext(logger, handlers.InfoHandler))
	r.Post("/api/sendCoin", middleware.WithSmartContext(logger, handlers.SendCoinHandler))
	r.Get("/api/buy/{item}", middleware.WithSmartContext(logger, handlers.BuyItemHandler))

	for i := 0; i < 10; i++ {
		safe_go.SafeGo(logger, func() {
			handlers.HandlersWorker()
		})
	}

	logger.Info("Server listening on port " + BACKEND_PORT)
	err = http.ListenAndServe(":"+BACKEND_PORT, r)
	logger.Fatal(err)
}
