package unit_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/titi-byte-dev/gorm-crm/internal/lead"
)

func TestBasicScorer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		value    float64
		status   lead.Status
		wantMin  int // score >= wantMin
		wantMax  int // score <= wantMax
	}{
		{"qualified + high value", 15000, lead.StatusQualified, 70, 100},
		{"new + zero value", 0, lead.StatusNew, -10, 10},
		{"lost penalizes", 5000, lead.StatusLost, -25, 5},
		{"contacted + medium value", 500, lead.StatusContacted, 5, 25},
	}

	scorer := lead.BasicScorer{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			l := &lead.Lead{
				ID:     uuid.New(),
				Value:  tt.value,
				Status: tt.status,
			}
			score := scorer.Score(l)
			if score < tt.wantMin || score > tt.wantMax {
				t.Errorf("Score() = %d, want [%d, %d]", score, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestWeightedScorer_QualifiedBeatsContacted(t *testing.T) {
	t.Parallel()
	scorer := lead.DefaultWeightedScorer()

	qualified := &lead.Lead{Value: 5000, Status: lead.StatusQualified}
	contacted := &lead.Lead{Value: 5000, Status: lead.StatusContacted}

	if scorer.Score(qualified) <= scorer.Score(contacted) {
		t.Error("qualified lead deve pontuar mais que contacted com mesmo valor")
	}
}

// TestScorer_IsSwappable verifica o nucleo do Strategy pattern:
// trocar o scorer nao requer mudar o Service.
func TestScorer_IsSwappable(t *testing.T) {
	t.Parallel()

	l := &lead.Lead{Value: 10000, Status: lead.StatusQualified}

	basic := lead.BasicScorer{}
	weighted := lead.DefaultWeightedScorer()

	// Ambos implementam Scorer — scorer diferente, pontuacao diferente, mesmo contrato
	var _ lead.Scorer = basic
	var _ lead.Scorer = weighted

	if basic.Score(l) == weighted.Score(l) {
		// Nao e um erro obrigatorio, mas demonstra que os algoritmos divergem
		t.Log("nota: BasicScorer e WeightedScorer produziram o mesmo score para este lead")
	}
}
