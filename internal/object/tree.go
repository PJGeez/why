package object
import (
	"fmt"
	"encoding/hex"
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