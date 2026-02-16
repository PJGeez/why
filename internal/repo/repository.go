package repo

import (
	"fmt"
	"os"
	"path/filepath"
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
	headContent := []byte("ref: refs/head/master\n")

	return os.WriteFile(headpath, headContent, 0644)
}