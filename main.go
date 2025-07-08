package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"github.com/valyala/fasthttp"
)

func main() {
	// Load environment variables

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

	requestHandler := func(ctx *fasthttp.RequestCtx) {
		db, err := sql.Open("postgres", postgresURL)
		if err != nil {
			ctx.Error("Failed to create database connection", fasthttp.StatusInternalServerError)
			log.Printf("sql.Open error: %v", err)
			return
		}
		defer db.Close()

		err = db.Ping()
		if err != nil {
			ctx.Error("Failed to ping database", fasthttp.StatusInternalServerError)
			log.Printf("db.Ping error: %v", err)
			return
		}

		nc, err := nats.Connect(natsURL)
		if err != nil {
			ctx.Error("Failed to connect to NATS", fasthttp.StatusInternalServerError)
			log.Printf("NATS connect error: %v", err)
			return
		}
		defer nc.Close()

		hostname, _ := os.Hostname()

		log.Println("I am One API, I am the only one!")

		ctx.SetContentType("text/plain; charset=utf-8")
		fmt.Fprintf(ctx, "Hello World from fasthttp!\nContainer: %s\nSuccessfully connected to PostgreSQL and NATS!", hostname)
	}

	log.Printf("Server starting with fasthttp on port %s", appPort)
	if err := fasthttp.ListenAndServe(":" + appPort, requestHandler); err != nil {
		log.Fatalf("fasthttp.ListenAndServe error: %v", err)
	}
}
