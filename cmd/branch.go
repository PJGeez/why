package cmd

import (
	"fmt"
	"git-from-scratch/internal/repo"
)

func Branch(repoPath string, args []string) error {
	r := &repo.Repository{WorkTree: repoPath, GitDir: repoPath + "/.why"}

	// 1. If no args, LIST branches (Phase 8.2)
	if len(args) == 0 {
		branches, current, err := r.ListBranches()
		if err != nil {
			return err
		}
		for _, b := range branches {
			prefix := "  "
			if b == current {
				prefix = "* "
			}
			fmt.Printf("%s%s\n", prefix, b)
		}
		if current == "" {
			fmt.Println("(Detached HEAD)")
		}
		return nil
	}

	// 2. If arg provided, CREATE branch (Phase 8.1)
	newBranch := args[0]
	currentCommit, err := r.GetHeadCommit()
	if err != nil {
		return err
	}
	if currentCommit == "" {
		return fmt.Errorf("fatal: not a valid object name: 'master'")
	}

	err = r.CreateBranch(newBranch, currentCommit)
	if err != nil {
		return err
	}

	fmt.Printf("Created branch %s at %s\n", newBranch, currentCommit[:7])
	return nil
}
