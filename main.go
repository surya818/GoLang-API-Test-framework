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
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/kong/candidate-take-home-exercise-sdet/internal/app"
	"github.com/kong/candidate-take-home-exercise-sdet/internal/config"
	"github.com/kong/candidate-take-home-exercise-sdet/internal/database"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

var (
	Version   string
	Commit    string
	OsArch    string
	GoVersion string
	BuildDate string
)

func main() {
	// Initialize the logger
	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("unable to create logger: %v", err))
	}
	logger.Info("starting candidate-take-home-exercise-sdet",
		zap.String("version", Version),
		zap.String("commit", Commit),
		zap.String("os-arch", OsArch),
		zap.String("go-version", GoVersion),
		zap.String("build-date", BuildDate),
	)

	db, err := database.NewDatabase()
	if err != nil {
		panic(fmt.Sprintf("unable to create database: %v", err))
	}

	// Load the configuration
	config, err := config.NewConfig()
	if err != nil {
		panic(fmt.Sprintf("unable to create config: %v", err))
	}

	// Create a new context with a cancel function
	ctx, cancel := context.WithCancel(context.Background())

	// Handle user break signals and ensure resources are properly cleaned
	breakSignal := make(chan os.Signal, 1)
	signal.Notify(breakSignal, os.Interrupt, syscall.SIGTERM)
	defer cancel()
	defer close(breakSignal)

	// Create the application
	app, err := app.NewApp(app.Opts{
		Config:   config,
		Database: db,
		Logger:   logger,
	})
	if err != nil {
		panic(fmt.Sprintf("unable to create application: %v", err))
	}

	// Start the application and wait for shutdown
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = app.Run(ctx)
	}()

	// Wait for a signal
	<-breakSignal
	logger.Info("user requested shutdown")
	cancel()
	wg.Wait()
	logger.Info("application shutdown")
}
