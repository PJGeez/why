package object

import (
	"fmt"
	"strings"
)

// ANSI Color Constants
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorBold   = "\033[1m"
)

type Edit struct {
	Type string // "EQUAL", "ADD", "DELETE"
	Line string
}

// DiffLines computes the Longest Common Subsequence (LCS) 
// to find the exact line-by-line differences between two files.
func DiffLines(a, b []string) []Edit {
	m, n := len(a), len(b)
	
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				if dp[i-1][j] > dp[i][j-1] {
					dp[i][j] = dp[i-1][j]
				} else {
					dp[i][j] = dp[i][j-1]
				}
			}
		}
	}

	var edits []Edit
	i, j := m, n
	for i > 0 || j > 0 {
		if i > 0 && j > 0 && a[i-1] == b[j-1] {
			edits = append([]Edit{{Type: "EQUAL", Line: a[i-1]}}, edits...)
			i--
			j--
		} else if j > 0 && (i == 0 || dp[i][j-1] >= dp[i-1][j]) {
			edits = append([]Edit{{Type: "ADD", Line: b[j-1]}}, edits...)
			j--
		} else if i > 0 && (j == 0 || dp[i][j-1] < dp[i-1][j]) {
			edits = append([]Edit{{Type: "DELETE", Line: a[i-1]}}, edits...)
			i--
		}
	}
	return edits
}

// GeneratePatch takes two raw file contents and returns a colored patch string
func GeneratePatch(a, b string) string {
	var aLines, bLines []string
	
	if a != "" {
		aLines = strings.Split(strings.TrimSuffix(a, "\n"), "\n")
	}
	if b != "" {
		bLines = strings.Split(strings.TrimSuffix(b, "\n"), "\n")
	}

	edits := DiffLines(aLines, bLines)
	var patch strings.Builder
	
	for _, edit := range edits {
		switch edit.Type {
		case "EQUAL":
			patch.WriteString(fmt.Sprintf("  %s\n", edit.Line))
		case "ADD":
			patch.WriteString(fmt.Sprintf("%s+ %s%s\n", ColorGreen, edit.Line, ColorReset))
		case "DELETE":
			patch.WriteString(fmt.Sprintf("%s- %s%s\n", ColorRed, edit.Line, ColorReset))
		}
	}
	
	return patch.String()
}

func GenerateConflictContent(ours, theirs string) string {
	return fmt.Sprintf("<<<<<<< HEAD\n%s\n=======\n%s\n>>>>>>> MERGE_TARGET\n", ours, theirs)
}
