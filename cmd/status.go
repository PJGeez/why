package cmd

import (
	"fmt"
	"git-from-scratch/internal/index"
	"os"
	"path/filepath"
)

func Status(repoPath string) error {
	idx, err := index.ReadIndex(repoPath)
	if err != nil {
		return fmt.Errorf("could not read index: %w", err)
	}

	// getting all files in the working directory (excluding .why, .git, and binary)
	workingFiles := make(map[string]bool)
	err = filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}


		// Skip directories and the .why/.git folders
		if info.IsDir() {
			if info.Name() == ".why" || info.Name() == ".git" || info.Name() == ".idea" {
				return filepath.SkipDir
			}
			return nil
		}

		//skip the 'why' binary
		if info.Name() == "why" || info.Name() == "go.mod" || info.Name() == "go.sum" {
			return nil
		}

		relPath, _ := filepath.Rel(repoPath, path)
		// Don't track Go source files if you want to keep it focused on your test files,
		// but usually, Git tracks everything.
		workingFiles[relPath] = true
		return nil
	})

	if err != nil {
		return err
	}

	fmt.Println("Changes to be committed:")
	for _, entry := range idx.Entries {
		fmt.Printf("  (staged)   %s\n", entry.Path)
		delete(workingFiles, entry.Path) // Remove from map so we can find untracked files
	}

	fmt.Println("\nUntracked files:")
	if len(workingFiles) == 0 {
		fmt.Println("  (none)")
	} else {
		for path := range workingFiles {
			// Skip .go files to keep output clean during development if you prefer
			fmt.Printf("  %s\n", path)
		}
	}

	return nil
}
