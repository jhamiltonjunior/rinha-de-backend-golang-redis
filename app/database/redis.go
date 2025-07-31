package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitializeRedis() *redis.Client {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "redis_cache:6379",
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	// RedisClient.Set(context.Background(), "payment_history", "", 0).Err()

	return RedisClient
}

func CreatePaymentHistoryInMemory(client *redis.Client, paymentData map[string]any, typeService string) {
	ctx := context.Background()

	key := "payment_history"

	newEntry := map[string]any{
		"correlationId": paymentData["correlationId"],
		"amount":        paymentData["amount"],
		"requestedAt":   paymentData["requestedAt"],
		"type":          typeService,
	}

	entryBytes, err := json.Marshal(newEntry)
	if err != nil {
		log.Printf("Erro ao serializar entrada: %v", err)
		return
	}

	fmt.Printf("Adding entry to payment history: %v\n", newEntry)
	err = client.LPush(ctx, key, entryBytes).Err()
	if err != nil {
		log.Printf("Erro ao adicionar entrada ao histórico de pagamentos: %v", err)
	}
}

func GetPaymentHistoryInMemory(client *redis.Client, from, to string) ([]PaymentHistory, error) {
	isoString := "2006-01-02T15:04:05.000Z"
	fromTime, err := time.Parse(isoString, from)
	if err != nil {
		log.Printf("Erro ao analisar data 'from': %v", err)
		return nil, err
	}

	toTime, err := time.Parse(isoString, to)
	if err != nil {
		log.Printf("Erro ao analisar data 'to': %v", err)
		return nil, err
	}

	ctx := context.Background()
	key := "payment_history"

	dataList, err := client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		log.Printf("Erro ao recuperar histórico do Redis: %v", err)
		return nil, err
	}

	var filteredHistory []PaymentHistory

	for _, item := range dataList {
		var entry map[string]any
		if err := json.Unmarshal([]byte(item), &entry); err != nil {
			continue
		}

		requestedAtStr, ok := entry["requestedAt"].(string)
		if !ok {
			continue
		}

		entryTime, err := time.Parse(isoString, requestedAtStr)
		if err != nil {
			continue
		}

		if entryTime.After(fromTime) && entryTime.Before(toTime) {
			payment := PaymentHistory{
				CorrelationId: entry["correlationId"].(string),
				Amount:        entry["amount"].(float64),
				RequestedAt:   requestedAtStr,
				Type:          entry["type"].(string),
			}
			filteredHistory = append(filteredHistory, payment)
		}
	}

	return filteredHistory, nil
}

func PurgePaymentHistoryInMemory(client *redis.Client) {
	ctx := context.Background()
	key := "payment_history"

	err := client.Del(ctx, key).Err()
	if err != nil {
		log.Printf("Erro ao limpar histórico de pagamentos: %v", err)
	} else {
		fmt.Println("Histórico de pagamentos limpo com sucesso.")
	}
}