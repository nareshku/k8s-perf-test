package types

import "time"

type ClusterConfig struct {
	APIServer       string   `yaml:"apiServer"`
	CAPath          string   `yaml:"caPath,omitempty"`   // Optional: path to CA cert
	Insecure        bool     `yaml:"insecure,omitempty"` // Optional: skip TLS verify
	QPS             float32  `yaml:"qps,omitempty"`      // Default will be 50 if not set
	Burst           int      `yaml:"burst,omitempty"`
	IgnoreResources []string `yaml:"ignoreResources,omitempty"`
}

type UserConfig struct {
	Username    string `yaml:"username"`
	Token       string `yaml:"token"`
	Concurrency int    `yaml:"concurrency"`
}

type Config struct {
	Cluster ClusterConfig `yaml:"cluster"`
	Users   []UserConfig  `yaml:"users"`
}

type CallStats struct {
	ResourceType string
	TotalCalls   int64
	Errors5xx    int64
	Errors4xx    int64
	StartTime    time.Time
	EndTime      time.Time
	CallsPerSec  float64
}
