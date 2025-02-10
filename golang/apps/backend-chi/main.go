package main

import (
	"bdv-avito-merch/libs/4_infrastructure/db_manager"
	"bdv-avito-merch/libs/5_common/env_vars"
	"bdv-avito-merch/libs/5_common/smart_context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	env_vars.LoadEnvVars() // load env vars from .env file if ENV_PATH is specified
	os.Setenv("LOG_LEVEL", "debug")

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

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"message": "Hello from Chi!"})
	})

	// wsupgrader := ws_server.NewWsUpgrader(logger)

	logger.Info("Server listening on port 9000")
	err = http.ListenAndServe(":9000", r)
	logger.Fatal(err)
}
