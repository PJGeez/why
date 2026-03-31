package tests

import (
	"git-from-scratch/cmd"
	"os"
	"path/filepath"
	"testing"
)

func TestFullWorkflow(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "why-integration-*")
	defer os.RemoveAll(tmpDir)

	// 1. Setup the basic directory structure
	dotWhy := filepath.Join(tmpDir, ".why")
	os.MkdirAll(filepath.Join(dotWhy, "objects"), 0755)
	os.MkdirAll(filepath.Join(dotWhy, "refs", "heads"), 0755)
	os.WriteFile(filepath.Join(dotWhy, "HEAD"), []byte("ref: refs/heads/master\n"), 0644)
	
	// 2. Add
	fname := "hello.txt"
	absFname := filepath.Join(tmpDir, fname)
	os.WriteFile(absFname, []byte("v1"), 0644)
	
	// Change to tmpDir for relative path tests
	oldCwd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldCwd)

	if err := cmd.Add(".", []string{fname}); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	// 3. Commit
	if err := cmd.Commit(".", "first commit"); err != nil {
		t.Fatalf("Commit failed: %v", err)
	}
    
    // 4. Verify checkout
    os.WriteFile(fname, []byte("v2"), 0644)
    if err := cmd.Checkout(".", "master"); err != nil {
        t.Fatalf("Checkout failed: %v", err)
    }

	content, _ := os.ReadFile(fname)
	if string(content) != "v1" {
		t.Errorf("Checkout did not restore file correctly: expected v1, got %s", string(content))
	}
}
