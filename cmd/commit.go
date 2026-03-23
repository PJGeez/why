package cmd

import (
	"fmt"
	"git-from-scratch/internal/object"
	"git-from-scratch/internal/repo"
)

func Commit(repoPath, message string) error {

	r := &repo.Repository{WorkTree: repoPath, GitDir: repoPath + "/.why"}

	treeHash, err := WriteTree(repoPath)
	if err !=nil{
		return err
	}

	parentHash, err := r.GetHeadCommit()
	if err != nil {
		return err
	}

	branch, err := r.GetCurrentBranch()
	if err != nil {
		return err
	}

	commit := object.Commit{
		Tree: treeHash,
		Parent: parentHash,
		Author: "Prajwal <crazyshit.dev@gmail.com>",
		Message: message,
	}

	data := commit.Serialize()

	commitHash, err := object.WriteObject(repoPath, data)
	if err != nil {
		return err
	}

	err = r.SetBranchCommit(branch, commitHash)
	if err != nil {
		return err
	}

	fmt.Println(commitHash)
	return nil
}