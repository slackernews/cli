# Release Setup: Fine-Grained Personal Access Token

This document covers setting up a **fine-grained Personal Access Token** for automated releases via GoReleaser. This is simpler than a GitHub App and doesn't require a callback URL.

> **Quick Start:** Run `scripts/setup-release-infra.sh` to automate repository creation and check token status.

---

## Overview

The release workflow uses a single fine-grained PAT with access to three repositories:

- `slackernews/cli` — Create releases
- `slackernews/homebrew-tap` — Update Homebrew formula
- `slackernews/scoop-bucket` — Update Scoop manifest

In addition, a **second** fine-grained PAT (`ENTERPRISE_PORTAL_TOKEN`) with access to `slackernews/slackernews-enterprise-portal` is required to publish release binaries to the enterprise portal.

---

## Automated Setup (gh CLI)

The `scripts/setup-release-infra.sh` script automates repository creation:

```bash
# Run from the repo root
./scripts/setup-release-infra.sh
```

This will:
- ✅ Check for `gh` CLI authentication
- ✅ Create `slackernews/homebrew-tap` (public repo with README)
- ✅ Create `slackernews/scoop-bucket` (public repo with README)
- ✅ Create `slackernews/slackernews-enterprise-portal` (public repo with README) if it doesn't exist
- ✅ Check if `TAP_TOKEN` secret is already configured
- ✅ Check if `ENTERPRISE_PORTAL_TOKEN` secret is already configured

---

## Step 1: Create a Fine-Grained PAT

1. Go to: https://github.com/settings/personal-access-tokens/new
2. Fill in the token details:

   | Field | Value |
   |-------|-------|
   | **Token name** | `SlackerNews Release Token` |
   | **Expiration** | 90 days (or 1 year) |
   | **Description** | `Automated releases for SlackerNews CLI` |
   | **Resource owner** | `slackernews` |

3. Under **Repository access**, select:
   - ☑️ **Only select repositories**
   - Select these repositories:
     - `slackernews/cli`
     - `slackernews/homebrew-tap`
     - `slackernews/scoop-bucket`

4. Under **Permissions**, expand **Repository permissions** and set:

   | Permission | Access Level |
   |------------|--------------|
   | **Contents** | Read and write |
   | **Metadata** | Read |

### Step 1b: Enterprise Portal Token

Create a second fine-grained PAT for the enterprise portal:

1. Go to: https://github.com/settings/personal-access-tokens/new
2. Fill in the details:
   - **Token name**: `Enterprise Portal Release Token`
   - **Resource owner**: `slackernews`
   - **Repositories**: `slackernews/slackernews-enterprise-portal`
4. Under **Repository permissions**:
   - **Contents**: Read and write
   - **Metadata**: Read

5. Click **Generate token** and **copy the token immediately**

---

## Step 2: Store the Token as a Repository Secret

Set the tokens as secrets on `slackernews/cli`:

```bash
# Via gh CLI
echo "ghp_xxxxxxxx" | gh secret set TAP_TOKEN --repo slackernews/cli
echo "ghp_xxxxxxxx" | gh secret set ENTERPRISE_PORTAL_TOKEN --repo slackernews/cli

# Or from a file
gh secret set TAP_TOKEN --repo slackernews/cli < /path/to/tap-token.txt
gh secret set ENTERPRISE_PORTAL_TOKEN --repo slackernews/cli < /path/to/portal-token.txt
```

Or via the web UI: https://github.com/slackernews/cli/settings/secrets/actions

---

## Step 3: Verify Setup

Push a test tag to confirm everything works:

```bash
git tag v0.1.0
git push origin v0.1.0
```

The workflow will:
1. Build binaries for all platforms
2. Create a GitHub Release
3. Publish Homebrew formula to `slackernews/homebrew-tap`
4. Publish Scoop manifest to `slackernews/scoop-bucket`
5. Generate Windows MSI installer
6. Open PRs to `slackernews/slackernews-enterprise-portal` updating `assets/` on both `main` and the version branch

---

## Required Permissions Summary

| Repository | Permission | Use Case | Token secret |
|------------|------------|----------|--------------|
| `slackernews/cli` | Contents: read/write | Create GitHub Releases | `GITHUB_TOKEN` |
| `slackernews/homebrew-tap` | Contents: read/write | Update Homebrew formula | `TAP_TOKEN` |
| `slackernews/scoop-bucket` | Contents: read/write | Update Scoop manifest | `TAP_TOKEN` |
| `slackernews/slackernews-enterprise-portal` | Contents: read/write | Publish release binaries to portal | `ENTERPRISE_PORTAL_TOKEN` |

---

## Troubleshooting

| Issue | Solution |
|-------|----------|
| `401 Bad credentials` | Regenerate the PAT — it may have expired or been revoked |
| `403 Resource not accessible` | Check the PAT has access to all required repos and Contents: read/write |
| Homebrew tap not updating | Verify `HOMEBREW_TAP_GITHUB_TOKEN` env var is set in the workflow |
| Scoop bucket not updating | Verify `SCOOP_BUCKET_GITHUB_TOKEN` env var is set in the workflow |
| Enterprise portal PRs not created | Verify `ENTERPRISE_PORTAL_TOKEN` secret is set and has access to `slackernews/slackernews-enterprise-portal` |

---

## References

- [Creating a fine-grained PAT](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token)
- [GoReleaser GitHub integration](https://goreleaser.com/ci/actions/)
