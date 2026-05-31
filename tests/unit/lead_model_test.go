package unit_test

import (
	"testing"

	"github.com/titi-byte-dev/gorm-crm/internal/lead"
)

// TestLeadStatus_CanTransitionTo usa table-driven tests — o idioma Go para
// testar múltiplos cenários de forma legível e sem repetição.
func TestLeadStatus_CanTransitionTo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		from     lead.Status
		to       lead.Status
		expected bool
	}{
		{"new can go to contacted", lead.StatusNew, lead.StatusContacted, true},
		{"new can go to lost", lead.StatusNew, lead.StatusLost, true},
		{"new cannot skip to qualified", lead.StatusNew, lead.StatusQualified, false},
		{"contacted can go to qualified", lead.StatusContacted, lead.StatusQualified, true},
		{"contacted can go to lost", lead.StatusContacted, lead.StatusLost, true},
		{"qualified can only go to lost", lead.StatusQualified, lead.StatusLost, true},
		{"qualified cannot go back to new", lead.StatusQualified, lead.StatusNew, false},
		{"lost is a final state", lead.StatusLost, lead.StatusNew, false},
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

func TestLeadStatus_IsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		status lead.Status
		valid  bool
	}{
		{lead.StatusNew, true},
		{lead.StatusContacted, true},
		{lead.StatusQualified, true},
		{lead.StatusLost, true},
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			t.Parallel()
			if got := tt.status.IsValid(); got != tt.valid {
				t.Errorf("Status(%q).IsValid() = %v, want %v", tt.status, got, tt.valid)
			}
		})
	}
}
