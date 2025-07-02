package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
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

	if err = repo.InitSchema(ctx); err != nil {
		log.Fatalf("failed to init table: %v", err.Error())
	}

	router := http.NewServeMux()

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
