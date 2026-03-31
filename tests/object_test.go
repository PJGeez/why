package tests

import (
	"bytes"
	"git-from-scratch/internal/object"
	"testing"
)

func TestBlobLifecycle(t *testing.T) {
	content := []byte("hello world")
	blob := object.Blob{Content: content}
	data := blob.Serialize()

	obj, err := object.ParseObject(data)
	if err != nil {
		t.Fatalf("Failed to parse object: %v", err)
	}

	if obj.Type != "blob" {
		t.Errorf("Expected type blob, got %s", obj.Type)
	}

	if !bytes.Equal(obj.Content, content) {
		t.Errorf("Content mismatch: expected %s, got %s", string(content), string(obj.Content))
	}
}

func TestCommitParsing(t *testing.T) {
	raw := "tree 12345\nparent 67890\nauthor Prajwal\n\nInitial commit"
	commit, err := object.ParseCommit([]byte(raw))
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if commit.Tree != "12345" {
		t.Errorf("Expected tree 12345, got %s", commit.Tree)
	}
	if commit.Message != "Initial commit" {
		t.Errorf("Expected message 'Initial commit', got '%s'", commit.Message)
	}
}
