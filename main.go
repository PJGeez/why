package main

import (
	"fmt"
	"os"
)

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
	repoPath := ".why"

	//checking if the repo already exists
	if _, err := os.Stat(repoPath); err == nil {
		fmt.Println("why repository already exists")
		os.Exit(1)
	}

	//creating directory structures
	dirs := []string{
		".why",
		".why/objects",
		".why/refs",
	}

	for _, dir := range dirs {
		if err := os.Mkdir(dir, 0755); err!=nil {
			fmt.Println("error creating dir %s: %v\n", dir, err)
			os.Exit(1)
		}
	}

	//create HEAD file
	headPath := ".why/HEAD"
	headContent := []byte("ref: refs/head/master\n")

	if err := os.WriteFile(headPath, headContent, 0644); err != nil {
		fmt.Println("error writing HEAD: %v\n",err)
		os.Exit(1)
	}

	fmt.Println("Initialize an empty why repository in .why")
}