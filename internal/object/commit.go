package object

import (
	"fmt"
	"time"
)

type Commit struct {
	Tree string
	Parent string
	Author string
	Message string
}

func (c *Commit) Serialize() []byte {
	now := time.Now()
	timestamp := now.Unix()
	
	_, offsetSeconds := now.Zone()
	
	offsetHours := offsetSeconds / 3600
	offsetMinutes := (offsetSeconds % 3600) / 60
	
	offset := fmt.Sprintf("%+03d%02d", offsetHours, offsetMinutes)

	content := fmt.Sprintf(
		"tree %s\n", c.Tree,
	)

	if c.Parent != "" {
		content += fmt.Sprintf("parent %s\n", c.Parent)
	}

	content += fmt.Sprintf(
		"author %s %d %s\n",
		c.Author,
		timestamp,
		offset,
	)

	content += fmt.Sprintf(
		"committer %s %d %s\n\n",
		c.Author,
		timestamp,
		offset,
	)

	content += c.Message + "\n"

	header := fmt.Sprintf("commit %d\x00", len(content))

	return append([]byte(header), []byte(content)...)
}