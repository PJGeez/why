package object

import (
	"compress/zlib"
	"io"
	"os"
	"path/filepath"
)

func ReadObject(repoPath, hash string) ([]byte, error) {
	dir := filepath.Join(repoPath, ".why", "objects", hash[:2])
	file := filepath.Join(dir, hash[2:])

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader, err := zlib.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return data, nil

}