package worker

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/nareshku/k8s-perf-test/pkg/discovery"
	"github.com/nareshku/k8s-perf-test/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
)

type Worker struct {
	client    dynamic.Interface
	stats     map[string]*types.CallStats
	statsLock sync.RWMutex
}

func NewWorker(client dynamic.Interface) *Worker {
	return &Worker{
		client: client,
		stats:  make(map[string]*types.CallStats),
	}
}

func (w *Worker) ExecuteCallsWithDuration(ctx context.Context, resources []discovery.APIResource, concurrency int, duration time.Duration) {
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()

	// Initialize stats for each resource
	w.statsLock.Lock()
	for _, resource := range resources {
		w.stats[resource.Name] = &types.CallStats{
			ResourceType: resource.Name,
			StartTime:    time.Now(),
		}
	}
	w.statsLock.Unlock()

	// Start worker goroutines for each resource
	for _, resource := range resources {
		wg.Add(1)
		go func(r discovery.APIResource) {
			defer wg.Done()

			// Keep making calls until context is cancelled
			for {
				select {
				case <-ctx.Done():
					return
				case sem <- struct{}{}: // Acquire semaphore
					func() {
						defer func() { <-sem }() // Release semaphore
						err := w.makeAPICall(ctx, r)

						w.statsLock.Lock()
						stats := w.stats[r.Name]
						atomic.AddInt64(&stats.TotalCalls, 1)
						if err != nil {
							if isStatus5xx(err) {
								atomic.AddInt64(&stats.Errors5xx, 1)
							} else if isStatus4xx(err) {
								atomic.AddInt64(&stats.Errors4xx, 1)
							}
						}
						w.statsLock.Unlock()
					}()
				}
			}
		}(resource)
	}

	// Wait for context to be done
	<-ctx.Done()
	wg.Wait()

	// Calculate final stats
	w.statsLock.Lock()
	for _, stats := range w.stats {
		stats.EndTime = time.Now()
		duration := stats.EndTime.Sub(stats.StartTime).Seconds()
		if duration > 0 {
			stats.CallsPerSec = float64(stats.TotalCalls) / duration
		}
	}
	w.statsLock.Unlock()
}

func (w *Worker) makeAPICall(ctx context.Context, resource discovery.APIResource) error {
	if resource.Namespaced {
		_, err := w.client.Resource(resource.GVR).Namespace("default").List(ctx, metav1.ListOptions{})
		return err
	}
	_, err := w.client.Resource(resource.GVR).List(ctx, metav1.ListOptions{})
	return err
}

func (w *Worker) GetStats() []*types.CallStats {
	w.statsLock.RLock()
	defer w.statsLock.RUnlock()

	stats := make([]*types.CallStats, 0, len(w.stats))
	for _, s := range w.stats {
		stats = append(stats, s)
	}
	return stats
}

// Helper functions for error checking
func isStatus5xx(err error) bool {
	// Implement status code checking logic
	return false
}

func isStatus4xx(err error) bool {
	// Implement status code checking logic
	return false
}
