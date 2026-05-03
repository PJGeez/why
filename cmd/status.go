package cmd

import (
	"fmt"
	"git-from-scratch/internal/index"
	"git-from-scratch/internal/object"
	"git-from-scratch/internal/repo"
	"os"
	"path/filepath"
	"sort"
)

// LoadIndexMap to load index entries into a map for fast comparison
func LoadIndexMap(idx *index.Index) map[string]string {
	m := make(map[string]string)
	for _, entry := range idx.Entries {
		m[entry.Path] = entry.Hash
	}
	return m
}

// LoadTreeMap to recursively walk a tree object and build a map of paths to hashes
func LoadTreeMap(repoPath, treeHash, currentPath string, m map[string]string) error {
	if treeHash == "" {
		return nil
	}

	data, err := object.ReadObject(repoPath, treeHash)
	if err != nil {
		return err
	}

	obj, err := object.ParseObject(data)
	if err != nil {
		return err
	}

	tree, err := object.ParseTree(obj.Content)
	if err != nil {
		return err
	}

	for _, entry := range tree.Entries {
		relPath := filepath.Join(currentPath, entry.Name)
		if entry.Mode == "040000" { // Directory
			if err := LoadTreeMap(repoPath, entry.Hash, relPath, m); err != nil {
				return err
			}
		} else { // File
			m[relPath] = entry.Hash
		}
	}
	return nil
}

// ScanWorkingDir to scan working directory and hash every file
func ScanWorkingDir(repoPath string) (map[string]string, error) {
	m := make(map[string]string)
	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip internal metadata and binary
		if info.IsDir() {
			if info.Name() == ".why" || info.Name() == ".git" || info.Name() == ".idea" {
				return filepath.SkipDir
			}
			return nil
		}
		// Skip project specific files
		if info.Name() == "why" || info.Name() == "go.mod" || info.Name() == "go.sum" {
			return nil
		}

		relPath, _ := filepath.Rel(repoPath, path)
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Hash using Git-style blob header to match index/storage
		blob := object.Blob{Content: data}
		hash := object.Hash(blob.Serialize())
		m[relPath] = hash
		return nil
	})
	return m, err
}

func Status(repoPath string) error {
	r := &repo.Repository{WorkTree: repoPath, GitDir: filepath.Join(repoPath, ".why")}

	// 1. Get current branch
	branch, err := r.GetCurrentBranch()
	if err != nil {
		branch = "master"
	}
	fmt.Printf("On branch %s\n\n", branch)

	// 2. Load the three snapshots
	idx, _ := index.ReadIndex(repoPath)
	indexMap := LoadIndexMap(idx)

	workingMap, err := ScanWorkingDir(repoPath)
	if err != nil {
		return fmt.Errorf("could not scan working directory: %w", err)
	}

	headMap := make(map[string]string)
	headCommit, _ := r.GetHeadCommit()
	if headCommit != "" {
		treeHash, err := GetTreeFromCommit(repoPath, headCommit)
		if err == nil {
			LoadTreeMap(repoPath, treeHash, "", headMap)
		}
	}

	// 3. Changes to be committed (HEAD vs Index)
	fmt.Println("Changes to be committed:")
	stagedChanges := false
	var stagedKeys []string
	stagedUnique := make(map[string]bool)
	for k := range indexMap { if !stagedUnique[k] { stagedKeys = append(stagedKeys, k); stagedUnique[k] = true } }
	for k := range headMap { if !stagedUnique[k] { stagedKeys = append(stagedKeys, k); stagedUnique[k] = true } }
	sort.Strings(stagedKeys)

	for _, path := range stagedKeys {
		idxHash, inIndex := indexMap[path]
		headHash, inHead := headMap[path]

		if inIndex && !inHead {
			fmt.Printf("  (staged)   new file:  %s\n", path)
			stagedChanges = true
		} else if inIndex && inHead && idxHash != headHash {
			fmt.Printf("  (staged)   modified:  %s\n", path)
			stagedChanges = true
		} else if !inIndex && inHead {
			fmt.Printf("  (staged)   deleted:   %s\n", path)
			stagedChanges = true
		}
	}
	if !stagedChanges { fmt.Println("  (none)") }

	// 4. Changes not staged for commit (Index vs Working Dir)
	fmt.Println("\nChanges not staged for commit:")
	unstagedChanges := false
	var unstagedKeys []string
	unstagedUnique := make(map[string]bool)
	for k := range indexMap { if !unstagedUnique[k] { unstagedKeys = append(unstagedKeys, k); unstagedUnique[k] = true } }
	for k := range workingMap { if !unstagedUnique[k] { unstagedKeys = append(unstagedKeys, k); unstagedUnique[k] = true } }
	sort.Strings(unstagedKeys)

	for _, path := range unstagedKeys {
		idxHash, inIndex := indexMap[path]
		workHash, inWork := workingMap[path]

		if inIndex && !inWork {
			fmt.Printf("  (unstaged) deleted:   %s\n", path)
			unstagedChanges = true
		} else if inIndex && inWork && workHash != idxHash {
			fmt.Printf("  (unstaged) modified:  %s\n", path)
			unstagedChanges = true
		}
	}
	if !unstagedChanges { fmt.Println("  (none)") }

	// 5. Untracked files
	fmt.Println("\nUntracked files:")
	untrackedFound := false
	var workingKeys []string
	for k := range workingMap { workingKeys = append(workingKeys, k) }
	sort.Strings(workingKeys)

	for _, path := range workingKeys {
		if _, inIndex := indexMap[path]; !inIndex {
			fmt.Printf("  %s\n", path)
			untrackedFound = true
		}
	}
	if !untrackedFound { fmt.Println("  (none)") }

	return nil
}
