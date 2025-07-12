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
	SegureOChann = make(chan PaymentWorker, 100)
)

func InitializeWorker(client *mongo.Client) {
	const numWorkers = 100

	for i := 1; i <= numWorkers; i++ {
		go func(id int) {
			for payment := range SegureOChann {
				body, ok := ProcessPayment(payment.Body)
				if !ok {
					log.Printf("Erro ao processar pagamento: %v\n", body)
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

	return payment, sendToPaymentService(paymentBytes)
}

func sendToPaymentService(paymentBytes []byte) bool {
	reqURL := os.Getenv("PAYMENT_PROCESSOR_URL_DEFAULT")
	_, status := utils.Request("POST", paymentBytes, reqURL+"/payments")
	fmt.Printf("status: %d\n", status)
	return status == 200 || status == 201
}
