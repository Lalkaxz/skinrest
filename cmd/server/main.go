package main

import (
	"SkinRest/config"
	"SkinRest/internal/api"
	"SkinRest/internal/database"
	"SkinRest/internal/middleware"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.GetConfig()           // get configuration
	logger := middleware.NewLogger(cfg) // get logger

	db := database.New() // initialize database

	r := api.NewRouter(logger, db) // initialize new router

	var addr string
	if cfg.Server.Host == "localhost" {
		if cfg.Server.ApiEnv == "local" {
			addr = fmt.Sprintf("localhost:%d", cfg.Server.Port)
		} else {
			addr = fmt.Sprintf("0.0.0.0:%d", cfg.Server.Port)
		}
	} else {
		addr = fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	}

	logger.Info("Starting on " + addr)
	logger.Sugar().Infof("Developed with Gin Framework version: %s", gin.Version)
	if err := r.Run(addr); err != nil { // run server
		log.Fatal(err)
	}
}

// $env:SERVER_HOST = "localhost"; `
// $env:SERVER_PORT = "8081"; `
// $env:GIN_MODE = "debug"; `
// $env:API_ENV = "local"; `
// $env:DATABASE_DRIVER = "postgres"; `
// $env:DATABASE_HOST = "localhost"; `
// $env:DATABASE_PORT = "5432"; `
// $env:DATABASE_USER = "postgres"; `
// $env:DATABASE_PASSWORD = "0000"; `
// $env:DATABASE_NAME = "skinRestDB"; `
// $env:DATABASE_SSL = "disable"; `
// $env:AUTH_JWT_SECRET = "8ddeefb1f8c17f17864b0512c5148319848614a11efaed0b247c5cb2e19122e2"; `
// go run ./cmd/server/main.go
