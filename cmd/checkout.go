package cmd

import (
	"fmt"
	"git-from-scratch/internal/object"
	"git-from-scratch/internal/repo"
)


// GetTreeFromCommit finds the tree hash associated with a given commit hash
func GetTreeFromCommit(repoPath, commitHash string) (string, error) {
	// If the commit hash is empty (initial state), return empty tree
	if commitHash == "" {
		return "", nil
	}

	data, err := object.ReadObject(repoPath, commitHash)
	if err != nil {
		return "", fmt.Errorf("could not read commit %s: %w", commitHash, err)
	}

	commit, err := object.ParseCommit(data)
	if err != nil {
		return "", fmt.Errorf("could not parse commit %s: %w", commitHash, err)
	}

	return commit.Tree, nil
}


func Checkout(repoPath, target string) error {
	r := &repo.Repository{WorkTree: repoPath, GitDir: repoPath + "/.why"}

	// Resolve Target
	commitHash, isBranch, err := r.ResolveTarget(target)
	if err != nil {
		return err
	}

	// Get Tree Hash
	treeHash, err := GetTreeFromCommit(repoPath, commitHash)
	if err != nil {
		return err
	}

	fmt.Printf("Resolved %s to commit %s (tree: %s)\n", target, commitHash, treeHash)

	//Update HEAD
	err = r.UpdateHead(target, isBranch)
	if err != nil {
		return err
	}

	fmt.Printf("Checked out %s\n", target)
	return nil
}
