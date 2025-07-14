package worker

import (
	"context"
	"log"
	"time"
)

type Job struct {
	ID        int
	CreatedAt time.Time
}

func worker(id int, jobs <-chan Job) {
	for j := range jobs {
		log.Printf("Worker %d: started job %d from time %v", id, j.ID, j.CreatedAt.Format(time.RFC3339))
		
		
		//worker logic here


		time.Sleep(2 * time.Second)
		log.Printf("Worker %d: finished job %d", id, j.ID)
	}
}

func dispatcher(_ context.Context, jobs chan<- Job) {
	log.Println("Dispatcher started.")
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	var jobCounter int

	for t := range ticker.C {
		log.Printf("Dispatcher tick at %v. Creating new job.", t.Format(time.RFC3339))
		jobCounter++
		jobs <- Job{ID: jobCounter, CreatedAt: t}
	}
}

func InitializeAndRunPool() {

	const numWorkers = 5
	const jobChannelBufferSize = 10

	jobs := make(chan Job, jobChannelBufferSize)

	for w := 1; w <= numWorkers; w++ {
		go worker(w, jobs)
	}
	log.Printf("Started %d background workers.", numWorkers)

	go dispatcher(context.Background(), jobs)
}

func main() {

	go InitializeAndRunPool()
}
