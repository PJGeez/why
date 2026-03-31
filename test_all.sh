#!/bin/bash

BINARY_PATH="$(pwd)/why"
TEST_REPO="test-repo"

# Build the latest binary
echo "Building 'why' tool..."
go build -o why main.go

mkdir -p "$TEST_REPO"
cp why "$TEST_REPO/"
cd "$TEST_REPO"

show_menu() {
    echo "--------------------------------"
    echo " WHY VCS - INTERACTIVE TESTER "
    echo "--------------------------------"
    echo "Current Branch: $(cat .why/HEAD 2>/dev/null || echo 'Not Init')"
    echo "--------------------------------"
    echo "1. Init Repository"
    echo "2. Create/Add File"
    echo "3. Commit Changes"
    echo "4. Show Status"
    echo "5. Show Log"
    echo "6. Checkout (Time Travel)"
    echo "7. Reset Test Repo (Wipe everything)"
    echo "q. Quit"
    echo "--------------------------------"
}

while true; do
    show_menu
    read -p "Choose an option: " choice
    case $choice in
        1)
            ./why init
            ;;
        2)
            read -p "Filename: " fname
            read -p "Content: " content
            echo "$content" > "$fname"
            ./why add "$fname"
            ;;
        3)
            read -p "Commit message: " msg
            ./why commit -m "$msg"
            ;;
        4)
            ./why status
            ;;
        5)
            ./why log
            ;;
        6)
            read -p "Enter Target (branch or hash): " target
            ./why checkout "$target"
            ;;
        7)
            rm -rf .why *
            cp ../why .
            echo "Test repo reset."
            ;;
        q)
            exit 0
            ;;
        *)
            echo "Invalid option."
            ;;
    esac
    echo ""
done
