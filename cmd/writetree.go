package cmd

import (
	"fmt"
	"git-from-scratch/internal/index"
	"git-from-scratch/internal/object"
	"sort"
)

func WriteTree(repoPath string) error {
	idx, err := index.ReadIndex(repoPath)
	if err != nil {
		return fmt.Errorf("could not read index: %w", err)
	}

	if len(idx.Entries) == 0 {
		return fmt.Errorf("nothing to commit work-tree is clean")
	}

	var entries []object.TreeEntry

	for _, idxEntry := range idx.Entries {
		entry := object.TreeEntry{
			Mode: idxEntry.Mode,
			Name: idxEntry.Path,
			Hash: idxEntry.Hash,
		}
		entries = append(entries, entry)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})

	tree := object.Tree{Entries: entries}
	treeData, err := tree.Serialize()
	if err != nil {
		return fmt.Errorf("error serializing tree: %w", err)
	}

	treeHash, err := object.WriteObject(repoPath, treeData)
	if err != nil {
		return fmt.Errorf("error writing tree object: %w", err)
	}
	fmt.Println(treeHash)
	return nil
}