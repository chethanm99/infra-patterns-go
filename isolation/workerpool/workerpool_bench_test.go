package workerpool_test

import (
	"fmt"
	"infra-patterns-go/isolation/workerpool"
	"testing"

	"go.uber.org/zap"
)

func BenchmarkWorkerPoolThroughPut(b *testing.B) {
	for _, workers := range []int{1, 2, 4, 8, 16} {
		b.Run(fmt.Sprintf("%d-workers", workers), func(b *testing.B) {
			jobs := make(chan int, b.N)
			results := make(chan int, b.N)
			logger := zap.NewNop()

			for w := 1; w <= workers; w++ {
				go workerpool.Worker(w, jobs, results, true, logger)
			}

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				jobs <- i
			}

			for i := 0; i < b.N; i++ {
				<-results
			}

			close(jobs)
		})
	}

}
