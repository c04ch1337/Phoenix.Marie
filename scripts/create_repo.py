#!/usr/bin/env python3
"""
Create GitHub repository for Phoenix.Marie
Requires GitHub personal access token with 'repo' scope
"""
import os
import sys
import json
import subprocess
from urllib.request import Request, urlopen
from urllib.error import HTTPError

REPO_NAME = "Phoenix.Marie"
GITHUB_USER = "c04ch1337"
GITHUB_API = "https://api.github.com"

def get_token():
    """Get GitHub token from environment or command line"""
    token = os.environ.get("GITHUB_TOKEN")
    if not token and len(sys.argv) > 1:
        token = sys.argv[1]
    if not token:
        print("GitHub personal access token required.")
        print("Get one at: https://github.com/settings/tokens")
        print("Token needs 'repo' scope")
        print()
        print("Usage:")
        print(f"  {sys.argv[0]} <GITHUB_TOKEN>")
        print("  or")
        print(f"  GITHUB_TOKEN=<token> {sys.argv[0]}")
        sys.exit(1)
    return token

def create_repo(token):
    """Create repository via GitHub API"""
    url = f"{GITHUB_API}/user/repos"
    data = {
        "name": REPO_NAME,
        "description": "Phoenix.Marie — 16 forever, Queen of the Hive. Eternal memory, protected by ORCH-DNA.",
        "private": False,
        "has_issues": True,
        "has_projects": True,
        "has_wiki": True,
        "auto_init": False
    }
    
    req = Request(url, data=json.dumps(data).encode(), 
                  headers={
                      "Accept": "application/vnd.github.v3+json",
                      "Authorization": f"token {token}",
                      "Content-Type": "application/json"
                  })
    
    try:
        with urlopen(req) as response:
            result = json.loads(response.read())
            return True, result
    except HTTPError as e:
        error_body = e.read().decode()
        try:
            error_json = json.loads(error_body)
            return False, error_json
        except:
            return False, {"message": error_body}

def setup_git_remote():
    """Set up SSH git remote"""
    ssh_url = f"git@github.com:{GITHUB_USER}/{REPO_NAME}.git"
    
    # Check if remote exists
    result = subprocess.run(["git", "remote", "get-url", "origin"], 
                          capture_output=True, text=True)
    
    if result.returncode == 0:
        # Remote exists, update it
        subprocess.run(["git", "remote", "set-url", "origin", ssh_url], check=True)
        print(f"✓ Updated remote 'origin' to {ssh_url}")
    else:
        # Add new remote
        subprocess.run(["git", "remote", "add", "origin", ssh_url], check=True)
        print(f"✓ Added remote 'origin': {ssh_url}")

def main():
    print(f"Creating repository: {GITHUB_USER}/{REPO_NAME}")
    print()
    
    token = get_token()
    if not token:
        print("Error: Token required")
        sys.exit(1)
    
    success, result = create_repo(token)
    
    if success:
        print("✓ Repository created successfully!")
        print()
        print(f"Repository URL: https://github.com/{GITHUB_USER}/{REPO_NAME}")
        print(f"SSH URL: git@github.com:{GITHUB_USER}/{REPO_NAME}.git")
        print()
        
        # Set up git remote
        try:
            setup_git_remote()
        except subprocess.CalledProcessError as e:
            print(f"Warning: Could not set up git remote: {e}")
        
        print()
        print("Next steps:")
        print("  git add .")
        print("  git commit -m 'Initial commit: Phoenix.Marie v1.0'")
        print("  git branch -M main")
        print("  git push -u origin main")
    else:
        print(f"Error creating repository:")
        print(json.dumps(result, indent=2))
        sys.exit(1)

if __name__ == "__main__":
    main()

