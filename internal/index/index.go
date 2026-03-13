package index

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type IndexEntry struct {
	Hash string `json:"hash"`
	Path string `json:"path"`
	Mode string `json:"mode"`
}

type Index struct {
	Entries []IndexEntry `json:"entries"`
}

func GetIndexPath(repoPath string) string {
	return filepath.Join(repoPath,".why", "index")
}

func ReadIndex(repoPath string) (*Index, error) {
	idx := &Index{Entries: []IndexEntry{}}
	path := GetIndexPath(repoPath)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return idx, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, idx)
	return idx, err
}

func WriteIndex(repoPath string, idx *Index) error {
	path := GetIndexPath(repoPath)
	data, err := json.MarshalIndent(idx, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
