package repo

type MergeAction string

const (
	ActionKeep       MergeAction = "KEEP"        // No changes anywhere
	ActionTakeOurs   MergeAction = "TAKE_OURS"   // Changed only in our branch
	ActionTakeTheirs MergeAction = "TAKE_THEIRS" // Changed only in their branch
	ActionDelete     MergeAction = "DELETE"      // Deleted in one branch
	ActionConflict   MergeAction = "CONFLICT"    // Both changed, different results
)

type MergeDecision struct {
	Path   string
	Action MergeAction
	Hash   string 
}

// intakes 3 snapshots and returns a list of decisions
func ComputeMergePlan(base, ours, theirs map[string]string) map[string]MergeDecision {
	plan := make(map[string]MergeDecision)


	allPaths := make(map[string]bool)
	for p := range base { allPaths[p] = true }
	for p := range ours { allPaths[p] = true }
	for p := range theirs { allPaths[p] = true }

	//apply the Decision Matrix for every single path
	for path := range allPaths {
		bHash := base[path]
		oHash := ours[path]
		tHash := theirs[path]

		//case A: Identical everywhere
		if oHash == bHash && tHash == bHash {
			plan[path] = MergeDecision{path, ActionKeep, oHash}
			continue
		}

		//case B: Changed only in Ours 
		if oHash != bHash && tHash == bHash {
			if oHash == "" {
				plan[path] = MergeDecision{path, ActionDelete, ""}
			} else {
				plan[path] = MergeDecision{path, ActionTakeOurs, oHash}
			}
			continue
		}

		//case C: Changed only in Theirs
		if tHash != bHash && oHash == bHash {
			if tHash == "" { 
				plan[path] = MergeDecision{path, ActionDelete, ""}
			} else {
				plan[path] = MergeDecision{path, ActionTakeTheirs, tHash}
			}
			continue
		}

		//case D: Same change in both branches
		if oHash == tHash {
			plan[path] = MergeDecision{path, ActionKeep, oHash}
			continue
		}
		
		//case E: Divergent Changes
		plan[path] = MergeDecision{path, ActionConflict, ""}
	}

	return plan
}
