package cmd

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
	"k8s.io/client-go/rest"

	"github.com/nareshku/k8s-perf-test/types"
)

func loadConfig(path string) (*types.Config, error) {
	data, err := os.ReadFile(path)
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
		caData, err := os.ReadFile(clusterConfig.CAPath)
		if err != nil {
			return nil, fmt.Errorf("error reading CA file: %v", err)
		}
		config.TLSClientConfig = rest.TLSClientConfig{
			CAData: caData,
		}
	}

	return config, nil
}
