package main

import (
    "context"
    "flag"
    "fmt"
    "os"
    "sync"
    "io/ioutil"
    "time"
    
    "gopkg.in/yaml.v2"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/dynamic"
    
    "github.com/nareshku/k8s-perf-test/types"
    "github.com/nareshku/k8s-perf-test/pkg/discovery"
    "github.com/nareshku/k8s-perf-test/pkg/worker"
    "github.com/nareshku/k8s-perf-test/pkg/summary"
)

func loadConfig(path string) (*types.Config, error) {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("error reading config file: %v", err)
    }

    var config types.Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("error parsing config file: %v", err)
    }

    return &config, nil
}

func createK8sConfig(clusterConfig types.ClusterConfig, userToken string) (*rest.Config, error) {
    config := &rest.Config{
        Host:        clusterConfig.APIServer,
        BearerToken: userToken,
    }

    // Set QPS and Burst with defaults if not configured
    if clusterConfig.QPS == 0 {
        config.QPS = 50
    } else {
        config.QPS = clusterConfig.QPS
    }
    
    if clusterConfig.Burst == 0 {
        config.Burst = 100
    } else {
        config.Burst = clusterConfig.Burst
    }

    if clusterConfig.Insecure {
        config.TLSClientConfig = rest.TLSClientConfig{
            Insecure: true,
        }
    } else if clusterConfig.CAPath != "" {
        caData, err := ioutil.ReadFile(clusterConfig.CAPath)
        if err != nil {
            return nil, fmt.Errorf("error reading CA file: %v", err)
        }
        config.TLSClientConfig = rest.TLSClientConfig{
            CAData: caData,
        }
    }

    return config, nil
}

func main() {
    configFile := flag.String("config", "config.yaml", "Path to config file")
    duration := flag.Duration("duration", 5*time.Minute, "Duration to run the test (e.g., 1h, 5m, 30s)")
    flag.Parse()

    // Load configuration
    config, err := loadConfig(*configFile)
    if err != nil {
        fmt.Printf("Error loading config: %v\n", err)
        os.Exit(1)
    }
    
    // Create root context
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    var wg sync.WaitGroup
    allStats := make(map[string][]*types.CallStats)
    
    fmt.Printf("Starting performance test for duration: %v\n", *duration)
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

            discoveryClient, err := discovery.NewResourceDiscovery(k8sConfig)
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
            w.ExecuteCallsWithDuration(ctx, resources, userConfig.Concurrency, *duration)
            
            allStats[userConfig.Username] = w.GetStats()
        }(user)
    }

    wg.Wait()

    totalDuration := time.Since(startTime)
    fmt.Printf("\nTest completed in %v\n", totalDuration)
    summary.PrintSummary(allStats)
} 