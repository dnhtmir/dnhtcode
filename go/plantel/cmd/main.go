package main

import (
	"fmt"
	"log"
	"os"
	"plantel/internal/player"
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
}
