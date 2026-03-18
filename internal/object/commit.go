package object

import (
	"fmt"
	"strings"
	"time"
)

type Commit struct {
	Tree    string
	Parent  string
	Author  string
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

func ParseCommit(data []byte) (*Commit, error) {
	c := &Commit{}
	content := string(data)
	lines := strings.Split(content, "\n")

	i := 0
	for ; i < len(lines); i++ {
		line := lines[i]
		if line == "" {
			//message starts after the blank line
			break
		}

		parts := strings.SplitN(line, " ", 2)
		if len(parts) < 2 {
			continue
		}

		switch parts[0] {
		case "tree":
			c.Tree = parts[1]
		case "parent":
			c.Parent = parts[1]
		case "author":
			// In a real Git parser, we'd extract the name, email, and timestamp.
			// For now, we'll store the whole line as author.
			c.Author = parts[1]
		}
	}

	// Join the rest as the message
	if i+1 < len(lines) {
		c.Message = strings.TrimSpace(strings.Join(lines[i+1:], "\n"))
	}

	return c, nil
}