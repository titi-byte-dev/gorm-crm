package unit_test

import (
	"testing"

	"github.com/titi-byte-dev/gorm-crm/internal/deal"
)

func TestDealStage_CanTransitionTo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		from     deal.Stage
		to       deal.Stage
		expected bool
	}{
		{"proposal → negotiation", deal.StageProposal, deal.StageNegotiation, true},
		{"proposal → lost", deal.StageProposal, deal.StageLost, true},
		{"proposal → won: inválido", deal.StageProposal, deal.StageWon, false},
		{"negotiation → won", deal.StageNegotiation, deal.StageWon, true},
		{"negotiation → lost", deal.StageNegotiation, deal.StageLost, true},
		{"won é final", deal.StageWon, deal.StageLost, false},
		{"lost é final", deal.StageLost, deal.StageWon, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.from.CanTransitionTo(tt.to)
			if got != tt.expected {
				t.Errorf("(%s).CanTransitionTo(%s) = %v, want %v",
					tt.from, tt.to, got, tt.expected)
			}
		})
	}
}

func TestDealStage_IsClosed(t *testing.T) {
	t.Parallel()

	tests := []struct {
		stage  deal.Stage
		closed bool
	}{
		{deal.StageProposal, false},
		{deal.StageNegotiation, false},
		{deal.StageWon, true},
		{deal.StageLost, true},
	}

	for _, tt := range tests {
		t.Run(string(tt.stage), func(t *testing.T) {
			t.Parallel()
			if got := tt.stage.IsClosed(); got != tt.closed {
				t.Errorf("Stage(%q).IsClosed() = %v, want %v", tt.stage, got, tt.closed)
			}
		})
	}
}
