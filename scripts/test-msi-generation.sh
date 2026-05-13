#!/bin/bash
# Test that the MSI generation workflow step will correctly replace the version
# placeholder in the WiX source file.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
WXS_FILE="${SCRIPT_DIR}/../packaging/windows/slackernews.wxs"

if [[ ! -f "$WXS_FILE" ]]; then
    echo "FAIL: WiX source file not found at $WXS_FILE"
    exit 1
fi

# Verify the .wxs contains the version placeholder
if ! grep -q 'Version="\$(var.Version)"' "$WXS_FILE"; then
    echo "FAIL: .wxs file does not contain Version=\"\$(var.Version)\" placeholder"
    exit 1
fi

echo "PASS: .wxs file contains the version placeholder"

# Simulate the workflow sed command (must escape \$ in shell)
VERSION="0.1.0"
TMP_WXS=$(mktemp)
trap 'rm -f "$TMP_WXS"' EXIT

sed "s|Version=\"\$(var.Version)\"|Version=\"$VERSION\"|g" "$WXS_FILE" > "$TMP_WXS"

# Verify the placeholder was replaced
if grep -q '\$(var.Version)' "$TMP_WXS"; then
    echo "FAIL: sed command did not replace the version placeholder"
    exit 1
fi

# Verify the correct version was inserted
if ! grep -q "Version=\"$VERSION\"" "$TMP_WXS"; then
    echo "FAIL: sed command did not insert the expected version ($VERSION)"
    exit 1
fi

echo "PASS: sed command correctly replaces version placeholder"
echo ""
echo "All MSI generation tests passed!"
