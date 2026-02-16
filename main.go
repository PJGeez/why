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
	switch os.Args[1] {
	case "init":
		runInit()
	case "hash-object":
		if len(os.Args) < 3 {
			fmt.Println("usage: hash-object <file>")
			return
		}
		err := cmd.HashObject(".", os.Args[2])
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
