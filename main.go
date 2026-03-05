package main

import (
	"fmt"
	"git-from-scratch/cmd"
	"git-from-scratch/internal/repo"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: why <command>")
		os.Exit(1)
	}

	//command registry...
	command := os.Args[1]
	switch command {
	case "init":
		runInit()
	case "hash-object":

		writeFlag := os.Args[2]
		filePath := os.Args[3]

		if writeFlag != "-w" {
			fmt.Println("error: only -w flag is supported")
			return
		}

		fmt.Println("Command detected: hash-object")
		fmt.Println("Write flag:", writeFlag)
		fmt.Println("Target file:", filePath)

		data, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error reading file: ", err)
			return
		}

		fmt.Println("File read successfully")
		fmt.Println("File size: ", len(data), "bytes")

	case "cat-file":
		if len(os.Args) < 3 {
			fmt.Println("usage: cat-file <hash>")
			return
		}
		err := cmd.CatFile(".", os.Args[2])
		if err != nil {
			fmt.Println("error: ", err)
		}

	case "write-tree":
		err := cmd.WriteTree(".")
		if err != nil {
			fmt.Println("error: ", err)
		}

	case "commit" :
		if len(os.Args)<4 {
			fmt.Println("usage: commit <tree_hash> <message>")
			return
		}
		err := cmd.Commit(".", os.Args[2], "", os.Args[3])
		if err != nil {
			fmt.Println("error: ", err)
		}
		
	default:
		fmt.Printf("unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

// soft coding the hardcoded directory values..
func runInit() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	repo, err := repo.NewRepository(cwd)
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}

	if err:=repo.Init(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Initialized an empty why repository in .why")
}
