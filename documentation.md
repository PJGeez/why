# "Why" VCS: Technical Architecture & Development Manual

This document is the definitive technical reference for the `why` version control system. It provides a granular breakdown of every phase and sub-step of development, detailing the objectives and the "Technical Why" behind each architectural decision.

---

## Phase 1 — Repository Initialization
**Objective:** Establish the workspace boundary and storage root.

### 1.1 — Hidden Root (`.why/`)
*   **Objective:** Create a hidden directory to house all repository metadata.
*   **Technical Why:** Separation of concerns. By hiding metadata in `.why/`, the tool avoids cluttering the user's working directory while maintaining a clear project boundary.

### 1.2 — Sub-directory Structure
*   **Objective:** Initialize `objects/`, `refs/`, and `refs/heads/`.
*   **Technical Why:** Pre-allocating the filesystem structure ensures that subsequent commands do not fail due to missing directories.

### 1.3 — Initial HEAD Setup
*   **Objective:** Create the `.why/HEAD` file pointing to `ref: refs/heads/master`.
*   **Technical Why:** Implements symbolic references, allowing the tool to track the branch even before the first commit.

---

## Phase 2 — Content-Addressable Storage (CAS)
**Objective:** Implement a system where data is retrieved by its "fingerprint" (hash).

### 2.1 — SHA-1 Hashing & Header Injection
*   **Objective:** Generate unique IDs and prefix objects with `type size\x00`.
*   **Technical Why:** Ensures data integrity. Hashing ensures that if one bit changes, the ID changes, preventing silent corruption.

### 2.2 — Zlib Compression
*   **Objective:** Compress object data before writing to disk.
*   **Technical Why:** Efficiency. Reduces the disk footprint of stored objects across many versions.

### 2.3 — Fan-out Storage (AB/CDEF...)
*   **Objective:** Split the 40-character hash into a 2-char directory and 38-char filename.
*   **Technical Why:** Performance. Distributes objects across 256 sub-directories to prevent OS slowdowns in large folders.

---

## Phase 3 — Trees & File Metadata
**Objective:** Map raw hashes back to filenames and directory structures.

### 3.1 — Binary Tree Serialization
*   **Objective:** Convert file lists into binary format: `<mode> <name>\x00<20-byte-hash>`.
*   **Technical Why:** Space efficiency. Binary hashes save 50% space compared to hex strings.

### 3.2 — Recursive Tree Snapshots
*   **Objective:** Enable tree objects to point to other tree objects (sub-directories).
*   **Technical Why:** Allows representing entire directory hierarchies as a single root hash.

---

## Phase 4 — Staging Area (The Index)
**Objective:** Decouple the Working Directory from the permanent history.

### 4.1 — The Index Manifest (`.why/index`)
*   **Objective:** Maintain a persistent list of staged files.
*   **Technical Why:** Atomic Commits. Allows developers to curate exactly which changes enter the next commit.

### 4.2 — Incremental Staging (`why add`)
*   **Objective:** Move data into the database and lock its version in the index.
*   **Technical Why:** Ensures the version in the index won't change even if the file on disk is edited further before committing.

---

## Phase 5 — Commits & History Traversal
**Objective:** Link snapshots into a logical timeline (The DAG).

### 5.1 — Commit Object Format
*   **Objective:** Wrap a tree hash with metadata (author, message, parent).
*   **Technical Why:** Context. Provides the "Who, When, and Why" behind every snapshot.

### 5.2 — History Walking (`why log`)
*   **Objective:** Implement a loop that follows parent pointers backward.
*   **Technical Why:** Enables full traceability of project evolution from the present to the start.

---

## Phase 6 — Automated Committing
**Objective:** Transform utilities into a cohesive workflow engine.

### 6.1 — Automatic Parent Resolution
*   **Objective:** Automatically find the current hash in `HEAD` to use as the parent.
*   **Technical Why:** Ensures the history chain is never broken by manual error.

### 6.2 — Atomic Workflow
*   **Objective:** Combine tree-writing, commit creation, and branch updating.
*   **Technical Why:** Orchestrates state management into a single user action.

---

## Phase 7 — Checkout & State Restoration
**Objective:** Restore any historical snapshot onto the disk (Time Travel).

### 7.1 — Recursive Unpacking (`UnpackTree`)
*   **Objective:** Recreate files and folders from a tree hash.
*   **Technical Why:** Reconstructs the physical world from the immutable database of hashes.

### 7.2 — Worktree Cleanup
*   **Objective:** Wipe existing files before restoration.
*   **Technical Why:** Guarantees the working directory matches the snapshot exactly, with no "stale" files.

---

## Phase 8 — Branching & Consistency
**Objective:** Enable parallel timelines and maintain system integrity.

