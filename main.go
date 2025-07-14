package main

import (
	"os"
	"github.com/jhamiltonjunior/rinha-de-backend/app/database"
	"github.com/jhamiltonjunior/rinha-de-backend/app/server"
	"github.com/jhamiltonjunior/rinha-de-backend/app/worker"
	_ "github.com/lib/pq"
)

func main() {
	client := database.InitializeMongoDB()
	clientRedis := database.InitializeRedis()

	if os.Getenv("RUN_VERIFY_PAYMENT_SERVICE") == "true" {
		go worker.InitializeAndRunPool(clientRedis)
	}
	worker.InitializeWorker(client, clientRedis)

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "3000"
	}

	server.ListenAndServe(appPort)
}
