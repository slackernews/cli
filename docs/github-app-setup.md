# GitHub App Setup: SlackerNews Releaser

This document covers the complete setup of the **SlackerNews Releaser** GitHub App, which enables automated releases for the CLI project using GoReleaser.

> **Quick Start:** Run `scripts/setup-release-infra.sh` to automate repository creation. The GitHub App itself must still be created manually (see Step 1).

> **Important:** These steps must be completed by a Slackernews organization admin (e.g., `chuck`) via the GitHub web UI. As an AI assistant, I cannot create GitHub Apps or manage organization-level permissions.

---

## Overview

The **SlackerNews Releaser** GitHub App provides:

- Automated release creation on `slackernews/cli`
- Homebrew formula updates on `slackernews/homebrew-tap`
- Scoop manifest updates on `slackernews/scoop-bucket`
- Secure, scoped access tokens for GoReleaser workflows

---

## Automated Setup (what gh CLI can do)

The `scripts/setup-release-infra.sh` script automates repository creation using the `gh` CLI:

```bash
# Run from the repo root
./scripts/setup-release-infra.sh
```

This script will:
- ✅ Check for `gh` CLI authentication
- ✅ Create `slackernews/homebrew-tap` (public repo with README)
- ✅ Create `slackernews/scoop-bucket` (public repo with README)
- ✅ Check if secrets are already configured on `slackernews/cli`

What it **cannot** do (requires web UI):
- ❌ Create the GitHub App (Step 1)
- ❌ Generate the private key (Step 2)
- ❌ Install the app on the organization (Step 4)

---

## Step 1: Create the GitHub App

1. Navigate to your **personal GitHub settings** or the **slackernews organization settings**
   - Personal: <https://github.com/settings/apps>
   - Organization: <https://github.com/organizations/slackernews/settings/apps>
