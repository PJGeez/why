package cmd

import (
	"fmt"
	"git-from-scratch/internal/object"
)

func CatFile(repoPath, hash string) error {
	data, err := object.ReadObject(repoPath, hash)
	if err != nil {
		return err
	}

	parsed, err := object.ParseObject(data)
	if err != nil {
		return err
	}

	fmt.Printf("Type: %s\n", parsed.Type)
	fmt.Printf("Size: %d\n", parsed.Size)
	fmt.Println("Content:")
	fmt.Println(string(parsed.Content))

	return nil
}