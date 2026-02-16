package object

import (
	"compress/zlib"
	"os"
	"path/filepath"
)

func WriteObject(repoPath string, data []byte) (string, error) {
	hash := Hash(data)

	dir := filepath.Join(repoPath, ".why", "objects", hash[:2])
	file := filepath.Join(dir, hash[2:])

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return "",err
	}

	f, err := os.Create(file)
	if err != nil {
		return "",err
	}
	defer f.close()

	writer := zlib.NewWriter(f)
	defer writer.close()

	_, err := writer.Write(data)
	if err != nil {
		return "", err
	}

	return hash, nil
}