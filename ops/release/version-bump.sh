#!/bin/bash

# Version bump utility for go-metadata-removal-toolkit
set -euo pipefail

# Current version from main.go
get_current_version() {
    grep 'appVersion =' src/main.go | cut -d'"' -f2
}

# Bump version based on type
bump_version() {
    local version=$1
    local bump_type=$2
    
    # Remove 'v' prefix if present
    version=${version#v}
    
    # Split version into parts
    IFS='.' read -r major minor patch <<< "$version"
    
    case "$bump_type" in
        major)
            major=$((major + 1))
            minor=0
            patch=0
            ;;
        minor)
            minor=$((minor + 1))
            patch=0
            ;;
        patch)
            patch=$((patch + 1))
            ;;
        *)
            echo "Invalid bump type: $bump_type"
            echo "Usage: $0 [major|minor|patch|current]"
            exit 1
            ;;
    esac
    
    echo "v${major}.${minor}.${patch}"
}

# Main logic
case "${1:-current}" in
    current)
        echo "Current version: $(get_current_version)"
        ;;
    major|minor|patch)
        current=$(get_current_version)
        new_version=$(bump_version "$current" "$1")
        echo "Bumping version from $current to $new_version"
        
        # Update version in main.go
        sed -i "s/appVersion = \".*\"/appVersion = \"${new_version#v}\"/" src/main.go
        
        echo "Version updated to $new_version"
        echo "Don't forget to commit the change!"
        ;;
    *)
        echo "Usage: $0 [major|minor|patch|current]"
        exit 1
        ;;
esac
