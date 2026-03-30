package repo

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Repository struct {
	WorkTree string
	GitDir string
}


func NewRepository(worktree string) (*Repository, error){
	gitDir := filepath.Join(worktree, ".why")

	if _, err := os.Stat(gitDir); err == nil {
		return nil, fmt.Errorf("why repository already exists...")
	}

	return &Repository {
		WorkTree: worktree,
		GitDir: gitDir,
	}, nil
}


func (r *Repository) Init() error {
	dirs := []string {
		r.GitDir,
		filepath.Join(r.GitDir, "objects"),
		filepath.Join(r.GitDir, "refs"),
		filepath.Join(r.GitDir, "refs", "heads"),
	}

	for _, dir := range dirs{
		if err := os.Mkdir(dir, 0755); err!=nil {
			return err
		}
	}

	headpath := filepath.Join(r.GitDir, "HEAD")
	headContent := []byte("ref: refs/heads/master\n")

	return os.WriteFile(headpath, headContent, 0644)
}


func (r *Repository) GetHeadCommit() (string, error){
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
			if os.IsNotExist(err){
				return "", nil
			}
			return "", err
		}
		return strings.TrimSpace(string(refData)), nil
	}
	// detached HEAD contains the hash directly
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

	if err := os.MkdirAll(filepath.Dir(refPath), 0755); err!=nil {
		return err
	}

	return os.WriteFile(refPath, []byte(commitHash+"\n"), 0644)
}


func (r *Repository) ResolveTarget(target string) (string, bool, error) {
	// checking if the branch is already present in the ref/heads
	refPath := filepath.Join(r.GitDir, "ref", "heads", target)
	if data, err := os.ReadFile(refPath); err == nil {
		return strings.TrimSpace(string(data)), true, nil
	}
	
	if len(target) == 40 {
		return target, false, nil
	}

	return "", false, fmt.Errorf("target '%s' is not a valid branch or a commit hash", target)
}


