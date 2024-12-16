// Copyright © 2024 Kong Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package app

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kong/candidate-take-home-exercise-sdet/internal/config"
	"github.com/kong/candidate-take-home-exercise-sdet/internal/server"
	"go.uber.org/zap"
	"golang.org/x/exp/rand"
)

const port = 18080

// Opts are the options used to create a new application.
type Opts struct {
	// Config is the configuration for the application.
	Config *config.Config
	// Database is the database to use for retrieving/storing data.
	Database *sql.DB
	// Logger is the logger to use for logging.
	Logger *zap.Logger
}

// Application instance.
type App struct {
	config   *config.Config
	database *sql.DB
	logger   *zap.Logger
	server   *http.Server
}

// NewApp creates and instance of the application.
func NewApp(opts Opts) (*App, error) {
	// Set up Gorilla Mux router
	router := mux.NewRouter()
	handlers, err := server.NewHandler(server.Opts{
		Config:   opts.Config,
		Database: opts.Database,
		Logger:   opts.Logger,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create handlers: %w", err)
	}

	// Register the /token endpoint
	router.HandleFunc("/v1/token", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateTokenHandler(w, r)
	}).Methods("POST")

	// Register endpoints for services
	// Create a new service
	router.HandleFunc("/v1/services", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateServiceHandler(w, r)
	}).Methods("POST")

	// List all services
	router.HandleFunc("/v1/services", func(w http.ResponseWriter, r *http.Request) {
		handlers.ListServicesHandler(w, r)
	}).Methods("GET")

	// Get a specific service by ID
	router.HandleFunc("/v1/services/{serviceId}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetServiceHandler(w, r)
	}).Methods("GET")

	// Update a specific service by ID
	router.HandleFunc("/v1/services/{serviceId}", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateServiceHandler(w, r)
	}).Methods("PATCH")

	// Delete a specific service by ID
	router.HandleFunc("/v1/services/{serviceId}", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteServiceHandler(w, r)
	}).Methods("DELETE")

	// Register endpoints for service versions
	// Create a new version for a specific service
	router.HandleFunc("/v1/services/{serviceId}/versions", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateServiceVersionHandler(w, r)
	}).Methods("POST")

	// List all versions for a specific service
	router.HandleFunc("/v1/services/{serviceId}/versions", func(w http.ResponseWriter, r *http.Request) {
		handlers.ListServiceVersionsHandler(w, r)
	}).Methods("GET")

	// Get a specific version by ID for a specific service
	router.HandleFunc("/v1/services/{serviceId}/versions/{versionId}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetServiceVersionHandler(w, r)
	}).Methods("GET")

	// Update a specific version by ID for a specific service
	router.HandleFunc("/v1/services/{serviceId}/versions/{versionId}", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateServiceVersionHandler(w, r)
	}).Methods("PATCH")

	// Delete a specific version by ID for a specific service
	router.HandleFunc("/v1/services/{serviceId}/versions/{versionId}", func(w http.ResponseWriter, r *http.Request) {
		// 20% chance to introduce a timeout for testing.
		if rand.Float64() < 0.2 {
			// Create a context with a timeout; simulate a long timeout for testing.
			timeoutDuration := opts.Config.RequestTimeout * 1000
			ctx, cancel := context.WithTimeout(r.Context(), timeoutDuration)
			_ = r.WithContext(ctx)
			defer cancel()

			// Simulate a request that hangs indefinitely unless canceled.
			<-ctx.Done()
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
			return
		}

		// No timeout introduced, proceed with the handler normally.
		handlers.DeleteServiceVersionHandler(w, r)
	}).Methods("DELETE")

	// Create an HTTP server with the router
	server := &http.Server{
		Handler:           router,
		ReadTimeout:       opts.Config.RequestTimeout,
		ReadHeaderTimeout: opts.Config.RequestTimeout,
		WriteTimeout:      opts.Config.RequestTimeout,
		Addr:              fmt.Sprintf(":%d", port),
	}
	return &App{
		config:   opts.Config,
		database: opts.Database,
		logger:   opts.Logger,
		server:   server,
	}, nil
}

// Run starts the HTTP server and listens for incoming requests.
func (a *App) Run(ctx context.Context) error {
	// Start the HTTP server
	go func() {
		a.logger.Info("starting server on", zap.Int("port", port))
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Fatal("server failed to start", zap.Error(err))
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()

	// Shutdown server gracefully
	a.logger.Info("shutdown signal received")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := a.server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("failed to gracefully shutdown server: %w", err)
	}
	a.logger.Info("server gracefully stopped")
	return nil
}
