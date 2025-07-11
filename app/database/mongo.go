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
	uri := "mongodb://root:very_hard_password@mongodb:27017/?authSource=admin"

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



	_, err := client.Database("payment").Collection("payment_history").InsertOne(ctx, bson.M{
		"correlationId": paymentData["correlationId"],
		"amount":        paymentData["amount"],
		"requestedAt":   paymentData["requestedAt"],
		"type":          typeService,
	})

	if err != nil {
		log.Printf("Erro ao inserir pagamento: %v", err)
		return nil
	}

	return client.Database("payment").Collection("payment_history")
}
