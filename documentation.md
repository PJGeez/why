# "Why" VCS: Technical Architecture & Development Guide

This document is a comprehensive guide to the internal mechanics of the `why` version control system, a functional clone of Git built in Go. It details the evolution of the project from a simple file-hasher into a multi-branch state machine.

---

## Phase 1 — Repository Initialization (`why init`)
**Objective:** Create the "Black Box" where history is stored.

### Milestones
1.  **Hidden Root:** Creation of the `.why/` directory.
2.  **The Object Database:** Initializing `.why/objects/` to store immutable data.
3.  **The Reference System:** Initializing `.why/refs/heads/` to store branch pointers.
4.  **The Context Pointer:** Creating the `.why/HEAD` file, initialized to `ref: refs/heads/master`.

### Technical "Why"
In Git-like systems, a repository is defined by its metadata. This phase establishes the "Symbolic Reference" system. By starting with `ref: refs/heads/master`, we tell the tool that we are on a branch named `master`, even before that branch actually exists on disk.

---

## Phase 2 — Content-Addressable Storage (Blobs)
**Objective:** Implement a system where data is retrieved by its "fingerprint," not its name.

### Milestones
1.  **Header Injection:** Every object is prefixed with `type size\x00` (e.g., `blob 12\x00hello world`).
2.  **SHA-1 Hashing:** The unique 160-bit identifier (40 hex characters) for every piece of data.
3.  **Zlib Compression:** Compressing content before writing to disk to save space.
4.  **Fan-out Storage:** Splitting the 40-char hash into `ab/cdef...` (2-char directory, 38-char filename) to avoid OS limits on files-per-folder.

### Technical "Why"
This is the **Content-Addressable Storage (CAS)** model. It provides **Automatic Deduplication**: if two users add the same 1GB file, `why` only stores it once because the hash (and thus the address) is identical.

---

## Phase 3 — Tree Objects (The Filesystem Map)
**Objective:** Map hashes back to filenames, paths, and permissions.

### Milestones
1.  **Binary Tree Format:** Entries are stored as `<mode> <name>\x00<20-byte-binary-hash>`.
2.  **Recursive Nesting:** Tree entries can point to either a `blob` (file) or another `tree` (sub-directory).
3.  **`write-tree` logic:** The command that snapshots the current "Index" and returns a single root hash.

### Technical "Why"
A "Blob" has no name—it is just data. The Tree object is the "glue" that gives data a location in a directory. By recursively hashing trees, we can represent an entire project of 10,000 files with a single 40-character root hash.

---

## Phase 4 — The Staging Area (The Index)
**Objective:** Create a persistent "Draft" of the next commit.

### Milestones
1.  **The Manifest:** A JSON or binary list of `(path, hash, mode)` tuples stored at `.why/index`.
2.  **State Separation:** The `add` command moves files from the **Working Directory** to the **Object Database** and updates the **Index**.
3.  **Incremental Staging:** Allowing the user to add files one-by-one.

### Technical "Why"
The Index solves the "Context Switch" problem. It allows a developer to work on 5 files but only "stage" 2 of them for the next commit. It is the bridge between the volatile disk and the permanent history.

---

## Phase 5 — Commits & History Traversal
**Objective:** Link snapshots into a timeline (The Directed Acyclic Graph).

### Milestones
1.  **Commit Format:**
    ```text
    tree <root-tree-hash>
    parent <previous-commit-hash>
    author <name> <timestamp>
    message <text>
    ```
2.  **`why log`:** A loop that resolves `HEAD`, reads a commit, finds its `parent`, and repeats.
3.  **Object Parsing:** Implementing logic to strip headers (`type size\x00`) to extract raw commit data.

### Technical "Why"
A commit is a "Wrapper." It gives a snapshot (the tree) a reason for existing (the message) and a place in history (the parent). This creates a "Chain of Truth" that cannot be altered without changing all subsequent hashes.

---

## Phase 6 — Automated Committing
**Objective:** Transform individual utilities into a cohesive workflow engine.

### Milestones
1.  **Parent Resolution:** Automatically finding the hash in `HEAD` to use as the `parent` for the next commit.
2.  **Ref Mutation:** Automatically overwriting the branch file (e.g., `master`) with the new commit hash.
3.  **Atomic Workflow:** Replacing 4 manual steps with 1 command: `why commit -m "msg"`.

### Technical "Why"
Automation ensures the **Consistency Invariant**: the current branch must always point to the most recent commit. Manual linking is error-prone; Phase 6 makes the tool a "User-Friendly" VCS.

---

## Phase 7 — Checkout (The Time Machine)
**Objective:** Materialize any historical snapshot onto the physical disk.

### Milestones
1.  **Destructive Cleanup:** Wiping the working directory (excluding `.why`) to ensure no "stale" files remain.
2.  **Recursive Unpacking:** A depth-first traversal of tree objects to recreate directories and write blobs.
3.  **Detached HEAD:** Allowing `HEAD` to point to a raw hash (direct) instead of a branch (symbolic).

### Technical "Why"
Checkout is the most dangerous and powerful command. It transforms the repository from a "history viewer" into a **State Machine**. It maps the "Immutable World" (objects) back into the "Mutable World" (your disk).

---

## Phase 8 — Branching & Multi-Timeline Development
**Objective:** Manage parallel streams of work and maintain state consistency.

### Milestones
1.  **Branch Pointers:** Creating new files in `refs/heads/` that point to the current commit.
2.  **Branch Listing:** Highlighting the active branch by comparing `HEAD` with the list of files in `refs/heads/`.
3.  **Index Synchronization (Phase 8.3):** Rebuilding the `.why/index` from the target tree during a `checkout`.

### Technical "Why"
Branches are "Cheap References." Creating a branch in `why` (and Git) takes 0.001 seconds because it is just writing 40 bytes to a new text file. 

**Index Synchronization** is the critical "Secret Sauce." Without it, checking out an old version would leave your staging area (the index) pointing to the *future* version, causing `why status` to show massive errors. By syncing the index, we ensure the "Golden Rule" is met: **Working Directory == Index == HEAD Commit Tree**.

---

## The "Consistency Invariant" (The Golden Rule)
For the system to be "Healthy," the following must always be true after a checkout or commit:
**Working Directory == Index == HEAD Commit Tree**

If these three layers diverge without the user knowing, the VCS has failed its primary job of state management.
