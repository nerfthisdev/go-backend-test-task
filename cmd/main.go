package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/nerfthisdev/go-backend-test-task/docs" // swagger docs
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/joho/godotenv"
	"github.com/nerfthisdev/go-backend-test-task/internal/auth"
	"github.com/nerfthisdev/go-backend-test-task/internal/config"
	"github.com/nerfthisdev/go-backend-test-task/internal/handler"
	"github.com/nerfthisdev/go-backend-test-task/internal/logger"
	"github.com/nerfthisdev/go-backend-test-task/internal/middleware"
	"github.com/nerfthisdev/go-backend-test-task/internal/repository"
	"go.uber.org/zap"
)

const defaultTimeout = 5 * time.Second

// @title           Go Backend Test Task API
// @version         1.0
// @description     This is the API for the authentication service.
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
func main() {
	// init env
	if err := godotenv.Load(); err != nil {
		log.Fatal(err.Error())
	}

	// init config
	cfg := config.InitConfig()

	// init context
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)

	defer cancel()

	logger := logger.GetLogger()

	dbpool, err := repository.InitDB(ctx, cfg)
	if err != nil {
		logger.Fatal("failed to connect to db", zap.Error(err))
	}

	logger.Info("successfully connected to db")

	defer dbpool.Close()

	tokenRepo := repository.NewTokenRepository(dbpool)
	userRepo := repository.NewUserRepository(dbpool)

	if err = repository.RunMigrations(dbpool, "migrations"); err != nil {
		logger.Fatal("failed to run migrations", zap.Error(err))
	}

	accessTTL, err := time.ParseDuration(cfg.AccessTTL)
	if err != nil {
		logger.Fatal("invalid ACCESS_TOKEN_TTL", zap.Error(err))
	}

	jwtService := auth.NewJwtService(cfg.JWTSecret, accessTTL)

	authService := auth.NewAuthService(tokenRepo, jwtService, userRepo, &logger)

	authHandler := handler.NewAuthHandler(authService)

	router := http.NewServeMux()
	router.HandleFunc("POST /api/v1/auth", authHandler.Authorize)
	router.Handle(
		"POST /api/v1/refresh",
		middleware.Auth(&logger, jwtService, tokenRepo, http.HandlerFunc(authHandler.Refresh)),
	)
	router.Handle(
		"GET /api/v1/me",
		middleware.Auth(&logger, jwtService, tokenRepo, http.HandlerFunc(authHandler.Me)),
	)

	router.Handle(
		"POST /api/v1/deauthorize",
		middleware.Auth(&logger, jwtService, tokenRepo, http.HandlerFunc(authHandler.Deauthorize)),
	)

	router.Handle("/swagger/", httpSwagger.WrapHandler)

	port := ":" + os.Getenv("HTTP_PORT")
	server := http.Server{
		Addr:    port,
		Handler: router,
	}

	logger.Info("starting server on ", zap.String("port", port))
	if err := server.ListenAndServe(); err != nil {
		logger.Fatal("failed to start server ", zap.String("reason", err.Error()))
	}
}
