#!/bin/bash

# --- 1. CLEANUP & INITIALIZATION ---
echo "--- Initializing Repo ---"
rm -rf .why hello.txt
go build -o why main.go
./why init

# --- 2. CREATE COMMIT #1 (Initial) ---
echo ""
echo "--- Creating First Commit ---"
echo "hello world v1" > hello.txt
./why add hello.txt
./why commit -m "First commit: added hello.txt"

# --- 3. CREATE COMMIT #2 (Linked to Commit #1) ---
echo ""
echo "--- Creating Second Commit ---"
echo "hello world v2" > hello.txt
./why add hello.txt
./why commit -m "Second commit: updated hello.txt"

# --- 4. CREATE COMMIT #3 ---
echo ""
echo "--- Creating Third Commit ---"
echo "hello world v3" > hello.txt
./why add hello.txt
./why commit -m "Third commit: updated hello.txt"

# --- 5. THE BIG TEST: RUNNING LOG ---
echo ""
echo "--- VERIFYING HISTORY WITH LOG ---"
./why log

# Cleanup
rm hello.txt
