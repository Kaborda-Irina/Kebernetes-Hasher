package services

import (
	"context"
	"sync"
	"time"

	"github.com/Kaborda-Irina/Kubernetes-Hasher/internal/core/consts"

	"github.com/Kaborda-Irina/Kubernetes-Hasher/pkg/api"

	"github.com/sirupsen/logrus"
)

// WorkerPool launches a certain number of workers for concurrent processing
func (hs HashService) WorkerPool(ctx context.Context, countWorkers int, jobs chan string, results chan api.HashData, logger *logrus.Logger) {
	ctx, cancel := context.WithTimeout(ctx, consts.TimeOut*time.Second)
	defer cancel()
	var wg sync.WaitGroup
	for w := 1; w <= countWorkers; w++ {
		wg.Add(1)
		go hs.Worker(ctx, &wg, jobs, results, logger)
	}
	defer close(results)
	wg.Wait()
}

// Worker gets jobs from a pipe and writes the result to stdout and database
func (hs HashService) Worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan string, results chan<- api.HashData, _ *logrus.Logger) {
	_, cancel := context.WithTimeout(ctx, consts.TimeOut*time.Second)
	defer cancel()
	defer wg.Done()
	for j := range jobs {
		results <- hs.CreateHash(j)
	}
}
