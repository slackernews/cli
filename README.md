# SlackerNews CLI

Browse top links, vote, and comment on [SlackerNews](https://slackernews.io) directly from your terminal.

## Installation

### macOS (Homebrew)

```bash
brew install slackernews/tap/slackernews
```

### Linux

#### Debian/Ubuntu (APT)

```bash
wget https://github.com/slackernews/cli/releases/latest/download/slackernews_linux_amd64.deb
sudo dpkg -i slackernews_linux_amd64.deb
```

#### Fedora/RHEL (DNF)

```bash
sudo rpm -i https://github.com/slackernews/cli/releases/latest/download/slackernews_linux_amd64.rpm
```

#### Alpine (APK)

```bash
wget https://github.com/slackernews/cli/releases/latest/download/slackernews_linux_amd64.apk
sudo apk add --allow-untrusted slackernews_linux_amd64.apk
```

#### Arch Linux

```bash
wget https://github.com/slackernews/cli/releases/latest/download/slackernews_linux_amd64.pkg.tar.zst
sudo pacman -U slackernews_linux_amd64.pkg.tar.zst
```

### Windows

#### Scoop

```powershell
scoop install slackernews
```

#### MSI Installer

Download `slackernews_*_windows_x86_64.msi` from the [latest release](https://github.com/slackernews/cli/releases/latest) and run it.

### Go Install

If you have Go installed:

```bash
go install github.com/slackernews/cli@latest
```

### Manual Download

Download the appropriate archive for your platform from the [releases page](https://github.com/slackernews/cli/releases/latest) and extract the `slackernews` binary to your `$PATH`.

## Configuration

Before using the CLI, configure it with your SlackerNews instance URL and API token:

```bash
slackernews configure --url https://your-instance.slackernews.io --token YOUR_API_TOKEN
```

The API token is stored securely in your OS keychain. For CI or headless environments, set the `SLACKERNEWS_TOKEN` environment variable instead.

## Usage

### Browse top links

```bash
# Default: last 7 days, human-readable table
slackernews top

# Custom duration
slackernews top --duration 30d

# JSON output (for scripting)
slackernews top --json
```

### Search links

```bash
slackernews search kubernetes
slackernews search "machine learning" --json
```

### Vote on links

```bash
slackernews upvote https://example.com/article
slackernews unvote https://example.com/article
```

### Comment on links

```bash
slackernews comment https://example.com/article "Great read!"
```

## Global Flags

| Flag | Description |
|------|-------------|
| `--insecure` | Allow HTTP URLs (development only) |
| `--help` | Show help for any command |

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Authentication failed |
| 3 | Network error |
| 4 | Server error |
| 5 | Rate limited |

## Environment Variables

| Variable | Description |
|----------|-------------|
| `SLACKERNEWS_TOKEN` | API token (overrides keychain) |
| `SLACKERNEWS_TIMEOUT` | HTTP timeout (default: 30s) |

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## License

MIT