2. Click **New GitHub App**
3. Fill in the basic information:

   | Field | Value |
   |-------|-------|
   | **GitHub App name** | `SlackerNews Releaser` |
   | **Description** | `Automated release app for SlackerNews CLI distribution` |
   | **Homepage URL** | `https://github.com/slackernews/cli` |
   | **Webhook URL** | *(leave blank or set a dummy URL — not used)* |
   | **Webhook active** | ❌ Uncheck (we don't use webhooks for this app) |

4. Under **Permissions**, expand **Repository permissions** and configure:

   | Permission | Access Level | Purpose |
   |------------|--------------|---------|
   | **Contents** | Read and write | Push release tags, update Homebrew formulas, update Scoop manifests |
   | **Metadata** | Read (default) | Required by GitHub for all apps |

   > **Note:** Only `Contents` read/write is required. No other repository permissions are needed.

5. Under **Where can this GitHub App be installed?**, select:
   - ☑️ **Only on this account** (if created under the org)
   - Or allow **Any account** if you want flexibility (can be restricted later)

6. Click **Create GitHub App**

---

## Step 2: Generate a Private Key

After creating the app, you'll be redirected to the app's settings page.

1. Scroll to the **Private keys** section
2. Click **Generate a private key**
3. GitHub will download a `.pem` file (e.g., `slacker-news-releaser.2025-01-01.private-key.pem`)
4. **Keep this file secure** — you cannot download it again

---

## Step 3: Record the App ID

On the app's settings page (General tab), note the **App ID** displayed near the top.

- Example: `App ID: 1234567`
- You'll need this ID to generate installation tokens

---

## Step 4: Install the App on the slackernews Organization

1. In the left sidebar of the app settings, click **Install App**
2. Click **Install** next to the `slackernews` organization
3. On the repository selection screen, choose:
   - ☑️ **Only select repositories**
   - Select:
     - `slackernews/cli`
     - `slackernews/homebrew-tap`
     - `slackernews/scoop-bucket`

4. Click **Install**
5. Note the **Installation ID** from the URL or page after installation:
   - Example: `https://github.com/organizations/slackernews/settings/installations/9876543`
   - Installation ID: `9876543`

---

## Step 5: Store Secrets in the Main Repository

Go to `slackernews/cli` and add the following secrets:

1. Navigate to: <https://github.com/slackernews/cli/settings/secrets/actions>
2. Click **New repository secret** and add:

   | Secret Name | Value | Description |
   |-------------|-------|-------------|
   | `APP_ID` | The **App ID** from Step 3 | GitHub App identifier |
   | `PRIVATE_KEY` | The entire contents of the `.pem` file from Step 2 | RSA private key for token generation |

   > **Tip:** When copying the private key, include the full contents:
   > ```
   > -----BEGIN RSA PRIVATE KEY-----
   > ...
   > -----END RSA PRIVATE KEY-----
   > ```

---

## Step 6: Configure the Workflow to Use the GitHub App Token

Update `.github/workflows/release.yml` (or create it) to generate an installation access token:

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Generate GitHub App Token
        id: generate-token
        uses: actions/create-github-app-token@v1
        with:
          app-id: ${{ secrets.APP_ID }}
          private-key: ${{ secrets.PRIVATE_KEY }}
          owner: slackernews
          repositories: |
            cli
            homebrew-tap
            scoop-bucket

      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v6
        with:
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ steps.generate-token.outputs.token }}
```

### Alternative: Using a Custom Token Action

If you prefer not to use `actions/create-github-app-token`, you can use a direct API call:

```yaml
      - name: Generate GitHub App Token
        id: generate-token
        run: |
          #!/bin/bash
          set -e

          # Create JWT from App ID and Private Key
          now=$(date +%s)
          iat=$((now - 60))
          exp=$((now + 600))

          header='{"alg":"RS256","typ":"JWT"}'
          payload="{\"iat\":${iat},\"exp\":${exp},\"iss\":${{ secrets.APP_ID }}}"

          b64enc() { openssl base64 -e -A | tr '+/' '-_' | tr -d '='; }

          header_b64=$(printf %s "$header" | b64enc)
          payload_b64=$(printf %s "$payload" | b64enc)
          signature=$(printf '%s.%s' "$header_b64" "$payload_b64" | openssl dgst -sha256 -sign <(echo "${{ secrets.PRIVATE_KEY }}") | b64enc)

          jwt="${header_b64}.${payload_b64}.${signature}"

          # Get installation ID (or hardcode from Step 4)
          installation_id="9876543"

          # Exchange JWT for installation access token
          token=$(curl -s -X POST \
            -H "Authorization: Bearer $jwt" \
            -H "Accept: application/vnd.github+json" \
            "https://api.github.com/app/installations/${installation_id}/access_tokens" | \
            jq -r '.token')

          echo "token=$token" >> "$GITHUB_OUTPUT"
```

---

## Required Permissions Summary

| Repository | Permission | Use Case |
|------------|------------|----------|
| `slackernews/cli` | Contents: read/write | Create GitHub Releases, push tags |
| `slackernews/homebrew-tap` | Contents: read/write | Update Homebrew formula on release |
| `slackernews/scoop-bucket` | Contents: read/write | Update Scoop manifest on release |

---

## Verification

After setup, verify the configuration by:

1. Pushing a test tag (e.g., `v0.0.0-test`) to `slackernews/cli`
2. Confirming the workflow runs successfully
3. Checking that the release appears on `slackernews/cli`
4. Verifying Homebrew formula updates on `slackernews/homebrew-tap`
5. Verifying Scoop manifest updates on `slackernews/scoop-bucket`

---

## Troubleshooting

| Issue | Solution |
|-------|----------|
| `401 Bad credentials` | Check that `APP_ID` and `PRIVATE_KEY` secrets are correctly set |
| `404 Not Found` on installation | Verify the app is installed on the `slackernews` org and has access to the selected repositories |
| `403 Resource not accessible` | Confirm the app has `Contents: read/write` permission and is installed on the target repositories |
| Token expires mid-workflow | Tokens are valid for 1 hour; for long builds, refresh the token or split the workflow |

---

## References

- [GitHub Apps documentation](https://docs.github.com/en/apps)
- [Creating a GitHub App](https://docs.github.com/en/apps/creating-github-apps)
- [Authenticating as a GitHub App](https://docs.github.com/en/apps/creating-github-apps/authenticating-with-a-github-app)
- [actions/create-github-app-token](https://github.com/actions/create-github-app-token)
- [GoReleaser GitHub integration](https://goreleaser.com/ci/actions/)
