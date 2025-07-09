package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"plantel/internal/player"
	"strconv"
	"strings"
	"time"
)

func readCSV(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

func main() {
	csvfp := os.Args[1]
	jsonfp := os.Args[2]

	data, err := readCSV(csvfp)
	if err != nil {
		log.Fatalf("failed to read csv: %v", err)
	}

	err = archivePlayersFile(jsonfp)
	if err != nil {
		log.Fatalf("failed to archive players file: %v", err)
	}

	err = createNewPlayersFile(jsonfp, data)
	if err != nil {
		log.Fatalf("failed to create new players file: %v", err)
	}
}

func archivePlayersFile(jsonfp string) error {
	_, err := os.Stat(jsonfp)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}

	// Step 1: Format timestamp
	timestamp := time.Now().Format("20060102_150405")
	newFileName := fmt.Sprintf("players_%s.json", timestamp)

	// Step 2: Create "old" directory if it doesn't exist
	f := filepath.Join(filepath.Dir(jsonfp), "old")
	if err := os.MkdirAll(f, os.ModePerm); err != nil {
		return err
	}

	// Step 3: Move and rename the original file
	newPath := filepath.Join(f, newFileName)
	if err := os.Rename(jsonfp, newPath); err != nil {
		return err
	}

	return nil
}

func createNewPlayersFile(jsonfp string, pl [][]string) error {
	j := make([]*player.RawPlayer, 0)
	for _, p := range pl {
		number, err := parseNumber(p[1], 1, 99)
		if err != nil {
			return err
		}
		pp, err := parseNumber(strings.ReplaceAll(p[7], "%", ""), 0, 100)
		if err != nil {
			return err
		}
		ce, err := parseNumber(p[8], 0, 2100)
		if err != nil {
			return err
		}
		a, err := getBool(p[9])
		if err != nil {
			return err
		}
		i, err := getBool(p[10])
		if err != nil {
			return err
		}
		ir, err := getBool(p[11])
		if err != nil {
			return err
		}
		s, err := getBool(p[13])
		if err != nil {
			return err
		}

		j = append(j, &player.RawPlayer{
			Position:       p[0],
			Number:         number,
			Name:           p[2],
			Birthdate:      p[3],
			Situation:      p[5],
			Country:        p[6],
			PassPercentage: pp,
			ContractEnds:   ce,
			FromAcademy:    a,
			Injured:        i,
			NT: player.National{
				IsRepresenting: ir,
				Team:           p[12],
			},
			Suspended:  s,
			NextSeason: p[14],
			Joined:     p[15],
		})
	}

	jd, err := json.MarshalIndent(j, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(jsonfp, jd, 0644)
	if err != nil {
		return err
	}

	return nil
}

func parseNumber(s string, l, h int) (int, error) {
	if s == "" {
		return 0, nil
	}
	number, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	if number < l || number > h {
		return number, fmt.Errorf("number '%d' is out of the %d-%d range", number, l, h)
	}
	return number, nil
}

func getBool(s string) (bool, error) {
	if s == "X" {
		return true, nil
	}

	if s == "" {
		return false, nil
	}

	return false, fmt.Errorf("failed to get bool from: %s", s)
}
