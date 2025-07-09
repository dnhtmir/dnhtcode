package player

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"plantel/internal/positions"
)

var (
	allowedSituations = map[string]bool{
		"No plantel": true,
		"Comprado":   true,
		"Promovido":  true,
		"Empréstimo": true,
		"Regresso":   true,
		"Emprestado": true,
		"Vai sair":   true,
	}
)

type National struct {
	IsRepresenting bool   `json:"is_representing"`
	Team           string `json:"team"`
}

// RawPlayer holds unmarshaled JSON as-is
type RawPlayer struct {
	Position       string   `json:"position"`
	Number         int      `json:"number"`
	Name           string   `json:"name"`
	Birthdate      string   `json:"birthdate"`
	Situation      string   `json:"situation"`
	Country        string   `json:"country"`
	PassPercentage int      `json:"pass_percentage"`
	ContractEnds   int      `json:"contract_ends"`
	FromAcademy    bool     `json:"academy"`
	Injured        bool     `json:"injured"`
	NT             National `json:"national"`
	Suspended      bool     `json:"suspended"`
	NextSeason     string   `json:"next_season"`
	Joined         string   `json:"joined"`
}

// Player is the validated and enriched model
type Player struct {
	Joined         time.Time
	Position       string
	Number         uint8
	Name           string
	Birthdate      time.Time
	Age            uint8
	Situation      string
	Country        string
	PassPercentage uint8
	ContractEnds   uint16
}

func LoadPlayers(fp string) ([]Player, error) {
	file, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	raw := make([]RawPlayer, 0)

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&raw); err != nil {
		return nil, err
	}

	var validated []Player
	now := time.Now()

	for _, rp := range raw {
		// Validate position
		if positions.FromString(strings.TrimSpace(rp.Position)) == positions.Invalid {
			return nil, fmt.Errorf("invalid position: %s", rp.Position)
		}

		// Validate situation
		if !allowedSituations[strings.TrimSpace(rp.Situation)] {
			return nil, fmt.Errorf("invalid situation: %s", rp.Situation)
		}

		// Validate number
		if rp.Number < 0 || rp.Number > 99 {
			return nil, fmt.Errorf("invalid number: %d", rp.Number)
		}

		// Validate passPercentage
		if rp.PassPercentage < 0 || rp.PassPercentage > 100 {
			return nil, fmt.Errorf("invalid pass percentage: %d", rp.PassPercentage)
		}

		// Validate contractEnds
		if rp.ContractEnds < now.Year() {
			return nil, fmt.Errorf("contract end is in the past: %d", rp.ContractEnds)
		}

		// Parse birthdate
		birthDate, err := parseDate(rp.Birthdate, "birthdate", "02/01/2006", now)
		if err != nil {
			return nil, err
		}

		// Parse joined date
		joined, err := parseDate(rp.Joined, "joined date", "2006-01-02", now)
		if err != nil {
			return nil, err
		}

		// Calculate age
		age := uint8(0)
		if !birthDate.IsZero() {
			age = uint8(now.Year() - birthDate.Year())
			if now.YearDay() < birthDate.YearDay() {
				age--
			}
		}

		// Append validated player
		validated = append(validated, Player{
			Joined:         joined,
			Position:       rp.Position,
			Number:         uint8(rp.Number),
			Name:           rp.Name,
			Birthdate:      birthDate,
			Age:            age,
			Situation:      rp.Situation,
			Country:        rp.Country,
			PassPercentage: uint8(rp.PassPercentage),
			ContractEnds:   uint16(rp.ContractEnds),
		})
	}

	// Sort first by Position, then by joined date
	sort.Slice(validated, func(i, j int) bool {
		pos_i := positions.FromString(validated[i].Position)
		pos_j := positions.FromString(validated[j].Position)
		if pos_i != pos_j {
			return pos_i < pos_j
		}
		return validated[i].Joined.Before(validated[j].Joined)
	})

	return validated, nil
}

func parseDate(dateStr, label, format string, now time.Time) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, nil
	}

	parsed, err := time.Parse(format, dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid %s format: %s", label, dateStr)
	}
	if parsed.After(now) {
		return time.Time{}, fmt.Errorf("%s cannot be in the future", label)
	}

	return parsed, nil
}
