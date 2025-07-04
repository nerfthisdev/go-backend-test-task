package main

import (
	"context"
	"log"
	"time"

	"github.com/joho/godotenv"
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

	// init context
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)

	defer cancel()

	logger := logger.GetLogger()

	repo, err := repository.Init(ctx)

	if err != nil {
		logger.Fatal("failed to connect to db", zap.String("reason", err.Error()))
	}
	defer repo.DB.Close()
	
	logger.Info("successfully connected to db")
	
	
}
