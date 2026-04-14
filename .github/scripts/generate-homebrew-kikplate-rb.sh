#!/usr/bin/env bash
set -euo pipefail

version="${1:?version required (e.g. 1.2.3)}"
sha256="${2:?sha256 required}"

{
  cat <<EOF
class Kikplate < Formula
  desc "Command-line interface for the Kikplate plate registry"
  homepage "https://github.com/kikplate/kikplate"
  url "https://github.com/kikplate/kikplate/archive/refs/tags/v${version}.tar.gz"
  sha256 "${sha256}"
EOF
  cat <<'RUBY'
  license "Apache-2.0"

  depends_on "go" => :build

  def install
    ENV["CGO_ENABLED"] = "0"
    cd "cli" do
      system "go", "build", *std_go_args(ldflags: "-s -w"), "."
    end
  end

  test do
    assert_match "Kikplate CLI", shell_output("#{bin}/kikplate --help")
  end
end
RUBY
}
