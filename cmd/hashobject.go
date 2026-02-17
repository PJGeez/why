package cmd

import (
	"fmt"
	"os"
	"git-from-scratch/internal/object"
)

func HashObject(repoPath, filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	blob := object.Blob{Content : content}
	serialized := blob.Serialize()

	hash, err := object.WriteObject(repoPath, serialized)
	if err != nil {
		return err
	}

	fmt.Println(hash)
	return nil
}