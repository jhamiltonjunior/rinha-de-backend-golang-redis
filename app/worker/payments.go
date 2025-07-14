package worker

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/jhamiltonjunior/rinha-de-backend/app/database"
	"github.com/jhamiltonjunior/rinha-de-backend/app/utils"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type PaymentWorker struct {
	Body []byte
	VouTeDarOContexto context.Context
}

var (
	SegureOChann = make(chan PaymentWorker, 300)
)

func InitializeWorker(client *mongo.Client, clientRedis *redis.Client) {
	const numWorkers = 300

	for i := 1; i <= numWorkers; i++ {
		go func(id int) {
			for payment := range SegureOChann {
				theBestURLEver := clientRedis.Get(payment.VouTeDarOContexto, "the_best_url_ever").Val()

				body, ok := ProcessPayment(payment.Body, payment.VouTeDarOContexto, theBestURLEver)
				if !ok {
					fallbackURL := os.Getenv("PAYMENT_PROCESSOR_URL_FALLBACK")

					// talvez isso seja um/o problema
					clientRedis.Set(context.Background(), "the_best_url_ever", fallbackURL, 0)

					sendToPaymentService(payment.Body, fallbackURL, payment.VouTeDarOContexto)
					database.CreatePaymentHistory(client, body, "fallback")

					continue
				}

				database.CreatePaymentHistory(client, body, "default")
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

func ProcessPayment(paymentBytes []byte, ctx context.Context, theBestURLEver string) (map[string]any, bool) {
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

	return payment, sendToPaymentService(paymentBytes, theBestURLEver, ctx)
}

func sendToPaymentService(paymentBytes []byte, reqURL string, ctx context.Context) bool {
	_, status := utils.Request("POST", paymentBytes, reqURL+"/payments", ctx)
	fmt.Printf("status: %d\n", status)
	return status == 200 || status == 201
}
