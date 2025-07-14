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

	go worker.InitializeAndRunPool(clientRedis)
	select {}
	return
	worker.InitializeWorker(client)

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "3000"
	}

	server.ListenAndServe(appPort)

}
