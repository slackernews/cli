# typed: false
# frozen_string_literal: true

# This is a GoReleaser Homebrew formula template.
# GoReleaser populates the URLs and SHA256 checksums during release.
class Slackernews < Formula
  desc "{{ .Desc }}"
  homepage "{{ .Homepage }}"
  version "{{ .Version }}"
  {{ if .License -}}
  license "{{ .License }}"
  {{ end -}}

  {{ range .Dependencies -}}
  depends_on "{{ . }}"
  {{ end -}}

  on_macos do
    on_intel do
      url "https://github.com/slackernews/slackernews-cli/releases/download/v{{ .Version }}/slackernews_{{ .Version }}_darwin_amd64.tar.gz"
      sha256 "DARWIN_AMD64_SHA256"
    end
    on_arm do
      url "https://github.com/slackernews/slackernews-cli/releases/download/v{{ .Version }}/slackernews_{{ .Version }}_darwin_arm64.tar.gz"
      sha256 "DARWIN_ARM64_SHA256"
    end
  end

  def install
    bin.install "slackernews"
  end

  test do
    system "#{bin}/slackernews", "--version"
  end
end
