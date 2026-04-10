# "Why" VCS: Technical Architecture & Development Manual

This document is the definitive technical reference for the `why` version control system. It provides a granular breakdown of every phase and sub-step of development, detailing the objectives and the "Technical Why" behind each architectural decision.

---

## Phase 1 — Repository Initialization
**Objective:** Establish the workspace boundary and storage root.

### 1.1 — Hidden Root (`.why/`)
*   **Objective:** Create a hidden directory to house all repository metadata.
*   **Technical Why:** Separation of concerns. By hiding metadata in `.why/`, the tool avoids cluttering the user's working directory while maintaining a clear project boundary. It also signals to other tools that this is a managed repository.

### 1.2 — Sub-directory Structure
*   **Objective:** Initialize `objects/`, `refs/`, and `refs/heads/`.
*   **Technical Why:** Pre-allocating the filesystem structure ensures that subsequent commands (like `add` or `commit`) do not fail due to missing directories. It establishes the "Layout Invariant" required for the rest of the tool.

### 1.3 — Initial HEAD Setup
*   **Objective:** Create the `.why/HEAD` file pointing to `ref: refs/heads/master`.
*   **Technical Why:** This implements the concept of a "Symbolic Reference." It allows the tool to track which branch the user is on even before any commits have been made or any branch files actually exist.

---

## Phase 2 — Content-Addressable Storage (CAS)
**Objective:** Implement a system where data is retrieved by its "fingerprint" (hash), not its name.

### 2.1 — SHA-1 Hashing & Header Injection
*   **Objective:** Generate unique IDs and prefix objects with `type size\x00`.
*   **Technical Why:** **Immutability and Integrity.** The header format mirrors Git and ensures the tool knows how to parse the data before reading the full content. Hashing ensures that if even one bit of a file changes, it receives a new ID, preventing "silent corruption."

### 2.2 — Zlib Compression
*   **Objective:** Compress object data before writing to disk.
*   **Technical Why:** **Storage Efficiency.** Version control systems store thousands of versions of a project. Without compression, the `.why` folder would quickly grow larger than the actual project.

### 2.3 — Fan-out Storage (AB/CDEF...)
*   **Objective:** Split the 40-character hash into a 2-char directory and 38-char filename.
*   **Technical Why:** **Filesystem Performance.** Many operating systems and filesystems (like FAT32 or older ext versions) suffer significant performance degradation when a single folder contains thousands of files. Fan-out distributes objects across 256 sub-directories.

---

## Phase 3 — Trees & File Metadata
**Objective:** Map raw hashes back to filenames, paths, and directory structures.

### 3.1 — Binary Tree Serialization
*   **Objective:** Convert file lists into the binary format: `<mode> <name>\x00<20-byte-binary-hash>`.
*   **Technical Why:** **Space Efficiency.** Storing the 20-byte raw binary hash instead of its 40-character hex representation saves 50% of the space for every entry in the tree.

### 3.2 — Recursive Tree Snapshots
*   **Objective:** Enable tree objects to point to other tree objects (sub-directories).
*   **Technical Why:** **Hierarchical Representation.** This allows the system to represent complex nested folder structures (e.g., `src/internal/utils`) as a single, recursive root hash. One hash can represent an entire filesystem state.

---

## Phase 4 — Staging Area (The Index)
**Objective:** Decouple the volatile Working Directory from the permanent Object Database.

### 4.1 — The Index Manifest (`.why/index`)
*   **Objective:** Maintain a persistent list of staged files and their hashes.
*   **Technical Why:** **Atomic Commits.** The Index acts as a "Drafting Table." It allows a developer to modify 10 files but only choose 3 of them to be part of the next commit, ensuring that commits are clean and focused.

### 4.2 — Incremental Staging (`why add`)
*   **Objective:** Move data into the object database and update the index record.
*   **Technical Why:** **Content Locking.** When a file is "Added," its content is hashed and stored immediately. This ensures that the version in the Index is locked and won't change even if the user continues editing the file on disk before committing.

---

## Phase 5 — Commits & History Traversal
**Objective:** Link individual snapshots into a logical timeline (The DAG).

### 5.1 — Commit Object Format
*   **Objective:** Wrap a tree hash with metadata (author, message, parent hash).
*   **Technical Why:** **Contextual History.** A Tree tells us *what* the project looks like; a Commit tells us *who* changed it, *when*, and *why*. The `parent` field is the "link" that creates the history chain.

