package stats

import (
	"fmt"
	"plantel/internal/player"
	"slices"
)

// AgeMetric represents one row in the output table
type AgeMetric struct {
	Age               uint8
	Frequency         uint8
	Percentage        float64
	CumulativePercent float64
}

// CreateAgeMetrics builds and prints age distribution table
func CreateAgeMetrics(pl []player.Player) {
	counts := make(map[uint8]uint8)
	total := len(pl)
	cumulative := uint8(0)

	// Count age occurrences
	for _, p := range pl {
		counts[p.Age]++
	}

	// Sorted unique ages
	ages := make([]uint8, 0, len(counts))
	for age := range counts {
		ages = append(ages, age)
	}
	slices.Sort(ages)

	// Build metrics
	metrics := make([]AgeMetric, 0, len(ages))
	for _, age := range ages {
		freq := counts[age]
		cumulative += freq
		metric := AgeMetric{
			Age:               age,
			Frequency:         freq,
			Percentage:        float64(freq) / float64(total) * 100,
			CumulativePercent: float64(cumulative) / float64(total) * 100,
		}
		metrics = append(metrics, metric)
	}

	// Print table
	fmt.Printf("| %-5s | %-9s | %-10s | %-12s |\n", "Age", "Frequency", "Percentage", "Cumulative %")
	fmt.Println("|-------|-----------|------------|--------------|")
	for _, m := range metrics {
		fmt.Printf("| %-5d | %-9d | %9.2f%% | %11.2f%% |\n",
			m.Age, m.Frequency, m.Percentage, m.CumulativePercent)
	}

}
