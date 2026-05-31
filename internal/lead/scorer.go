package lead

// Scorer é a interface Strategy para calcular a pontuação de um lead.
//
// Strategy pattern: define uma família de algoritmos, encapsula cada um,
// e torna-os intercambiáveis. O Service não sabe qual algoritmo usa —
// só sabe que tem um Scorer.
type Scorer interface {
	Score(lead *Lead) int
}

// BasicScorer implementa uma pontuação simples baseada em regras fixas.
// É o Scorer padrão — sem configuração, sem dependências.
type BasicScorer struct{}

func (BasicScorer) Score(lead *Lead) int {
	score := 0

	switch lead.Status {
	case StatusContacted:
		score += 10
	case StatusQualified:
		score += 30
	case StatusLost:
		score -= 20
	}

	switch {
	case lead.Value >= 10000:
		score += 50
	case lead.Value >= 1000:
		score += 20
	case lead.Value >= 100:
		score += 5
	}

	return score
}

// WeightedScorer permite configurar o peso de cada critério.
// Útil quando diferentes equipas de vendas usam escalas diferentes.
type WeightedScorer struct {
	ValueWeight  float64 // peso aplicado ao valor normalizado
	StatusWeight float64 // peso aplicado ao status
}

// DefaultWeightedScorer devolve pesos equilibrados para uso geral.
func DefaultWeightedScorer() WeightedScorer {
	return WeightedScorer{ValueWeight: 0.6, StatusWeight: 0.4}
}

func (s WeightedScorer) Score(lead *Lead) int {
	valueScore := s.scoreValue(lead.Value)
	statusScore := s.scoreStatus(lead.Status)

	return int(float64(valueScore)*s.ValueWeight + float64(statusScore)*s.StatusWeight)
}

func (s WeightedScorer) scoreValue(value float64) int {
	switch {
	case value >= 50000:
		return 100
	case value >= 10000:
		return 70
	case value >= 1000:
		return 40
	default:
		return 10
	}
}

func (s WeightedScorer) scoreStatus(status Status) int {
	scores := map[Status]int{
		StatusNew:       10,
		StatusContacted: 30,
		StatusQualified: 80,
		StatusLost:      0,
	}
	return scores[status]
}
