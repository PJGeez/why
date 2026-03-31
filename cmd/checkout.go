package cmd

import (
	"fmt"
	"git-from-scratch/internal/object"
	"git-from-scratch/internal/repo"
	"os"
	"path/filepath"
)

// cleanupWorkingDir removes all files and directories in the worktree EXCEPT .why
func cleanupWorkingDir(repoPath string) error {
	entries, err := os.ReadDir(repoPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		name := entry.Name()
		// CRITICAL: Do not delete internal metadata, Git history, or the binary!
		if name == ".why" || name == ".git" || name == ".idea" || name == "why" || name == "go.mod" {
			continue
		}

		path := filepath.Join(repoPath, name)
		if entry.IsDir() {
			err = os.RemoveAll(path)
		} else {
			err = os.Remove(path)
		}
		if err != nil {
			return fmt.Errorf("could not remove %s: %w", path, err)
		}
	}
	return nil
}

// UnpackTree recursively restores files and folders from a tree object
func UnpackTree(repoPath, treeHash, currentPath string) error {
	if treeHash == "" {
		return nil
	}

	data, err := object.ReadObject(repoPath, treeHash)
	if err != nil {
		return err
	}

	tree, err := object.ParseTree(data)
	if err != nil {
		return err
	}

	for _, entry := range tree.Entries {
		fullPath := filepath.Join(currentPath, entry.Name)

		if entry.Mode == "040000" { 
			if err := os.MkdirAll(fullPath, 0755); err != nil {
				return err
			}
			if err := UnpackTree(repoPath, entry.Hash, fullPath); err != nil {
				return err
			}
		} else { 
			blobData, err := object.ReadObject(repoPath, entry.Hash)
			if err != nil {
				return err
			}
			if err := os.WriteFile(fullPath, blobData, 0644); err != nil {
				return err
			}
		}
	}
	return nil
}


// GetTreeFromCommit finds the tree hash associated with a given commit hash
func GetTreeFromCommit(repoPath, commitHash string) (string, error) {
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

	//Resolve Target
	commitHash, isBranch, err := r.ResolveTarget(target)
	if err != nil {
		return err
	}

	//Get Tree Hash
	treeHash, err := GetTreeFromCommit(repoPath, commitHash)
	if err != nil {
		return err
	}

	//Cleanup Working Directory
	if err := cleanupWorkingDir(repoPath); err != nil {
		return fmt.Errorf("cleanup failed: %w", err)
	}

	//Restore Files
	if treeHash != "" {
		if err := UnpackTree(repoPath, treeHash, repoPath); err != nil {
			return fmt.Errorf("unpack failed: %w", err)
		}
	}

	//Update HEAD
	if err := r.UpdateHead(target, isBranch); err != nil {
		return err
	}

	fmt.Printf("Checked out %s (commit %s)\n", target, commitHash)
	return nil
}
