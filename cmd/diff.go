package cmd

import (
	"fmt"
	"git-from-scratch/internal/index"
	"git-from-scratch/internal/object"
	"git-from-scratch/internal/repo"
	"os"
	"path/filepath"
)

func Diff(repoPath string, args []string) error {
	r := &repo.Repository{WorkTree: repoPath, GitDir: filepath.Join(repoPath, ".why")}
	isStaged := false

	if len(args) > 0 && args[0] == "--staged" {
		isStaged = true
	}

	idx, err := index.ReadIndex(repoPath)
	if err != nil {
		return err
	}
	indexMap := LoadIndexMap(idx)

	headMap := make(map[string]string)
	headCommit, _ := r.GetHeadCommit()

	if headCommit != "" {
		treeHash, _ := GetTreeFromCommit(repoPath, headCommit)
		LoadTreeMap(repoPath, treeHash, "", headMap)
	}

	if isStaged {
		// COMPARE INDEX vs HEAD (Staged changes)
		for path, idxHash := range indexMap {
			headHash, inHead := headMap[path]

			if !inHead || headHash != idxHash {
				fmt.Printf("diff --staged %s\n", path)

				var oldContent string
				if inHead {
					data, _ := object.ReadObject(repoPath, headHash)
					obj, _ := object.ParseObject(data)
					oldContent = string(obj.Content)
				}

				data, _ := object.ReadObject(repoPath, idxHash)
				obj, _ := object.ParseObject(data)
				newContent := string(obj.Content)

				fmt.Println(object.GeneratePatch(oldContent, newContent))
			}
		}
	} else {
		// COMPARE WORKING DIR vs INDEX (Unstaged changes)
		workingMap, _ := ScanWorkingDir(repoPath)
		for path, idxHash := range indexMap {
			workHash, inWork := workingMap[path]

			if !inWork || workHash != idxHash {
				fmt.Printf("diff %s\n", path)

				data, _ := object.ReadObject(repoPath, idxHash)
				obj, _ := object.ParseObject(data)
				oldContent := string(obj.Content)

				var newContent string
				if inWork {
					newContentBytes, _ := os.ReadFile(filepath.Join(repoPath, path))
					newContent = string(newContentBytes)
				}

				fmt.Println(object.GeneratePatch(oldContent, newContent))
			}
		}
	}
	return nil
}
