#!/bin/bash

# --- SETUP ---
echo "--- Starting Grand Cycle Test ---"
rm -rf test-final
mkdir test-final
go build -o test-final/why main.go
cd test-final

# 1. INITIALIZATION (Phase 1)
echo -e "\n1. Initializing..."
./why init

# 2. FIRST COMMIT (Phase 2-6)
echo -e "\n2. Creating Initial Commit..."
echo "Base Content" > README.md
./why add README.md
./why commit -m "initial commit"

# 3. BRANCHING (Phase 8)
echo -e "\n3. Creating 'feature' branch..."
./why branch feature

# 4. MODIFY ON MASTER (Phase 9)
echo -e "\n4. Modifying README on 'master'..."
echo "Base Content + Master Edit" > README.md
./why add README.md
./why commit -m "master edit"

# 5. CHECKOUT FEATURE (Phase 7)
echo -e "\n5. Switching to 'feature' branch..."
./why checkout feature
echo "Current File Content: $(cat README.md)"

# 6. MODIFY ON FEATURE (Phase 10 Conflict Prep)
echo -e "\n6. Modifying README on 'feature'..."
echo "Base Content + Feature Edit" > README.md
./why add README.md
./why commit -m "feature edit"

# 7. THE MERGE (Phase 10)
echo -e "\n7. Switching back to master and merging 'feature'..."
./why checkout master
./why merge feature

# 8. VERIFY CONFLICT (Phase 10.4)
echo -e "\n8. Verifying Conflict Markers..."
cat README.md

# 9. FINAL HISTORY (Phase 5)
echo -e "\n9. Full History Log:"
./why log

echo -e "\n--- Grand Cycle Test Complete ---"
