package worker

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/jhamiltonjunior/rinha-de-backend/app/database"
	"github.com/jhamiltonjunior/rinha-de-backend/app/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type PaymentWorker struct {
	Body []byte
	VouTeDarOContexto context.Context
}

var (
	SegureOChann = make(chan PaymentWorker, 300)
)

func InitializeWorker(client *mongo.Client) {
	const numWorkers = 300

	for i := 1; i <= numWorkers; i++ {
		go func(id int) {
			for payment := range SegureOChann {
				start := time.Now()
				log.Printf("Worker %d started processing payment: %d", id, start.UnixNano())
				body, ok := ProcessPayment(payment.Body, payment.VouTeDarOContexto)
				if !ok {
					sendToPaymentService(payment.Body, os.Getenv("PAYMENT_PROCESSOR_URL_FALLBACK"), payment.VouTeDarOContexto)
					database.CreatePaymentHistory(client, body, "fallback")
					duration := time.Since(start)
					log.Printf("Worker finished processing payment: %d, duration: %s", id, duration)
					continue
				}
				duration := time.Since(start)
				log.Printf("Worker finished processing payment: %d, duration: %s", id, duration)

				database.CreatePaymentHistory(client, body, "default")
				// log.Printf("%v\n", body["amount"])
			}
		}(i)
	}
}

func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	uuid[8] = uuid[8]&^0xc0 | 0x80
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func ProcessPayment(paymentBytes []byte, ctx context.Context) (map[string]any, bool) {
	var payment map[string]any
	if err := json.Unmarshal(paymentBytes, &payment); err != nil {
		return nil, false
	}

	// payment["correlationId"], _ = newUUID()

	now := time.Now().UTC()
	payment["requestedAt"] = now.Format(time.RFC3339Nano)

	paymentBytes, err := json.Marshal(payment)
	if err != nil {
		fmt.Println("Erro ao serializar o pagamento:", err)
		return nil, false
	}
	reqURL := os.Getenv("PAYMENT_PROCESSOR_URL_DEFAULT")

	return payment, sendToPaymentService(paymentBytes, reqURL, ctx)
}

func sendToPaymentService(paymentBytes []byte, reqURL string, ctx context.Context) bool {
	_, status := utils.Request("POST", paymentBytes, reqURL+"/payments", ctx)
	fmt.Printf("status: %d\n", status)
	return status == 200 || status == 201
}
