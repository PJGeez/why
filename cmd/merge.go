package cmd

import (
	"fmt"
	"git-from-scratch/internal/object"
	"git-from-scratch/internal/repo"
	"os"
	"path/filepath"
)

func Merge(repoPath, targetBranch string) error {
	r := &repo.Repository{WorkTree: repoPath, GitDir: filepath.Join(repoPath, ".why")}

	
	oursHash, err := r.GetHeadCommit()
	if err != nil {
		return fmt.Errorf("could not resolve HEAD: %w", err)
	}

	theirsHash, _, err := r.ResolveTarget(targetBranch)
	if err != nil {
		return fmt.Errorf("could not resolve merge target %s: %w", targetBranch, err)
	}


	baseHash, err := r.FindMergeBase(oursHash, theirsHash)
	if err != nil {
		return fmt.Errorf("could not find shared history: %w", err)
	}


	baseMap := make(map[string]string)
	headMap := make(map[string]string)
	targetMap := make(map[string]string)

	bCommit, _ := r.GetCommit(baseHash)
	if bCommit != nil {
		LoadTreeMap(repoPath, bCommit.Tree, "", baseMap)
	}

	oCommit, _ := r.GetCommit(oursHash)
	if oCommit != nil {
		LoadTreeMap(repoPath, oCommit.Tree, "", headMap)
	}

	tCommit, _ := r.GetCommit(theirsHash)
	if tCommit != nil {
		LoadTreeMap(repoPath, tCommit.Tree, "", targetMap)
	}


	plan := repo.ComputeMergePlan(baseMap, headMap, targetMap)

	// Execute Plan (Mutation Layer)
	conflicts := 0
	for path, decision := range plan {
		fullPath := filepath.Join(repoPath, path)

		switch decision.Action {
		case repo.ActionTakeTheirs:
			data, err := object.ReadObject(repoPath, decision.Hash)
			if err != nil {
				return err
			}
			obj, err := object.ParseObject(data)
			if err != nil {
				return err
			}
			os.WriteFile(fullPath, obj.Content, 0644)
			fmt.Printf("Updating %s\n", path)

		case repo.ActionDelete:
			os.Remove(fullPath)
			fmt.Printf("Removing %s\n", path)

		case repo.ActionConflict:
			fmt.Printf("CONFLICT in %s\n", path)
			// Generate conflict content with markers
			oData, _ := object.ReadObject(repoPath, headMap[path])
			oObj, _ := object.ParseObject(oData)
			tData, _ := object.ReadObject(repoPath, targetMap[path])
			tObj, _ := object.ParseObject(tData)

			conflictContent := fmt.Sprintf("<<<<<<< HEAD\n%s\n=======\n%s\n>>>>>>> %s\n",
				string(oObj.Content), string(tObj.Content), targetBranch)
			os.WriteFile(fullPath, []byte(conflictContent), 0644)
			conflicts++
		}
	}

	if conflicts > 0 {
		return fmt.Errorf("automatic merge failed; fix conflicts and then commit the result")
	}

	fmt.Println("Automatic merge successful. Ready to commit.")
	return nil
}
