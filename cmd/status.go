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

// ANSI Color Constants
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorBold   = "\033[1m"
	ColorWhite  = "\033[37m"
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

		blob := object.Blob{Content: data}
		hash := object.Hash(blob.Serialize())
		m[relPath] = hash
		return nil
	})
	return m, err
}

func Status(repoPath string) error {
	r := &repo.Repository{WorkTree: repoPath, GitDir: filepath.Join(repoPath, ".why")}

	branch, err := r.GetCurrentBranch()
	if err != nil {
		branch = "master"
	}
	fmt.Printf("On branch %s\n\n", branch)

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

	// 1. Changes to be committed (HEAD vs Index) - GREEN
	fmt.Printf("%sChanges to be committed:%s\n", ColorBold, ColorReset)
	stagedChanges := false
	var stagedKeys []string
	stagedUnique := make(map[string]bool)
	for k := range indexMap { if !stagedUnique[k] { stagedKeys = append(stagedKeys, k); stagedUnique[k] = true } }
	for k := range headMap { if !stagedUnique[k] { stagedKeys = append(stagedKeys, k); stagedUnique[k] = true } }
	sort.Strings(stagedKeys)

	for _, path := range stagedKeys {
		if seen := stagedUnique[path]; !seen { continue }
		idxHash, inIndex := indexMap[path]
		headHash, inHead := headMap[path]

		if inIndex && !inHead {
			fmt.Printf("  %snew file:  %s%s\n", ColorGreen, path, ColorReset)
			stagedChanges = true
		} else if inIndex && inHead && idxHash != headHash {
			fmt.Printf("  %smodified:  %s%s\n", ColorGreen, path, ColorReset)
			stagedChanges = true
		} else if !inIndex && inHead {
			fmt.Printf("  %sdeleted:   %s%s\n", ColorGreen, path, ColorReset)
			stagedChanges = true
		}
	}
	if !stagedChanges { fmt.Println("  (none)") }

	// 2. Changes not staged for commit (Index vs Working Dir) - RED
	fmt.Printf("\n%sChanges not staged for commit:%s\n", ColorBold, ColorReset)
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
			fmt.Printf("  %sdeleted:   %s%s\n", ColorRed, path, ColorReset)
			unstagedChanges = true
		} else if inIndex && inWork && workHash != idxHash {
			fmt.Printf("  %smodified:  %s%s\n", ColorRed, path, ColorReset)
			unstagedChanges = true
		}
	}
	if !unstagedChanges { fmt.Println("  (none)") }

	// 3. Untracked files - RED
	fmt.Printf("\n%sUntracked files:%s\n", ColorBold, ColorReset)
	untrackedFound := false
	var workingKeys []string
	for k := range workingMap { workingKeys = append(workingKeys, k) }
	sort.Strings(workingKeys)

	for _, path := range workingKeys {
		if _, inIndex := indexMap[path]; !inIndex {
			fmt.Printf("  %s%s%s\n", ColorRed, path, ColorReset)
			untrackedFound = true
		}
	}
	if !untrackedFound { fmt.Println("  (none)") }

	return nil
}
