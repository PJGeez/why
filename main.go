package main

import (
	"fmt"
	"os"
)
import "git-from-scratch/internal/repo"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: why <command>")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		runInit()
	default:
		fmt.Println("unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}


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