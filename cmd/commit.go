package cmd

import (
	"fmt"
	"git-from-scratch/internal/object"
)

func Commit(repoPath, treeHash, parentHash, message string) error {
	commit := object.Commit{
		Tree: treeHash,
		Parent: parentHash,
		Author: "Prajwal <crazyshit.dev@gmail.com>",
		Message: message,
	}

	data := commit.Serialize()

	hash, err := object.WriteObject(repoPath, data)
	if err != nil {
		return err
	}

	fmt.Println(hash)
	return nil
}