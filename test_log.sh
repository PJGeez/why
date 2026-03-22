#!/bin/bash

# --- 1. CLEANUP & INITIALIZATION ---
echo "--- Initializing Repo ---"
rm -rf .why hello.txt test_commit_gen.go
go build -o why main.go
./why init

# --- 2. TEST PHASE 5.7 (Empty Repo) ---
echo "--- Testing Phase 5.7 (Empty Log) ---"
./why log
echo ""

# --- 3. CREATE TEMPORARY GO TOOL TO GENERATE COMMITS ---
# This tool uses your actual 'internal/object' package to create real objects
cat <<EOF > test_commit_gen.go
package main

import (
	"fmt"
	"git-from-scratch/internal/object"
	"os"
)

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: go run test_commit_gen.go <repo_path> <tree_hash> <parent_hash> <message>")
		os.Exit(1)
	}
	
	repo := os.Args[1]
	tree := os.Args[2]
	parent := os.Args[3]
	message := os.Args[4]

	c := object.Commit{
		Tree:    tree,
		Parent:  parent,
		Author:  "Prajwal <prajwal@example.com>",
		Message: message,
	}

	data := c.Serialize()
	hash, err := object.WriteObject(repo, data)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Print(hash)
}
EOF

# --- 4. CREATE COMMIT #1 ---
echo "--- Creating Commit #1 ---"
echo "hello world v1" > hello.txt
./why add hello.txt
TREE1=$(./why write-tree)
COMMIT1=$(go run test_commit_gen.go "." "$TREE1" "" "First commit: added hello.txt")

# Manually update the master branch to point to commit 1
mkdir -p .why/refs/heads
echo "$COMMIT1" > .why/refs/heads/master
echo "Commit 1 created: $COMMIT1"

# --- 5. CREATE COMMIT #2 (Linked to Commit #1) ---
echo "--- Creating Commit #2 ---"
echo "hello world v2" > hello.txt
./why add hello.txt
TREE2=$(./why write-tree)
COMMIT2=$(go run test_commit_gen.go "." "$TREE2" "$COMMIT1" "Second commit: updated hello.txt")

# Update master branch to point to commit 2
echo "$COMMIT2" > .why/refs/heads/master
echo "Commit 2 created: $COMMIT2"

# --- 6. TEST PHASE 5.8 (History Traversal) ---
echo ""
echo "--- THE BIG TEST: RUNNING LOG ---"
./why log

# Cleanup
rm test_commit_gen.go
rm hello.txt