### 8.1 — Branch Management
*   **Objective:** Create named pointers (`refs/heads/`) to specific commits.
*   **Technical Why:** Enables cheap, instantaneous divergence for non-linear development.

### 8.2 — Index Synchronization
*   **Objective:** Rebuild the `.why/index` from the target tree during a checkout.
*   **Technical Why:** The Golden Invariant. Ensures **Working Directory == Index == HEAD**. Without this, status reports would be incorrect after a checkout.

---

## Phase 9 — Diff & Status Engine
**Objective:** Transition from simple tracking to intelligent state comparison.

### 9.1 — Smart Status (Hash Comparison)
*   **Objective:** Detect content modifications by comparing disk hashes vs. index hashes.
*   **Technical Why:** Moves from "File Tracking" to "Content Tracking." The tool becomes aware of every byte changed, even if the filename remains the same.

### 9.2 — The Diff Engine (LCS Algorithm)
*   **Objective:** Implement line-by-line comparison to generate patches.
*   **Technical Why:** Uses the **Longest Common Subsequence** to find the most efficient path of edits, making diffs human-readable and minimal.

### 9.3 — Understanding the DP Table & Backtracking
*   **The Scoring Logic (The Memory Map):**
    *   The engine builds a 2D grid comparing the Old version (Vertical) vs. the New version (Horizontal).
    *   **Match:** If lines match, add 1 to the diagonal score. This "anchors" the match.
    *   **No Match:** Carry over the max score from the cell Above or Left.
*   **The Backtracking Rules (Generating the Path):**
    *   **Diagonal Move (↖):** Both versions share the line (**EQUAL**).
    *   **Horizontal Move (←):** Line exists only in the new version (**ADD +**).
    *   **Vertical Move (↑):** Line existed only in the old version (**DELETE -**).

### 9.4 — The `why diff` Command
*   **Objective:** Expose line-level differences to the user via the CLI.
*   **Technical Why:**
    *   **Unstaged Diff:** Compares the Working Directory vs. the Index. This shows what the user has changed but not yet "locked in" with `add`.
    *   **Staged Diff (`--staged`):** Compares the Index vs. the HEAD commit. This allows the user to double-check exactly what they are about to commit.
*   **Outcome:** The tool now provides full observability into the project's evolution, matching the standard Git developer experience.

---

## Phase 10 — The Reconciliation Engine (Merging)
**Objective:** Deterministically reconcile two divergent states of history.

### 10.1 — Graph Traversal & Merge Base
*   **Objective:** Identify the Lowest Common Ancestor (LCA) of two commits.
*   **Technical Why:** To merge two branches, you must first find the last point in history where they were identical. This provides the "Base" state. 
*   **The Algorithm (BFS Queue):**
    *   The engine uses a **Breadth-First Search (BFS)** with a queue to walk backward through parent pointers. 
    *   **Why a Queue?** Using a queue makes the system "Merge-Aware." When we later support merge commits (which have multiple parents), the BFS logic can explore multiple branches of history simultaneously.

### 10.2 — The Decision Matrix (Pure Engine)
*   **Objective:** Reconcile three different states (Base, Ours, Theirs) without touching the disk.
*   **Technical Why:** By separating the **computation** from the **execution**, the merge process becomes deterministic and highly testable. The engine produces a "Merge Plan" before a single file is written.
*   **The Decision Matrix:**
    | Comparison | Result | Logic |
    | :--- | :--- | :--- |
    | Ours == Base, Theirs == Base | **Keep Base** | No changes occurred anywhere. |
    | Ours != Base, Theirs == Base | **Take Ours** | Change happened only in our timeline. |
    | Theirs != Base, Ours == Base | **Take Theirs** | Change happened only in their timeline. |
    | Ours != Base, Theirs != Base, Ours == Theirs | **Take Either** | Both branches made the identical change. |
    | Ours != Base, Theirs != Base, Ours != Theirs | **CONFLICT** | Divergent changes detected. |
*   **Constraint:** At this stage, any concurrent modification to the same file path results in a conflict (file-level merging).

---

### 10.3 — The Execution Pipeline (Mutation Layer)
*   **Objective:** Apply the calculated merge plan to the physical files on disk.
*   **Technical Why:** This is the materialization step. It automates the updating of files that changed in only one branch and deletes files that were removed.

### 10.4 — Conflict Resolution Markers
*   **Objective:** Generate human-readable conflict regions using `<<<<<<< HEAD` markers.
*   **Technical Why:** When the engine cannot decide (because both branches changed the same file), it must safely halt and present both versions to the user, ensuring no data is lost during the merge.

## Operational Guarantees
Every architectural choice in this tool serves one goal: **State Integrity.** By ensuring the object database is content-addressable and the references are atomic, the system guarantees that history is immutable, verifiable, and permanent.
