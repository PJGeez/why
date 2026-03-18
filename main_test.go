package main

import (
	"git-from-scratch/internal/object"
	"git-from-scratch/internal/repo"
	"os"
	"path/filepath"
	"testing"
)

func TestGetHeadCommit(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "why-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Mock a .why repo
	dotWhy := filepath.Join(tempDir, ".why")
	os.MkdirAll(filepath.Join(dotWhy, "refs", "heads"), 0755)

	r := &repo.Repository{WorkTree: tempDir, GitDir: dotWhy}

	// 1. Test empty repo
	os.WriteFile(filepath.Join(dotWhy, "HEAD"), []byte("ref: refs/heads/master\n"), 0644)
	hash, err := r.GetHeadCommit()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if hash != "" {
		t.Errorf("Expected empty hash for new repo, got %s", hash)
	}

	// 2. Test valid ref
	expectedHash := "abcdef1234567890abcdef1234567890abcdef12"
	os.WriteFile(filepath.Join(dotWhy, "refs", "heads", "master"), []byte(expectedHash), 0644)
	hash, _ = r.GetHeadCommit()
	if hash != expectedHash {
		t.Errorf("Expected %s, got %s", expectedHash, hash)
	}

	// 3. Test detached HEAD
	detachedHash := "0000000000000000000000000000000000000000"
	os.WriteFile(filepath.Join(dotWhy, "HEAD"), []byte(detachedHash), 0644)
	hash, _ = r.GetHeadCommit()
	if hash != detachedHash {
		t.Errorf("Expected %s, got %s", detachedHash, hash)
	}
}

func TestParseCommit(t *testing.T) {
	raw := "tree 12345\nparent 67890\nauthor Prajwal <prajwal@example.com> 1234567890 +0530\ncommitter Prajwal <prajwal@example.com> 1234567890 +0530\n\nInitial commit message"
	
	c, err := object.ParseCommit([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}

	if c.Tree != "12345" {
		t.Errorf("Expected tree 12345, got %s", c.Tree)
	}
	if c.Parent != "67890" {
		t.Errorf("Expected parent 67890, got %s", c.Parent)
	}
	if c.Message != "Initial commit message" {
		t.Errorf("Expected message 'Initial commit message', got '%s'", c.Message)
	}
}
