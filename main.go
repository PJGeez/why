package main

import (
	"fmt"
	"git-from-scratch/internal/repo"
	"git-from-scratch/cmd"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("usage: why <command>")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
		
		case "init":
			runInit()
		
		case "hash-object":
			hash, err := cmd.HashObject(".", os.Args[2:])
			if err != nil {
				fmt.Println("error:", err)
			}
			fmt.Println(hash)
		
		case "cat-file":
			if len(os.Args) < 4 {
				fmt.Println("usage: why cat-file -p <hash>")
				return
			}
		
			flag := os.Args[2]
			hash := os.Args[3]
		
			err := cmd.CatFile(".", flag, hash)
			if err != nil {
				fmt.Println("error:", err)
			}
		
		case "write-tree":
			err := cmd.WriteTree(".")
			if err != nil {
				fmt.Println("error:", err)
			}
		
		case "status":
			err := cmd.Status(".")
			if err!=nil{
				fmt.Println("error:", err)
			}
		
		case "commit":
			if len(os.Args) < 4 {
				fmt.Println("usage: why commit <tree_hash> <message>")
				return
			}
		
			err := cmd.Commit(".", os.Args[2], "", os.Args[3])
			if err != nil {
				fmt.Println("error:", err)
			}
		
		case "log":
			err := cmd.Log(".")
			if err != nil {
				fmt.Println("error: ", err)
			}
		
		case "add":
			if len(os.Args) < 3{
				fmt.Println("usage: why add <file1> <file2> ...")
				return
			}
			err := cmd.Add(".", os.Args[2:])
			if err != nil {
				fmt.Println("error: ", err)
			}
		
		case "debug-head":
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Println("error:", err)
				return
			}
		
			repository, err := repo.NewRepository(cwd)
			if err != nil {
				fmt.Println("error:", err)
				return
			}
		
			hash, err := repository.GetHeadCommit()
			if err != nil {
				fmt.Println("error:", err)
			}
			fmt.Println(hash)
		
		default:
			fmt.Printf("unknown command: %s\n", command)
			os.Exit(1)
	}
}

func runInit() {

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	repository, err := repo.NewRepository(cwd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := repository.Init(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Initialized an empty why repository in .why")
}