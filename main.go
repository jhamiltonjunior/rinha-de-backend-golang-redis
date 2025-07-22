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

	worker.InitializeWorker(client, clientRedis)

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "3000"
	}

	// go pingQuantityOfSegureOChann()

	server.ListenAndServe(appPort)
}

// func pingQuantityOfSegureOChann() {
// 	for {
// 		fmt.Printf("Quantidade de pagamentos pendentes: %d\n", len(worker.SegureOChann))
// 		fmt.Printf("Quantidade de pagamentos em retry: %d\n", len(worker.SegureOChann2))
// 		fmt.Println("========================================")
// 		time.Sleep(5 * time.Second)
// 	}
// }
