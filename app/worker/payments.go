package worker

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jhamiltonjunior/rinha-de-backend/app/database"
	"github.com/jhamiltonjunior/rinha-de-backend/app/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type PaymentWorker struct {
	Body map[string]any
}

var (
	Queue = make(chan PaymentWorker, 1000)
)

func InitializeWorker(client *mongo.Client) {
	const numWorkers = 1000

	for i := 1; i <= numWorkers; i++ {
		go func(id int) {
			for payment := range Queue {
				database.CreatePaymentHistory(client, payment.Body, "default")
				log.Printf("%v\n", payment.Body["amount"])
			}
		}(i)
	}
}

func registerPayment(paymentBytes []byte) {
	reqURL := os.Getenv("PAYMENT_PROCESSOR_URL_DEFAULT")
	_, status := utils.Request("POST", paymentBytes, reqURL+"/payments")
	fmt.Printf("status: %d\n", status)
}

func ProcessPayment(paymentBytes []byte) {
	var payment map[string]any
	if err := json.Unmarshal(paymentBytes, &payment); err != nil {
		return
	}

	now := time.Now().UTC()
	payment["requestedAt"] = now.Format(time.RFC3339Nano)

	Queue <- PaymentWorker{
		Body: payment,
	}

	paymentBytes, err := json.Marshal(payment)
	if err != nil {
		fmt.Println("Erro ao serializar o pagamento:", err)
		return
	}

	registerPayment(paymentBytes)
}
