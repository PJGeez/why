package cmd

import (
	"fmt"
	"git-from-scratch/internal/index"
	"git-from-scratch/internal/object"
	"git-from-scratch/internal/repo"
	"os"
	"path/filepath"
)


// cleanupWorkingDir removes all files and directories in the worktree EXCEPT .why and project essentials
func cleanupWorkingDir(repoPath string) error {
	entries, err := os.ReadDir(repoPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		name := entry.Name()
		// Protect internal metadata and the tool itself
		if name == ".why" || name == ".git" || name == ".idea" || name == "why" || name == "go.mod" || name == "go.sum" {
			continue
		}

		path := filepath.Join(repoPath, name)
		var removeErr error
		if entry.IsDir() {
			removeErr = os.RemoveAll(path)
		} else {
			removeErr = os.Remove(path)
		}

		if removeErr != nil {
			return fmt.Errorf("could not remove %s: %w", path, removeErr)
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

	// 1. Strip the object header from the tree data
	obj, err := object.ParseObject(data)
	if err != nil {
		return err
	}

	// 2. Parse the content into a Tree struct
	tree, err := object.ParseTree(obj.Content)
	if err != nil {
		return err
	}

	for _, entry := range tree.Entries {
		fullPath := filepath.Join(currentPath, entry.Name)

		if entry.Mode == "040000" { // Directory
			if err := os.MkdirAll(fullPath, 0755); err != nil {
				return err
			}
			if err := UnpackTree(repoPath, entry.Hash, fullPath); err != nil {
				return err
			}
		} else { // File (Blob)
			blobData, err := object.ReadObject(repoPath, entry.Hash)
			if err != nil {
				return err
			}
			// 3. Strip blob header before writing the file
			blobObj, err := object.ParseObject(blobData)
			if err != nil {
				return err
			}
			if err := os.WriteFile(fullPath, blobObj.Content, 0644); err != nil {
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

	// 4. Strip commit header
	obj, err := object.ParseObject(data)
	if err != nil {
		return "", err
	}

	commit, err := object.ParseCommit(obj.Content)
	if err != nil {
		return "", fmt.Errorf("could not parse commit %s: %w", commitHash, err)
	}

	return commit.Tree, nil
}


func Checkout(repoPath, target string) error {
	r := &repo.Repository{WorkTree: repoPath, GitDir: repoPath + "/.why"}

	// 1. Resolve Target (branch or hash)
	commitHash, isBranch, err := r.ResolveTarget(target)
	if err != nil {
		return err
	}

	// 2. Get Tree Hash from the resolved commit
	treeHash, err := GetTreeFromCommit(repoPath, commitHash)
	if err != nil {
		return err
	}

	// 3. Cleanup Working Directory (Phase 7.4)
	if err := cleanupWorkingDir(repoPath); err != nil {
		return fmt.Errorf("cleanup failed: %w", err)
	}

	// 4. Restore Files from the Tree (Phase 7.3)
	newIdx := &index.Index{Entries: []index.IndexEntry{}}
	if treeHash != "" {
		if err := UnpackTree(repoPath, treeHash, repoPath); err != nil {
			return fmt.Errorf("unpack failed: %w", err)
		}
		// 5. Sync Index (Phase 8.3)
		if err := newIdx.FromTree(repoPath, treeHash, ""); err != nil {
			return fmt.Errorf("index sync failed: %w", err)
		}
	}

	if err := index.WriteIndex(repoPath, newIdx); err != nil {
		return fmt.Errorf("could not write new index: %w", err)
	}

	// 6. Update HEAD (Phase 7.5)
	if err := r.UpdateHead(target, isBranch); err != nil {
		return err
	}

	fmt.Printf("Checked out %s (commit %s)\n", target, commitHash)
	return nil
}
