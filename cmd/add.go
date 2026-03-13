package cmd

import (
	"fmt"
	"git-from-scratch/internal/index"
	"git-from-scratch/internal/object"
	"os"
)

func Add(repoPath string, files []string) error {
	idx, err := index.ReadIndex(repoPath)
	if err != nil {
		return fmt.Errorf("could not read index: %w", err)
	}

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("error reading %s: %w", file, err)
		}

		blob := object.Blob{Content: data}
		hash, err := object.WriteObject(repoPath, blob.Serialize())
		if err != nil {
			return fmt.Errorf("error writing blob for %s: %w", file, err)
		}

		// Update or add entry in index
		found := false
		for i, entry := range idx.Entries {
			if entry.Path == file {
				idx.Entries[i].Hash = hash
				found = true
				break
			}
		}

		if !found {
			idx.Entries = append(idx.Entries, index.IndexEntry{
				Path: file,
				Hash: hash,
				Mode: "100644",
			})
		}
		fmt.Printf("staged %s\n", file)
	}

	return index.WriteIndex(repoPath, idx)
}
