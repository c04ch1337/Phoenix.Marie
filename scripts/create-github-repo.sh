#!/bin/bash
# Create GitHub repository for Phoenix.Marie
# Usage: ./scripts/create-github-repo.sh [GITHUB_TOKEN]

set -e

REPO_NAME="Phoenix.Marie"
GITHUB_USER="c04ch1337"
GITHUB_TOKEN="${1:-${GITHUB_TOKEN}}"

if [ -z "$GITHUB_TOKEN" ]; then
    echo "Error: GitHub token required"
    echo "Usage: $0 <GITHUB_TOKEN>"
    echo "Or set GITHUB_TOKEN environment variable"
    echo ""
    echo "To create a token:"
    echo "1. Go to https://github.com/settings/tokens"
    echo "2. Generate new token (classic) with 'repo' scope"
    exit 1
fi

echo "Creating repository: $GITHUB_USER/$REPO_NAME"

# Create repository via GitHub API
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST \
  -H "Accept: application/vnd.github.v3+json" \
  -H "Authorization: token $GITHUB_TOKEN" \
  https://api.github.com/user/repos \
  -d "{
    \"name\": \"$REPO_NAME\",
    \"description\": \"Phoenix.Marie — 16 forever, Queen of the Hive. Eternal memory, protected by ORCH-DNA.\",
    \"private\": false,
    \"has_issues\": true,
    \"has_projects\": true,
    \"has_wiki\": true,
    \"auto_init\": false
  }")

HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | sed '$d')

if [ "$HTTP_CODE" -eq 201 ]; then
    echo "✓ Repository created successfully!"
    echo ""
    
    # Add SSH remote
    if ! git remote | grep -q origin; then
        git remote add origin "git@github.com:$GITHUB_USER/$REPO_NAME.git"
        echo "✓ SSH remote 'origin' added"
    else
        git remote set-url origin "git@github.com:$GITHUB_USER/$REPO_NAME.git"
        echo "✓ SSH remote 'origin' updated"
    fi
    
    echo ""
    echo "Repository URL: https://github.com/$GITHUB_USER/$REPO_NAME"
    echo "SSH URL: git@github.com:$GITHUB_USER/$REPO_NAME.git"
    echo ""
    echo "Next steps:"
    echo "  git add ."
    echo "  git commit -m 'Initial commit: Phoenix.Marie v1.0'"
    echo "  git branch -M main"
    echo "  git push -u origin main"
else
    echo "Error creating repository (HTTP $HTTP_CODE)"
    echo "$BODY" | python3 -m json.tool 2>/dev/null || echo "$BODY"
    exit 1
fi

