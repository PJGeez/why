package tests

import (
	"git-from-scratch/internal/repo"
	"os"
	"strings"
	"testing"
)

func TestRepositoryState(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "why-repo-test-*")
	defer os.RemoveAll(tmpDir)

	r, _ := repo.NewRepository(tmpDir)
	if err := r.Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// 1. Check Default Branch
	branch, err := r.GetCurrentBranch()
	if err != nil || branch != "master" {
		t.Errorf("Expected branch master, got %s (err: %v)", branch, err)
	}

	// 2. Check SetBranchCommit
	hash := "abc1234567890abcdef1234567890abcdef12345"
	if err := r.SetBranchCommit("master", hash); err != nil {
		t.Fatalf("Failed to set branch commit: %v", err)
	}

	// 3. Check ResolveTarget
	resolvedHash, isBranch, err := r.ResolveTarget("master")
	if err != nil || !isBranch || strings.TrimSpace(resolvedHash) != hash {
		t.Errorf("ResolveTarget failed: expected %s, got %s (err: %v)", hash, resolvedHash, err)
	}
}
