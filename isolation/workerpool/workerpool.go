package workerpool

import (
	"time"

	"go.uber.org/zap"
)

func Worker(id int, jobs <-chan int, results chan<- int, isBenchmark bool, logger *zap.Logger) {

	for j := range jobs {
		if !isBenchmark {
			logger.Info("Job started by worker", zap.Int("job_id", j), zap.Int("worker_id", id))
		}

		if !isBenchmark {
			time.Sleep(time.Second)
		}

		if !isBenchmark {
			logger.Info("Job completed by worker", zap.Int("job_id", j), zap.Int("worker_id", id))
		}

		results <- cpuwork(j)
	}
}

func cpuwork(n int) int {
	x := 0
	for i := 0; i < 10000; i++ {
		x += i * n
	}
	return x
}
