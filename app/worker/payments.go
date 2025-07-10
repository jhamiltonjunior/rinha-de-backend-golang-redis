package worker

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jhamiltonjunior/rinha-de-backend/app/utils"
)

type Worker struct{
	ID   int
	Data string
}

var (
	queueJob   = make(chan Worker, 1000)
)

func InitializeWorker() {
	const numWorkers = 1000 

	for i := 1; i <= numWorkers; i++ {
		go func(id int) {
			for job := range queueJob {
				log.Printf("[Worker %d] job #%d (%s)\n", id, job.ID, job.Data)
				time.Sleep(2 * time.Second)
				log.Printf("[Worker %d] job #%d\n", id, job.ID)
			}
		}(i)
	}

	log.Printf("🚀 %d workers tests\n", numWorkers)
}

func registerPayment(paymentBytes []byte, defaultURL bool) {
	reqURL := os.Getenv("PAYMENT_PROCESSOR_URL_DEFAULT")
	if !defaultURL {
		reqURL = os.Getenv("PAYMENT_PROCESSOR_URL_FALLBACK")
	}

	_, statusCode := utils.Request("POST", paymentBytes, reqURL)
	if statusCode < 200 || statusCode > 299 {
		registerPayment(paymentBytes, !defaultURL)
		return
	}
}

func ProcessPayment(paymentBytes []byte) {
	var payment map[string]any
	if err := json.Unmarshal(paymentBytes, &payment); err != nil {
		return
	}

	now := time.Now().UTC()
	payment["requestedAt"] = now

	paymentBytes, err := json.Marshal(payment)
	if err != nil {
		fmt.Println("Erro ao serializar o pagamento:", err)
		return
	}

	registerPayment(paymentBytes, true)
}