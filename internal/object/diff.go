package object

import (
	"fmt"
	"strings"
)

type Edit struct {
	Type string // "EQUAL", "ADD", "DELETE"
	Line string
}

// DiffLines computes the Longest Common Subsequence (LCS) 
// to find the exact line-by-line differences between two files.
func DiffLines(a, b []string) []Edit {
	m, n := len(a), len(b)
	
	// 1. Build the Dynamic Programming (DP) table
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

	// 2. Backtrack to find the path of edits
	var edits []Edit
	i, j := m, n
	for i > 0 || j > 0 {
		if i > 0 && j > 0 && a[i-1] == b[j-1] {
			// Lines are identical
			edits = append([]Edit{{Type: "EQUAL", Line: a[i-1]}}, edits...)
			i--
			j--
		} else if j > 0 && (i == 0 || dp[i][j-1] >= dp[i-1][j]) {
			// Line exists in new version but not in old (Addition)
			edits = append([]Edit{{Type: "ADD", Line: b[j-1]}}, edits...)
			j--
		} else if i > 0 && (j == 0 || dp[i][j-1] < dp[i-1][j]) {
			// Line exists in old version but not in new (Deletion)
			edits = append([]Edit{{Type: "DELETE", Line: a[i-1]}}, edits...)
			i--
		}
	}
	return edits
}

// GeneratePatch takes two raw file contents and returns a Git-style patch string (+ / -)
func GeneratePatch(a, b string) string {
	var aLines, bLines []string
	
	// Split into lines while handling trailing newlines
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
			patch.WriteString(fmt.Sprintf("+ %s\n", edit.Line))
		case "DELETE":
			patch.WriteString(fmt.Sprintf("- %s\n", edit.Line))
		}
	}
	
	return patch.String()
}
