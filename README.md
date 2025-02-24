# Kubernetes API Performance Testing Tool

A tool for performance testing Kubernetes API server by making concurrent API calls across different resources and users.

## Features

- Auto-discovery of available Kubernetes API resources
- Multi-user support with configurable concurrency
- Duration-based test execution
- Configurable QPS and burst rates
- Automatic filtering of deprecated API resources
- Detailed performance metrics and statistics
- Support for secure and insecure connections

## Prerequisites

- Go 1.21 or higher
- Access to a Kubernetes cluster
- Valid authentication tokens for the users

## Installation

```bash
git clone https://github.com/nareshku/k8s-perf-test
cd k8s-perf-test
go mod tidy
```

## Configuration

Create a `config.yaml` file:

```yaml
cluster:
  apiServer: "https://kubernetes.example.com:6443"
  caPath: "/path/to/ca.crt"  # Optional
  insecure: false           # Optional
  qps: 100                  # Higher QPS to avoid throttling
  burst: 200                # Higher burst to avoid throttling
  ignoreResources:
    - "v1/componentstatuses"
    - "v1/componentstatus"
    - "v1/events"

users:
  - username: "user-1"
    token: "token-1"
    concurrency: 10
  - username: "user-2"
    token: "token-2"
    concurrency: 5 
```

## Usage

Run the performance test:

```bash
# Run with default duration (5 minutes)
go run main.go -config config.yaml

# Run with custom duration (30 seconds)
go run main.go -config config.yaml -duration 30s
```

## Output
```
-----------------------------------------------------------------------------------------------------------------------
Username Resource Type Total Calls Calls/sec 4xx Errors 5xx Errors Success Rate
------------------------------------------------------------------------------------------------------------------------
user1 pods 944 31.47 0 0 100.00%
user1 services 468 15.60 0 0 100.00%
user1 configmaps 470 15.67 0 0 100.00%
```


## Metrics Collected

- Total number of API calls per resource
- Calls per second
- 4xx errors count
- 5xx errors count
- Success rate percentage

## Architecture

The tool consists of several key components:

- **Resource Discovery**: Automatically discovers available API resources
- **Worker**: Manages concurrent API calls for each user
- **Stats Collector**: Tracks and aggregates performance metrics
- **Summary Generator**: Produces formatted test results

## Error Handling

The tool handles various error scenarios:

- API server connection failures
- Authentication errors
- Rate limiting and throttling
- Resource-specific errors

## Best Practices

1. Start with lower concurrency values and gradually increase
2. Monitor cluster resources during testing
3. Use appropriate QPS and burst settings
4. Consider cluster size and capacity when setting test parameters

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
