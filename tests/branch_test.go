package tests

import (
	"git-from-scratch/internal/repo"
	"os"
	"testing"
)

func TestBranching(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "why-branch-test-*")
	defer os.RemoveAll(tmpDir)

	r, _ := repo.NewRepository(tmpDir)
	r.Init()

	// 1. Initial State
	branches, current, _ := r.ListBranches()
	if len(branches) != 1 || branches[0] != "master" {
		t.Errorf("Expected only master branch, got %v", branches)
	}
	if current != "master" {
		t.Errorf("Expected current branch master, got %s", current)
	}

	// 2. Create New Branch
	hash := "1234567890123456789012345678901234567890"
	err := r.CreateBranch("feature", hash)
	if err != nil {
		t.Fatalf("Failed to create branch: %v", err)
	}

	// 3. Verify List
	branches, _, _ = r.ListBranches()
	found := false
	for _, b := range branches {
		if b == "feature" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("New branch 'feature' not found in listing")
	}

	// 4. Verify Branch Content
	resolvedHash, isBranch, _ := r.ResolveTarget("feature")
	if !isBranch || resolvedHash != hash {
		t.Errorf("Branch 'feature' points to wrong hash: expected %s, got %s", hash, resolvedHash)
	}
}
