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

	if parsed.Type == "tree" {
		data := parsed.Content
		i := 0
		for i < len(data) {

			start := i
			for data[i] != ' ' {
				i++
			}
			mode := string(data[start:i])
			i++

			start = i
			for data[i] != 0 {
				i++
			}
			name := string(data[start:i])
			i++

			hashBytes := data[i : i+20]
			hash := fmt.Sprintf("%x", hashBytes)
			i += 20

			fmt.Printf("%s %s %s\n", mode, hash, name)
		}
	} else {
		fmt.Print(string(parsed.Content))
	}

	return nil
}