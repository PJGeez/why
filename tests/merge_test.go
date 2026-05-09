package tests

import (
	"git-from-scratch/internal/repo"
	"testing"
)

func TestMergeDecisionMatrix(t *testing.T) {
	// Setup dummy hashes
	baseHash := "BASE"
	oursHash := "OURS"
	theirsHash := "THEIRS"

	tests := []struct {
		name     string
		base     map[string]string
		ours     map[string]string
		theirs   map[string]string
		expected repo.MergeAction
		expHash  string
	}{
		{
			name:     "Case A: No changes anywhere",
			base:     map[string]string{"file.txt": baseHash},
			ours:     map[string]string{"file.txt": baseHash},
			theirs:   map[string]string{"file.txt": baseHash},
			expected: repo.ActionKeep,
			expHash:  baseHash,
		},
		{
			name:     "Case B: Change only in Ours",
			base:     map[string]string{"file.txt": baseHash},
			ours:     map[string]string{"file.txt": oursHash},
			theirs:   map[string]string{"file.txt": baseHash},
			expected: repo.ActionTakeOurs,
			expHash:  oursHash,
		},
		{
			name:     "Case C: Change only in Theirs",
			base:     map[string]string{"file.txt": baseHash},
			ours:     map[string]string{"file.txt": baseHash},
			theirs:   map[string]string{"file.txt": theirsHash},
			expected: repo.ActionTakeTheirs,
			expHash:  theirsHash,
		},
		{
			name:     "Case D: Same change in both branches",
			base:     map[string]string{"file.txt": baseHash},
			ours:     map[string]string{"file.txt": oursHash},
			theirs:   map[string]string{"file.txt": oursHash},
			expected: repo.ActionKeep,
			expHash:  oursHash,
		},
		{
			name:     "Case E: True Conflict (Divergent)",
			base:     map[string]string{"file.txt": baseHash},
			ours:     map[string]string{"file.txt": oursHash},
			theirs:   map[string]string{"file.txt": theirsHash},
			expected: repo.ActionConflict,
			expHash:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan := repo.ComputeMergePlan(tt.base, tt.ours, tt.theirs)
			decision := plan["file.txt"]
			
			if decision.Action != tt.expected {
				t.Errorf("%s: expected action %v, got %v", tt.name, tt.expected, decision.Action)
			}
			if decision.Hash != tt.expHash {
				t.Errorf("%s: expected hash %s, got %s", tt.name, tt.expHash, decision.Hash)
			}
		})
	}
}
