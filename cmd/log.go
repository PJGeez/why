package cmd

import (
	"fmt"
	"git-from-scratch/internal/object"
	"git-from-scratch/internal/repo"
)

func Log(repoPath string) error {
	r := &repo.Repository{WorkTree: repoPath, GitDir: repoPath + "/.why"}

	currentHash, err := r.GetHeadCommit()
	if err != nil {
		return fmt.Errorf("could not read HEAD: %w", err)
	}

	if currentHash == "" {
		fmt.Println("No commits yet.")
		return nil
	}

	for currentHash != "" {
		data, err := object.ReadObject(repoPath, currentHash)
		if err != nil {
			return fmt.Errorf("error reading commit %s: %w", currentHash, err)
		}

		commit, err := object.ParseCommit(data)
		if err != nil {
			return fmt.Errorf("error parsing commit %s: %w", currentHash, err)
		}

		fmt.Printf("commit %s\n", currentHash)
		// Note: we could further parse the author line to show just the name
		fmt.Printf("Author: %s\n", commit.Author)
		fmt.Printf("\n    %s\n\n", commit.Message)

		currentHash = commit.Parent
	}

	return nil
}
