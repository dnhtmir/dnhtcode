package positions

type Position int

const (
	Invalid Position = iota
	GuardaRedes
	LateralDireito
	LateralEsquerdo
	DefesaCentral
	MedioDefensivo
	MedioCentro
	MedioOfensivo
	Extremo
	Avancado
)

var positionNames = map[Position]string{
	GuardaRedes:     "Guarda Redes",
	LateralDireito:  "Lateral Direito",
	LateralEsquerdo: "Lateral Esquerdo",
	DefesaCentral:   "Defesa Central",
	MedioDefensivo:  "Médio Defensivo",
	MedioCentro:     "Médio Centro",
	MedioOfensivo:   "Médio Ofensivo",
	Extremo:         "Extremo",
	Avancado:        "Avançado",
}

func FromString(name string) Position {
	for pos, label := range positionNames {
		if label == name {
			return pos
		}
	}
	return Invalid
}
