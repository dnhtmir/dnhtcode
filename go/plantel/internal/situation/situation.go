package situation

import "slices"

type Situation int

const (
	Invalid Situation = iota
	NoPlantel
	Comprado
	Promovido
	Regresso
	Empréstimo

	Divider // Marker between retained and outgoing

	Emprestado
	VaiSair
)

var situationNames = map[Situation]string{
	NoPlantel:  "No plantel",
	Comprado:   "Comprado",
	Promovido:  "Promovido",
	Regresso:   "Regresso",
	Empréstimo: "Empréstimo",
	Emprestado: "Emprestado",
	VaiSair:    "Vai sair",
}

func FromString(name string) Situation {
	for sit, label := range situationNames {
		if label == name {
			return sit
		}
	}
	return Invalid
}

func Labels() []Situation {
	allSituations := []Situation{}
	for sit := range situationNames {
		if sit == Invalid || sit == Divider {
			continue
		}
		allSituations = append(allSituations, sit)
	}
	slices.Sort(allSituations)
	return allSituations
}

func Label(s Situation) string {
	return situationNames[s]
}

func IsRetained(s Situation) bool {
	return s > Invalid && s < Divider
}

func IsOutgoing(s Situation) bool {
	return s > Divider
}
