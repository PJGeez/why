package tests

import (
	"bytes"
	"git-from-scratch/cmd"
	"git-from-scratch/internal/index"
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

	// 5. Verify Index Sync (Phase 8.3)
	idx, err := index.ReadIndex(".")
	if err != nil {
		t.Fatalf("Could not read index: %v", err)
	}

	found := false
	for _, entry := range idx.Entries {
		if entry.Path == fname {
			found = true
			// The hash in the index should match the content "v1"
			if entry.Hash == "" {
				t.Errorf("Index entry for %s has empty hash", fname)
			}
			break
		}
	}
	if !found {
		t.Errorf("Index was not synchronized after checkout: %s missing", fname)
	}
}

func TestCatFile(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "why-catfile-*")
	defer os.RemoveAll(tmpDir)

	dotWhy := filepath.Join(tmpDir, ".why")
	os.MkdirAll(filepath.Join(dotWhy, "objects"), 0755)

	fname := "test.txt"
	content := []byte("hello catfile")
	os.WriteFile(filepath.Join(tmpDir, fname), content, 0644)

	oldCwd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldCwd)

	// Hash-object to write it
	hash, err := cmd.HashObject(".", []string{"-w", fname})
	if err != nil {
		t.Fatalf("HashObject failed: %v", err)
	}

	// Capture output of CatFile
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err = cmd.CatFile(".", "-p", hash)
	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("CatFile failed: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if output != string(content) {
		t.Errorf("CatFile returned wrong content: expected %q, got %q", string(content), output)
	}
}