### 5.2 — Object Parsing & Header Stripping
*   **Objective:** Implement logic to strip the `type size\x00` header to access raw content.
*   **Technical Why:** **Generic Storage.** By having a unified parser that strips headers, the tool can treat any object (Blob, Tree, or Commit) using the same underlying read logic before passing it to specific type-parsers.

### 5.3 — History Walking (`why log`)
*   **Objective:** Implement a recursive loop that follows parent pointers backward.
*   **Technical Why:** **Traceability.** This milestone enables the user to view the entire evolution of the project from the current moment back to the very first commit.

---

## Phase 6 — Automated Committing
**Objective:** Transform independent utilities into a cohesive workflow engine.

### 6.1 — Automatic Parent Resolution
*   **Objective:** Automatically find the current hash in `HEAD` to use as the `parent`.
*   **Technical Why:** **Automation.** Users should not have to manually look up hashes to link their history. This ensures the chain of history is never broken by human error.

### 6.2 — Atomic Workflow
*   **Objective:** Combine `write-tree`, `commit`, and branch updating into one command.
*   **Technical Why:** **User Experience.** This turns the tool into a professional VCS where a single action (`why commit`) handles the complex internal orchestration of state management.

---

## Phase 7 — Checkout & State Restoration
**Objective:** Materialize any historical snapshot onto the physical disk (Time Travel).

### 7.1 — Target Resolution (Branch vs. Hash)
*   **Objective:** Build a resolver that distinguishes between names (master) and IDs (hashes).
*   **Technical Why:** **Flexibility.** This provides the "Human Layer" of the tool, allowing users to think in terms of names (branches) while the system operates on hashes.

### 7.2 — Recursive Unpacking (`UnpackTree`)
*   **Objective:** Depth-first traversal of tree objects to recreate directories and write blobs.
*   **Technical Why:** **Reconstruction.** This is the critical act of turning the "Immutable World" of the database back into the "Mutable World" of files you can actually edit and run.

### 7.3 — Destructive Worktree Cleanup
*   **Objective:** Wipe existing files (except `.why` and essentials) before restoration.
*   **Technical Why:** **Determinism.** To guarantee the working directory looks *exactly* like the chosen snapshot, we must remove "stale" files that were not part of that historical version.

---

## Phase 8 — Branching & Consistency
**Objective:** Enable parallel streams of work while maintaining system integrity.

### 8.1 — Branch Management
*   **Objective:** Create cheap, named pointers (`refs/heads/`) to specific commits.
*   **Technical Why:** **Cheap Divergence.** In `why` (and Git), creating a branch is nearly instantaneous because it only involves writing 40 bytes to a new text file. This enables non-linear development.

### 8.2 — Index Synchronization (Phase 8.3)
*   **Objective:** Rebuild the `.why/index` from the target tree during a `checkout`.
*   **Technical Why:** **The Golden Invariant.** For the system to be healthy, the **Working Directory == Index == HEAD Commit Tree**. Without syncing the index, `why status` would report false changes after every checkout.

---

## Phase 9 — Diff & Status Engine
**Objective:** Transition from simple tracking to intelligent state comparison.

### 9.1 — File-Level Diff (Status Engine)
*   **Objective:** Detect content-aware modifications using hash comparison.
*   **Technical Why:** Previously, `status` only checked if a file existed. By hashing live files and comparing them to the Index and HEAD, the tool can now identify `modified`, `deleted`, and `new file` states accurately.
*   **Why this is Necessary:**
    *   **From File Tracking to Content Tracking:** Before 9.1, the tool was "blind" to edits. If a file's name didn't change, the tool assumed the content hadn't either. Hashing provides a "fingerprint" that makes the tool aware of every single byte changed.
    *   **Enabling the Triangular Model:** Professional VCS workflow relies on the relationship between HEAD (Past), Index (Draft), and Working Directory (Present). Phase 9.1 computes the differences between these layers, allowing for "Staged" vs "Unstaged" visibility.
    *   **Foundation for Line-Diffs:** You cannot perform a line-by-line comparison (Phase 9.2) until the system first identifies *which* files actually contain differences.

### 9.2 — Triangular Comparison Logic
*   **Objective:** Compare Working Directory vs. Index (Unstaged) and Index vs. HEAD (Staged).
*   **Technical Why:** This establishes the foundational "intelligence layer," allowing the tool to interpret history rather than just recording it.

---

## The "Consistency Invariant" (The Golden Rule)
Every architectural choice in this tool serves one goal: **State Integrity.** By ensuring the object database is content-addressable and the references are atomic, the system guarantees that history is immutable, verifiable, and permanent.
