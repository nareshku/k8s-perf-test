package cmd

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/nareshku/k8s-perf-test/pkg/discovery"
	"github.com/nareshku/k8s-perf-test/pkg/summary"
	"github.com/nareshku/k8s-perf-test/pkg/worker"
	"github.com/nareshku/k8s-perf-test/types"
	"k8s.io/client-go/dynamic"
)

func runPerformanceTest(configFile string, duration time.Duration) error {
	// Load configuration
	config, err := loadConfig(configFile)
	if err != nil {
		return fmt.Errorf("error loading config: %v", err)
	}

	// Create root context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	allStats := make(map[string][]*types.CallStats)

	fmt.Printf("Starting performance test for duration: %v\n", duration)
	startTime := time.Now()

	// Start workers for each user
	for _, user := range config.Users {
		wg.Add(1)
		go func(userConfig types.UserConfig) {
			defer wg.Done()

			k8sConfig, err := createK8sConfig(config.Cluster, userConfig.Token)
			if err != nil {
				fmt.Printf("Error creating k8s config for user %s: %v\n", userConfig.Username, err)
				return
			}

			discoveryClient, err := discovery.NewResourceDiscovery(k8sConfig, config.Cluster.IgnoreResources)
			if err != nil {
				fmt.Printf("Error creating discovery client for user %s: %v\n", userConfig.Username, err)
				return
			}

			resources, err := discoveryClient.GetAPIResources()
			if err != nil {
				fmt.Printf("Error discovering resources for user %s: %v\n", userConfig.Username, err)
				return
			}

			dynamicClient, err := dynamic.NewForConfig(k8sConfig)
			if err != nil {
				fmt.Printf("Error creating dynamic client for user %s: %v\n", userConfig.Username, err)
				return
			}

			w := worker.NewWorker(dynamicClient)
			w.ExecuteCallsWithDuration(ctx, resources, userConfig.Concurrency, duration)

			allStats[userConfig.Username] = w.GetStats()
		}(user)
	}

	wg.Wait()

	totalDuration := time.Since(startTime)
	fmt.Printf("\nTest completed in %v\n", totalDuration)
	summary.PrintSummary(allStats)

	return nil
}
