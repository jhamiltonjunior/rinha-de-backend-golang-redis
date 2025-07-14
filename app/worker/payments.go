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
	Body []byte
}

var (
	SegureOChann = make(chan PaymentWorker, 150)
)

func InitializeWorker(client *mongo.Client) {
	const numWorkers = 150

	for i := 1; i <= numWorkers; i++ {
		go func(id int) {
			for payment := range SegureOChann {
				body, ok := ProcessPayment(payment.Body)
				if !ok {
					sendToPaymentService(payment.Body, os.Getenv("PAYMENT_PROCESSOR_URL_FALLBACK"))
					database.CreatePaymentHistory(client, body, "fallback")
					continue
				}

				database.CreatePaymentHistory(client, body, "default")
				log.Printf("%v\n", body["amount"])
			}
		}(i)
	}
}

func ProcessPayment(paymentBytes []byte) (map[string]any, bool) {
	var payment map[string]any
	if err := json.Unmarshal(paymentBytes, &payment); err != nil {
		return nil, false
	}

	now := time.Now().UTC()
	payment["requestedAt"] = now.Format(time.RFC3339Nano)

	paymentBytes, err := json.Marshal(payment)
	if err != nil {
		fmt.Println("Erro ao serializar o pagamento:", err)
		return nil, false
	}
	reqURL := os.Getenv("PAYMENT_PROCESSOR_URL_DEFAULT")

	return payment, sendToPaymentService(paymentBytes, reqURL)
}

func sendToPaymentService(paymentBytes []byte, reqURL string) bool {
	_, status := utils.Request("POST", paymentBytes, reqURL+"/payments")
	fmt.Printf("status: %d\n", status)
	return status == 200 || status == 201
}
