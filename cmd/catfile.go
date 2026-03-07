package cmd

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"os"
	"strings"
)

func CatFile(repoPath, flag, hash string) error {
	dir := hash[:2]
	file := hash[2:]

	path := repoPath + "/.why/objects" + dir + "/" + file

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	reader, err := zlib.NewReader(bytes.NewReader(data))
	if err!= nil {
		return err
	}

	defer reader.Close()

	decompressed := new(bytes.Buffer)
	_, err = decompressed.ReadFrom(reader)
	if err!= nil {
		return err
	}

	content := decompressed.Bytes()
	nullIndex := bytes.IndexByte(content, 0)
	header := string(content[:nullIndex])
	body := content[nullIndex+1:]

	parts := strings.Split(header, " ")
	objType := parts[0]
	size := parts[1]

	if len(parts) != 2 {
		return fmt.Errorf("invalid object header")
	}

	switch flag {

	case "-p":
		fmt.Printf("%s", body)

	case "-t":
		fmt.Println(objType)

	case "-s":
		fmt.Println(size)

	default:
		return fmt.Errorf("unknown flag %s", flag)
	}

	return nil
}