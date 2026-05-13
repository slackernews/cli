#!/bin/bash
set -euo pipefail

# Slackernews Release Infrastructure Setup Script
# Uses `gh` CLI where possible. Manual steps are documented.

echo "=== Slackernews Release Setup ==="
echo ""
echo "This script automates what it can via gh CLI."
echo "Some steps MUST be done manually in the GitHub web UI."
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
echo "Step 2: Setting up repository secrets..."
echo ""
echo "NOTE: You need the GitHub App ID and private key first."
echo "If you haven't created the app yet, stop here and follow the manual steps below."
echo ""

# Check if secrets are already set
if gh secret list --repo "$ORG/cli" | grep -q "APP_ID\|PRIVATE_KEY"; then
    echo "  ✓ Secrets already configured on $ORG/cli"
else
    echo "  Secrets not found on $ORG/cli"
    echo ""
    echo "  To set them via gh CLI:"
    echo "    gh secret set APP_ID --repo $ORG/cli"
    echo "    gh secret set PRIVATE_KEY --repo $ORG/cli < /path/to/private-key.pem"
    echo ""
fi

echo ""
echo "=== MANUAL STEPS (must be done via GitHub web UI) ==="
echo ""
echo "1. Create the GitHub App 'SlackerNews Releaser':"
echo "   https://github.com/organizations/$ORG/settings/apps/new"
echo ""
echo "   Settings:"
echo "   - Name: SlackerNews Releaser"
echo "   - Homepage URL: https://slackernews.io"
echo "   - Webhook: Uncheck 'Active'"
echo "   - Permissions → Repository permissions:"
echo "     - Contents: Read and write"
echo "     - Metadata: Read"
echo "   - 'Where can this GitHub App be installed?': Only on this account"
echo ""
echo "2. Generate a private key on the app's settings page"
echo "   (download the .pem file)"
echo ""
echo "3. Install the app on the $ORG organization:"
echo "   https://github.com/organizations/$ORG/settings/installations"
echo "   Select these repositories: cli, homebrew-tap, scoop-bucket"
echo ""
echo "4. Get the App ID from the app's settings page"
echo ""
echo "5. Set the secrets on $ORG/cli:"
echo "   gh secret set APP_ID --repo $ORG/cli"
echo "     (paste the App ID number)"
echo "   gh secret set PRIVATE_KEY --repo $ORG/cli < /path/to/private-key.pem"
echo ""
echo "=== Done! ==="
echo ""
echo "After completing the manual steps, push a tag to test:"
echo "  git tag v0.1.0"
echo "  git push origin v0.1.0"
echo ""
