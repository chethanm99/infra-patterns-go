package workerpool_test

import (
	"infra-patterns-go/isolation/workerpool"
	"testing"

	"go.uber.org/zap"
)

func TestWorker(t *testing.T) {
	numOfJobs := 5
	nopLogger := zap.NewNop()
	jobs := make(chan int, numOfJobs)
	results := make(chan int, numOfJobs)

	for w := 1; w <= 5; w++ {
		go workerpool.Worker(w, jobs, results, true, nopLogger)
	}

	for j := 1; j <= numOfJobs; j++ {
		jobs <- j
	}
	close(jobs)

	for r := 1; r <= numOfJobs; r++ {
		<-results
	}
	close(results)
}
