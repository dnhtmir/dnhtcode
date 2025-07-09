package stats

import "plantel/internal/player"

func PlayerNumbering(pl []player.Player) map[uint8]string {
	m := make(map[uint8]string)
	for _, p := range pl {
		if p.Number != 0 {
			m[p.Number] = p.Name
		}
	}
	return m
}
