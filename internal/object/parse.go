package object

import (
	"bytes"
	"fmt"
)

type ParsedObject struct {
	Type string
	Size int
	Content []byte
}

func ParseObject(data []byte) (*ParsedObject, error) {
	nullIndex := bytes.IndexByte(data, 0)
	if nullIndex == -1 {
		return nil, fmt.Errorf("invalid object format")
	}

	header := string(data[:nullIndex])
	content := data[nullIndex+1 :]

	var objType string
	var size int

	_, err := fmt.Sscanf(header, "%s %d", &objType, &size)
	if err != nil {
		return nil, err
	}

	if size != len(content) {
		return nil, fmt.Errorf("size mismatch")
	}

	return &ParsedObject{
		Type: objType,
		Size: size,
		Content: content,
	}, nil
}