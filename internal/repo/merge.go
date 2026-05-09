package repo

// MergeAction represents the deterministic decision made for a specific file path.
type MergeAction string

const (
	// ActionKeep: No changes were made in either branch, or both made the same change.
	ActionKeep MergeAction = "KEEP"
	
	// ActionTakeOurs: Change occurred only in our current branch (HEAD).
	ActionTakeOurs MergeAction = "TAKE_OURS"
	
	// ActionTakeTheirs: Change occurred only in the target branch we are merging in.
	ActionTakeTheirs MergeAction = "TAKE_THEIRS"
	
	// ActionDelete: The file was deleted in one branch but remained unchanged in the other.
	ActionDelete MergeAction = "DELETE"
	
	// ActionConflict: Both branches modified the file differently relative to the base.
	ActionConflict MergeAction = "CONFLICT"
)

// MergeDecision holds the final verdict for a file and the hash of the version to use.
type MergeDecision struct {
	Path   string
	Action MergeAction
	Hash   string
}

func ComputeMergePlan(base, ours, theirs map[string]string) map[string]MergeDecision {
	plan := make(map[string]MergeDecision)

	// 1. Create a master list of every file path involved in the merge.
	// This ensures we detect files that were added or deleted in any branch.
	allPaths := make(map[string]bool)
	for p := range base { allPaths[p] = true }
	for p := range ours { allPaths[p] = true }
	for p := range theirs { allPaths[p] = true }

	// 2. Iterate through every path and apply the Decision Matrix.
	for path := range allPaths {
		bHash := base[path] // Version in the shared history
		oHash := ours[path] // Version in our current branch
		tHash := theirs[path] // Version in the branch we want to merge

		// Case A: Perfect Harmony
		// If everyone has the same hash, nothing changed.
		if oHash == bHash && tHash == bHash {
			plan[path] = MergeDecision{path, ActionKeep, oHash}
			continue
		}

		// Case B: Ours is newer, Theirs is stale
		// If we changed/deleted it but they didn't touch it relative to base.
		if oHash != bHash && tHash == bHash {
			if oHash == "" {
				plan[path] = MergeDecision{path, ActionDelete, ""}
			} else {
				plan[path] = MergeDecision{path, ActionTakeOurs, oHash}
			}
			continue
		}

		// Case C: Theirs is newer, Ours is stale
		// If they changed/deleted it but we didn't touch it.
		if tHash != bHash && oHash == bHash {
			if tHash == "" {
				plan[path] = MergeDecision{path, ActionDelete, ""}
			} else {
				plan[path] = MergeDecision{path, ActionTakeTheirs, tHash}
			}
			continue
		}

		// Case D: Parallel Evolution (Identical)
		// Both branches made the exact same modification (same resulting hash).
		if oHash == tHash {
			plan[path] = MergeDecision{path, ActionKeep, oHash}
			continue
		}

		// Case E: DIVERGENT EVOLUTION (Conflict!)
		// Both branches changed the file relative to base, but they produced 
		// different hashes. At this stage, we flag this for human resolution.
		plan[path] = MergeDecision{path, ActionConflict, ""}
	}

	return plan
}
