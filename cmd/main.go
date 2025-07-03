package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/nerfthisdev/go-backend-test-task/internal/auth"
	"github.com/nerfthisdev/go-backend-test-task/internal/handler"
	"github.com/nerfthisdev/go-backend-test-task/internal/repository"
)

const defaultTimeout = 5 * time.Second

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err.Error())
	}

	appCtx := context.Background()

	ctx, cancel := context.WithTimeout(appCtx, defaultTimeout)
	defer cancel()

	repo, err := repository.Init(ctx)
	if err != nil {
		log.Fatalf("failed to init repository: %v", err.Error())
	}
	defer repo.DB.Close()
	log.Println("successfully connected to db")

	if err = repository.RunMigrations(repo.DB, "migrations"); err != nil {
		log.Fatalf("failed to init table: %v", err.Error())
	}

	accessTTL, _ := time.ParseDuration(os.Getenv("ACCESS_TOKEN_TTL"))
	refreshTTL, _ := time.ParseDuration(os.Getenv("REFRESH_TOKEN_TTL"))
	signingKey := os.Getenv("JWT_SECRET")

	authservice := auth.NewAuthService(repo, signingKey, accessTTL, refreshTTL)
	authhandler := handler.NewAuthHandler(authservice)

	router := http.NewServeMux()
	router.HandleFunc("GET /api/v1/auth/token", authhandler.AddUser)
	router.HandleFunc("GET /api/v1/auth/user/{id}", authhandler.FetchUser)

	port := ":" + os.Getenv("HTTP_PORT")
	server := http.Server{
		Addr:    port,
		Handler: router,
	}

	log.Printf("starting server on %s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("failed to start server: %v", err.Error())
	}
}
