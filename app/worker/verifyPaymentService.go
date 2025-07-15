package worker

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type Job struct{}

type PaymentServiceResponse struct {
	Failing         bool `json:"failing"`
	MinResponseTime int  `json:"minResponseTime"`
}

func worker(id int, jobs <-chan Job, clientRedis *redis.Client) {
	defaultURL := os.Getenv("PAYMENT_PROCESSOR_URL_DEFAULT")
	fallbackURL := os.Getenv("PAYMENT_PROCESSOR_URL_FALLBACK")

	for range jobs {
		resp, err := http.Get(defaultURL + "/payments/service-health")

		if err != nil {
			log.Printf("Worker %d: Error fetching from default URL: %v", id, err)
			continue
		}

		responseBody, err := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			log.Printf("Worker %d: Error reading response body: %v", id, err)
			continue
		}

		var paymentServiceResponse PaymentServiceResponse
		if err := json.Unmarshal(responseBody, &paymentServiceResponse); err != nil {
			log.Printf("Worker %d: Error unmarshalling response: %v", id, err)
			continue
		}

		if !paymentServiceResponse.Failing {
			clientRedis.Set(context.Background(), "the_best_url_ever", defaultURL, 0)
			clientRedis.Set(context.Background(), "type_of_processor", "default", 0)
			continue
		}

		clientRedis.Set(context.Background(), "the_best_url_ever", fallbackURL, 0)
		clientRedis.Set(context.Background(), "type_of_processor", "fallback", 0)
	}
}

func dispatcher(_ context.Context, jobs chan<- Job) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		jobs <- Job{}
	}
}

func InitializeAndRunPool(clientRedis *redis.Client) {
	const numWorkers = 5
	const jobChannelBufferSize = 10

	jobs := make(chan Job, jobChannelBufferSize)

	for w := 1; w <= numWorkers; w++ {
		go worker(w, jobs, clientRedis)
	}

	go dispatcher(context.Background(), jobs)
}
