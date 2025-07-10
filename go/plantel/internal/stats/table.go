package stats

import (
	"fmt"
	"plantel/internal/player"
	"plantel/internal/situation"
	"slices"
)

type FrequencyTable struct {
	matrix       map[uint16]map[situation.Situation]int
	rowTotals    map[uint16]int
	columnTotals map[situation.Situation]int
	total        int
	contractKeys []uint16
}

func NewFrequencyTable(pl []player.Player, thisSeason bool) *FrequencyTable {
	ft := &FrequencyTable{
		matrix:       make(map[uint16]map[situation.Situation]int),
		rowTotals:    make(map[uint16]int),
		columnTotals: make(map[situation.Situation]int),
		total:        0,
	}

	seen := make(map[uint16]bool)
	for _, p := range pl {
		sit := p.Situation
		if !thisSeason {
			if p.NextSeason != "" {
				sit = p.NextSeason
			}
		}

		if ft.matrix[p.ContractEnds] == nil {
			ft.matrix[p.ContractEnds] = make(map[situation.Situation]int)
		}
		ft.matrix[p.ContractEnds][situation.FromString(sit)]++
		ft.columnTotals[situation.FromString(sit)]++
		ft.rowTotals[p.ContractEnds]++
		ft.total++
		if !seen[p.ContractEnds] {
			seen[p.ContractEnds] = true
			ft.contractKeys = append(ft.contractKeys, p.ContractEnds)
		}
	}
	slices.Sort(ft.contractKeys)
	return ft
}

func (ft *FrequencyTable) PrintTable() {
	retained := []situation.Situation{}
	outgoing := []situation.Situation{}

	// Collect and sort situations based on enum values
	for _, sit := range situation.Labels() {
		if sit == situation.Invalid || sit == situation.Divider {
			continue
		}
		if sit < situation.Divider {
			retained = append(retained, sit)
		} else {
			outgoing = append(outgoing, sit)
		}
	}

	slices.Sort(retained)
	slices.Sort(outgoing)

	// Header
	fmt.Printf("| %-10s", "Contract")
	for _, sit := range retained {
		fmt.Printf("| %-11s ", situation.Label(sit))
	}
	fmt.Printf("| %-6s ", "Total")
	for _, sit := range outgoing {
		fmt.Printf("| %-11s ", situation.Label(sit))
	}
	fmt.Printf("| %-6s |\n", "Total")

	// Rows
	for _, contract := range ft.contractKeys {
		fmt.Printf("| %-10d", contract)
		rowRetained := 0
		for _, sit := range retained {
			count := ft.matrix[contract][sit]
			fmt.Printf("| %-11d ", count)
			rowRetained += count
		}
		fmt.Printf("| %-6d ", rowRetained)
		for _, sit := range outgoing {
			count := ft.matrix[contract][sit]
			fmt.Printf("| %-11d ", count)
		}
		fmt.Printf("| %-6d |\n", ft.rowTotals[contract])
	}

	// Totals
	fmt.Printf("| %-10s", "Total")
	totalRetained := 0
	for _, sit := range retained {
		count := ft.columnTotals[sit]
		fmt.Printf("| %-11d ", count)
		totalRetained += count
	}
	fmt.Printf("| %-6d ", totalRetained)
	for _, sit := range outgoing {
		count := ft.columnTotals[sit]
		fmt.Printf("| %-11d ", count)
	}
	fmt.Printf("| %-6d |\n", ft.total)
}
