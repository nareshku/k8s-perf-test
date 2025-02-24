package summary

import (
	"fmt"
	"strings"
	"text/tabwriter"
	"os"
	"github.com/nareshku/k8s-perf-test/types"
)

func PrintSummary(allStats map[string][]*types.CallStats) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	
	fmt.Fprintln(w, "\nPerformance Test Results:")
	fmt.Fprintln(w, strings.Repeat("-", 120))
	fmt.Fprintf(w, "%-30s\t%-25s\t%12s\t%10s\t%10s\t%10s\t%12s\n",
		"Username",
		"Resource Type",
		"Total Calls",
		"Calls/sec",
		"4xx Errors",
		"5xx Errors",
		"Success Rate")
	fmt.Fprintln(w, strings.Repeat("-", 120))

	for username, stats := range allStats {
		for _, stat := range stats {
			successRate := float64(stat.TotalCalls-stat.Errors4xx-stat.Errors5xx) / float64(stat.TotalCalls) * 100
			fmt.Fprintf(w, "%-30s\t%-25s\t%12d\t%10.2f\t%10d\t%10d\t%12.2f%%\n",
				username,
				stat.ResourceType,
				stat.TotalCalls,
				stat.CallsPerSec,
				stat.Errors4xx,
				stat.Errors5xx,
				successRate,
			)
		}
	}
	
	w.Flush()
} 