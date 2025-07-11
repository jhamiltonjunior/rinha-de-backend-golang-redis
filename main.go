package main

import (
	"log"
	"os"

	"github.com/jhamiltonjunior/rinha-de-backend/app/server"
	"github.com/jhamiltonjunior/rinha-de-backend/app/worker"
	_ "github.com/lib/pq"
)

func main() {
	worker.InitializeWorker()

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "3000"
	}

	postgresURL := os.Getenv("POSTGRES_URL")
	if postgresURL == "" {
		log.Fatal("POSTGRES_URL environment variable is not set")
	}

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		log.Fatal("NATS_URL environment variable is not set")
	}

	server.ListenAndServe(appPort, postgresURL, natsURL)
}
