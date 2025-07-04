package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/nerfthisdev/go-backend-test-task/internal/auth"
	"github.com/nerfthisdev/go-backend-test-task/internal/config"
	"github.com/nerfthisdev/go-backend-test-task/internal/handler"
	"github.com/nerfthisdev/go-backend-test-task/internal/logger"
	"github.com/nerfthisdev/go-backend-test-task/internal/repository"
	"go.uber.org/zap"
)

const defaultTimeout = 5 * time.Second

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
		logger.Fatal("failed to initialize db", zap.Error(err))
	}

	tokenRepo := repository.NewTokenRepository(dbpool)
	userRepo := repository.NewUserRepository(dbpool)

	if err != nil {
		logger.Fatal("failed to connect to db", zap.String("reason", err.Error()))
	}

	logger.Info("successfully connected to db")

	accessTTL, err := time.ParseDuration(cfg.AccessTTL)
	if err != nil {
		logger.Fatal("invalid ACCESS_TOKEN_TTL", zap.Error(err))
	}

	jwtService := auth.NewJwtService(cfg.JWTSecret, accessTTL)

	authService := auth.NewAuthService(tokenRepo, jwtService, userRepo, &logger)

	authhandler := handler.NewAuthHandler(authService)

	router := http.NewServeMux()
	router.HandleFunc("GET /api/v1/auth/{guid}", authhandler.Authorize)

	port := ":" + os.Getenv("HTTP_PORT")
	server := http.Server{
		Addr:    port,
		Handler: router,
	}

	log.Printf("starting server on %s", port)
	if err := server.ListenAndServe(); err != nil {
		logger.Fatal("failed to start server ", zap.String("reason", err.Error()))
	}

}
