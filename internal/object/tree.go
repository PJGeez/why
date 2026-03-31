package object
import (
	"fmt"
	"encoding/hex"
	"bytes"
)

type TreeEntry struct {
	Mode string
	Name string
	Hash string //hash string
}

type Tree struct {
	Entries []TreeEntry
}

func (t *Tree) Serialize() ([]byte, error) {
	var content []byte

	for _, entry := range t.Entries {
		line := fmt.Sprintf("%s %s\x00", entry.Mode, entry.Name)
		content = append(content, []byte(line)...)

		hashBytes, err := hex.DecodeString(entry.Hash)
		if err != nil {
			return nil, err
		}

		content = append(content, hashBytes...)
	}

	header := fmt.Sprintf("tree %d\x00", len(content))
	return append([]byte(header), content...), nil
}

func ParseTree(data []byte) (*Tree, error) {
	var entries []TreeEntry
	i := 0
	for i < len(data) {
		// 1. Find the space between "mode" and "filename"
		spaceIndex := bytes.IndexByte(data[i:], ' ')
		if spaceIndex == -1 {
			break
		}
		mode := string(data[i : i+spaceIndex])
		i += spaceIndex + 1

		// 2. Find the null terminator after "filename"
		nullIndex := bytes.IndexByte(data[i:], 0)
		if nullIndex == -1 {
			break
		}
		name := string(data[i : i+nullIndex])
		i += nullIndex + 1

		// 3. The next 20 bytes are the binary SHA-1 hash
		if i+20 > len(data) {
			break
		}
		hashBytes := data[i : i+20]
		hashHex := hex.EncodeToString(hashBytes)
		i += 20

		entries = append(entries, TreeEntry{
			Mode: mode,
			Name: name,
			Hash: hashHex,
		})
	}
	return &Tree{Entries: entries}, nil
}

