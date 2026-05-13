#!/bin/bash
set -euo pipefail

# Slackernews Release Infrastructure Setup Script
# Uses `gh` CLI to automate repository creation and secret management.

echo "=== Slackernews Release Setup ==="
echo ""

# Check prerequisites
if ! command -v gh &> /dev/null; then
    echo "ERROR: gh CLI not found. Install it first:"
    echo "  brew install gh  # macOS"
    echo "  or https://github.com/cli/cli#installation"
    exit 1
fi

if ! gh auth status &> /dev/null; then
    echo "ERROR: Not authenticated with gh. Run: gh auth login"
    exit 1
fi

# Get org name
ORG="slackernews"
echo "Using organization: $ORG"
echo ""

# Step 1: Create tap repositories (automated via gh)
echo "Step 1: Creating tap repositories..."

for repo in "homebrew-tap" "scoop-bucket"; do
    if gh repo view "$ORG/$repo" &> /dev/null; then
        echo "  ✓ $ORG/$repo already exists"
    else
        echo "  Creating $ORG/$repo..."
        gh repo create "$ORG/$repo" \
            --public \
            --description "SlackerNews $repo" \
            --add-readme \
            --disable-wiki \
            --disable-issues
        echo "  ✓ Created $ORG/$repo"
    fi
done

echo ""
echo "Step 2: Checking repository secrets..."

# Check if TAP_TOKEN is already set
if gh secret list --repo "$ORG/cli" | grep -q "TAP_TOKEN"; then
    echo "  ✓ TAP_TOKEN already configured on $ORG/cli"
else
    echo ""
    echo "  TAP_TOKEN not found on $ORG/cli"
    echo ""
    echo "  You need a fine-grained Personal Access Token with:"
    echo "    - Resource owner: slackernews"
    echo "    - Repository access: cli, homebrew-tap, scoop-bucket"
    echo "    - Permissions: Contents (read and write)"
    echo ""
    echo "  Create one at: https://github.com/settings/personal-access-tokens/new"
    echo ""
    echo "  Then set it via gh CLI:"
    echo "    gh secret set TAP_TOKEN --repo $ORG/cli < /path/to/your-token.txt"
    echo ""
fi

echo "=== Done! ==="
echo ""
echo "After setting TAP_TOKEN, push a tag to test:"
echo "  git tag v0.1.0"
echo "  git push origin v0.1.0"
echo ""
