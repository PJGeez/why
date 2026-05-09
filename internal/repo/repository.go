package repo

import (
	"fmt"
	"git-from-scratch/internal/object"
	"os"
	"path/filepath"
	"strings"
)

type Repository struct {
	WorkTree string
	GitDir   string
}

func NewRepository(worktree string) (*Repository, error) {
	gitDir := filepath.Join(worktree, ".why")

	if _, err := os.Stat(gitDir); err == nil {
		return nil, fmt.Errorf("why repository already exists...")
	}

	return &Repository{
		WorkTree: worktree,
		GitDir:   gitDir,
	}, nil
}

func (r *Repository) Init() error {
	dirs := []string{
		r.GitDir,
		filepath.Join(r.GitDir, "objects"),
		filepath.Join(r.GitDir, "refs"),
		filepath.Join(r.GitDir, "refs", "heads"),
	}

	for _, dir := range dirs {
		if err := os.Mkdir(dir, 0755); err != nil {
			return err
		}
	}

	headpath := filepath.Join(r.GitDir, "HEAD")
	headContent := []byte("ref: refs/heads/master\n")

	return os.WriteFile(headpath, headContent, 0644)
}

func (r *Repository) GetHeadCommit() (string, error) {
	headPath := filepath.Join(r.GitDir, "HEAD")
	data, err := os.ReadFile(headPath)
	if err != nil {
		return "", err
	}

	content := strings.TrimSpace(string(data))

	const refPrefix = "ref: "
	if strings.HasPrefix(content, refPrefix) {
		refPath := strings.TrimSpace(strings.TrimPrefix(content, refPrefix))
		fullRefPath := filepath.Join(r.GitDir, refPath)
		refData, err := os.ReadFile(fullRefPath)
		if err != nil {
			if os.IsNotExist(err) {
				return "", nil
			}
			return "", err
		}
		return strings.TrimSpace(string(refData)), nil
	}
	return content, nil
}

func (r *Repository) GetCurrentBranch() (string, error) {
	headPath := filepath.Join(r.GitDir, "HEAD")
	data, err := os.ReadFile(headPath)
	if err != nil {
		return "", err
	}

	content := strings.TrimSpace(string(data))
	if strings.HasPrefix(content, "ref: refs/heads/") {
		return strings.TrimPrefix(content, "ref: refs/heads/"), nil
	}

	return "", fmt.Errorf("detached HEAD or unknown branch format")
}

func (r *Repository) SetBranchCommit(branch string, commitHash string) error {
	refPath := filepath.Join(r.GitDir, "refs", "heads", branch)

	if err := os.MkdirAll(filepath.Dir(refPath), 0755); err != nil {
		return err
	}

	return os.WriteFile(refPath, []byte(commitHash+"\n"), 0644)
}

func (r *Repository) ResolveTarget(target string) (string, bool, error) {
	refPath := filepath.Join(r.GitDir, "refs", "heads", target)
	if data, err := os.ReadFile(refPath); err == nil {
		return strings.TrimSpace(string(data)), true, nil
	}

	if len(target) == 40 {
		return target, false, nil
	}

	return "", false, fmt.Errorf("target '%s' is not a valid branch or a commit hash", target)
}

func (r *Repository) UpdateHead(target string, isBranch bool) error {
	headPath := filepath.Join(r.GitDir, "HEAD")
	var content string

	if isBranch {
		content = fmt.Sprintf("ref: refs/heads/%s\n", target)
	} else {
		content = target + "\n"
	}

	return os.WriteFile(headPath, []byte(content), 0644)
}

func (r *Repository) CreateBranch(name string, commitHash string) error {
	branchPath := filepath.Join(r.GitDir, "refs", "heads", name)

	if _, err := os.Stat(branchPath); err == nil {
		return fmt.Errorf("branch %s already exists", name)
	}

	return os.WriteFile(branchPath, []byte(commitHash+"\n"), 0644)
}

func (r *Repository) ListBranches() ([]string, string, error) {
	var branches []string
	headsDir := filepath.Join(r.GitDir, "refs", "heads")

	entries, err := os.ReadDir(headsDir)
	if err != nil && !os.IsNotExist(err) {
		return nil, "", err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			branches = append(branches, entry.Name())
		}
	}

	current, _ := r.GetCurrentBranch()

	if current != "" {
		found := false
		for _, b := range branches {
			if b == current {
				found = true
				break
			}
		}
		if !found {
			branches = append(branches, current)
		}
	}

	return branches, current, nil
}

func (r *Repository) GetCommit(hash string) (*object.Commit, error) {
	data, err := object.ReadObject(r.WorkTree, hash)
	if err != nil {
		return nil, err
	}

	obj, err := object.ParseObject(data)
	if err != nil {
		return nil, err
	}

	return object.ParseCommit(obj.Content)
}

func (r *Repository) FindMergeBase(hash1, hash2 string) (string, error) {
	visited := make(map[string]bool)

	queue1 := []string{hash1}
	for len(queue1) > 0 {
		curr := queue1[0]
		queue1 = queue1[1:]

		if curr == "" || visited[curr] {
			continue
		}
		visited[curr] = true

		commitObj, err := r.GetCommit(curr)
		if err == nil && commitObj != nil && commitObj.Parent != "" {
			queue1 = append(queue1, commitObj.Parent)
		}
	}

	visited2 := make(map[string]bool)
	queue2 := []string{hash2}
	for len(queue2) > 0 {
		curr := queue2[0]
		queue2 = queue2[1:]

		if curr == "" || visited2[curr] {
			continue
		}
		visited2[curr] = true

		if visited[curr] {
			return curr, nil
		}

		commitObj, err := r.GetCommit(curr)
		if err == nil && commitObj != nil && commitObj.Parent != "" {
			queue2 = append(queue2, commitObj.Parent)
		}
	}
	return "", fmt.Errorf("no common ancestor found")
}
