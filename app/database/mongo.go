package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitializeMongoDB() *mongo.Client {
	uri := "mongodb://mongodb:27017"

	clientOptions := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func CreatePaymentHistory(client *mongo.Client, paymentData map[string]interface{}, typeService string) *mongo.Collection {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client.Database("meu_banco").Collection("payment_history").InsertOne(ctx, bson.M{
		"correlationId": paymentData["correlationId"],
		"amount":        paymentData["amount"],
		"requestedAt":   paymentData["requestedAt"],
		"type":         typeService,
	})

	return client.Database("meu_banco").Collection("payment_history")
}
