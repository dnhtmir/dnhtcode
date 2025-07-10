package main

import (
	"fmt"
	"log"
	"os"
	"plantel/internal/player"
	"plantel/internal/stats"
)

func main() {
	fp := os.Args[1]
	ps, err := player.LoadPlayers(fp)
	if err != nil {
		log.Fatalf("Error: %v, Exiting", err)
	}
	for _, p := range ps {
		fmt.Println(p)
	}

	numbering := stats.PlayerNumbering(ps)
	fmt.Println(numbering)

	ft := stats.NewFrequencyTable(ps, true)
	ft.PrintTable()

	ft2 := stats.NewFrequencyTable(ps, false)
	ft2.PrintTable()

	stats.CreateAgeMetrics(ps)
}
