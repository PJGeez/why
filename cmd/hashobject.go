package cmd

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
)

func HashObject(repoPath string, args []string) (string, error) {

	if len(args) < 2 {
		return "", fmt.Errorf("usage: why hash-object -w <file>")
	}

	flag := args[0]
	file := args[1]

	if flag != "-w" {
		return "",fmt.Errorf("only -w supported")
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return "",err
	}

	blob := buildBlob(data)

	hash := hashObject(blob)

	compressed, err := compressData(blob)
	if err != nil {
		return "",err
	}

	err = writeObject(hash, compressed)
	if err != nil {
		return "",err
	}

	fmt.Println(hash)

	return hash,nil
}


func buildBlob(content []byte) []byte {

	header := fmt.Sprintf("blob %d\x00", len(content))

	return append([]byte(header), content...)
}


func hashObject(data []byte) string {
	hash := sha1.Sum(data)

	return hex.EncodeToString(hash[:])
}


func compressData(data []byte) ([]byte, error) {
	var buffer bytes.Buffer

	writer := zlib.NewWriter(&buffer)

	_, err := writer.Write(data)
	if err != nil {
		return nil, err
	}

	writer.Close()

	return buffer.Bytes(), nil
}


func writeObject(hash string, compressed []byte) error {
	dir := ".why/objects/" + hash[:2]
	file := dir + "/" + hash[2:]

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	if _, err := os.Stat(file); err == nil {
		return nil
	}

	return os.WriteFile(file, compressed, 0644)
}