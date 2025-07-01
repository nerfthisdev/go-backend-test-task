package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/nerfthisdev/go-backend-test-task/internal/repository"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	repo, err := repository.Init(ctx)
	if err != nil {
		log.Fatalf("failed to init repository: %v", err.Error())
	}

	defer repo.DB.Close()

	log.Println("successfully connected to db")

	router := http.NewServeMux()

	port := ":" + os.Getenv("HTTP_PORT")

	server := http.Server{
		Addr:    port,
		Handler: router,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("failed to start http server: %v", err.Error())
	}
}
