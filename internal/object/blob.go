package object

import (
	"fmt"
	"bytes"
)

type Blob struct {
	Content []byte
}

func (b *Blob) Serialize() []byte {
	//Sprintf lets the raw string
	header := fmt.Sprintf("blob %d\x00", len(b.Content))
	return append([]byte(header), b.Content...)
	//the contents are stored as blob 7\x00homework, blob 5000\x00[photo bytes]...
}