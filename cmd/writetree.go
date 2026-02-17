package cmd

import (
	"fmt"
	"os"
	"git-from-scratch/internal/object"
)

func WriteTree(repoPath string) error {
	files, err := os.ReadDir(repoPath)
	if err!= nil {
		return err
	}

	var entries []object.TreeEntry

	for _, file := range files {
		if file.IsDir() || file.Name() == ".git" {
			continue
		}

		content, err := os.ReadFile(file.Name())
		if err != nil {
			return err
		}

		blob := object.Blob{Content: content}
		serialized := blob.Serialize()

		hash, err := object.WriteObject(repoPath, serialized)
		if err != nil {
			return err
		}

		//hard coded values
		entry := object.TreeEntry {
			Mode: "100644",
			Name: file.Name(),
			Hash: hash,
		}

		entries = append(entries, entry)
	}

	tree := object.Tree{Entries: entries}
	treeData, err := tree.Serialize()
	if err != nil {
		return err
	}

	treeHash, err := object.WriteObject(repoPath, treeData)
	if err != nil {
		return err
	}
	fmt.Println(treeHash)
	return nil
}