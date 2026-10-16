#!/bin/sh
# ⚡ lazygcp — 1-Line Universal Installer for macOS & Linux
# Usage: curl -sSL https://raw.githubusercontent.com/damian-sztankowski/lazygcp/main/install.sh | sh

set -e

REPO="damian-sztankowski/lazygcp"
BIN_NAME="lazygcp"
INSTALL_DIR="/usr/local/bin"

printf "\033[1;36m╭──────────────────────────────────────────────────────────────╮\033[0m\n"
printf "\033[1;36m│  ⚡ Installing lazygcp — Terminal Cockpit for Google Cloud   │\033[0m\n"
printf "\033[1;36m╰──────────────────────────────────────────────────────────────╯\033[0m\n"

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
  x86_64|amd64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *)
    printf "\033[1;31m✗ Unsupported architecture: %s\033[0m\n" "$ARCH"
    exit 1
    ;;
esac

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

INSTALLED_BIN=""

# 1. Try downloading pre-built binary from latest GitHub Release first (no Go required)
LATEST_TAG="$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" 2>/dev/null | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' || true)"
if [ -n "$LATEST_TAG" ]; then
  VERSION_NO_V="${LATEST_TAG#v}"
  TARBALL="${BIN_NAME}_${VERSION_NO_V}_${OS}_${ARCH}.tar.gz"
  DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_TAG/$TARBALL"
  printf "\033[0;36m→ Downloading prebuilt %s (%s/%s) from GitHub Release %s...\033[0m\n" "$BIN_NAME" "$OS" "$ARCH" "$LATEST_TAG"
  if curl -fsSL "$DOWNLOAD_URL" -o "$TMP_DIR/$TARBALL" 2>/dev/null; then
    tar -xzf "$TMP_DIR/$TARBALL" -C "$TMP_DIR"
    if [ -f "$TMP_DIR/$BIN_NAME" ]; then
      INSTALLED_BIN="$TMP_DIR/$BIN_NAME"
    fi
  fi
fi

# 2. Fallback: Build cleanly from source using Go if prebuilt binary was not downloaded
if [ -z "$INSTALLED_BIN" ]; then
  if command -v go >/dev/null 2>&1; then
    printf "\033[0;36m→ Compiling %s natively with Go (%s/%s)...\033[0m\n" "$BIN_NAME" "$OS" "$ARCH"
    GOBIN="$TMP_DIR" go install "github.com/$REPO@latest"
    if [ -f "$TMP_DIR/$BIN_NAME" ]; then
      INSTALLED_BIN="$TMP_DIR/$BIN_NAME"
    fi
  else
    printf "\033[1;31m✗ Could not find a prebuilt release binary and 'go' is not installed.\033[0m\n"
    printf "  Install Go from https://go.dev/dl/ or check https://github.com/%s/releases\n" "$REPO"
    exit 1
  fi
fi

chmod +x "$INSTALLED_BIN"

# 3. Place binary directly into /usr/local/bin (already in macOS/Linux PATH)
if [ -w "$INSTALL_DIR" ]; then
  mv "$INSTALLED_BIN" "$INSTALL_DIR/$BIN_NAME"
  FINAL_PATH="$INSTALL_DIR/$BIN_NAME"
elif command -v sudo >/dev/null 2>&1; then
  printf "\033[0;33m→ Moving %s to %s (may prompt for sudo password so it is in your PATH)...\033[0m\n" "$BIN_NAME" "$INSTALL_DIR"
  if sudo mv "$INSTALLED_BIN" "$INSTALL_DIR/$BIN_NAME" 2>/dev/null; then
    FINAL_PATH="$INSTALL_DIR/$BIN_NAME"
  fi
fi

# 4. Fallback to ~/.local/bin or ~/go/bin if sudo is unavailable
if [ -z "$FINAL_PATH" ]; then
  FALLBACK_DIR="$HOME/.local/bin"
  mkdir -p "$FALLBACK_DIR"
  mv "$INSTALLED_BIN" "$FALLBACK_DIR/$BIN_NAME"
  FINAL_PATH="$FALLBACK_DIR/$BIN_NAME"

  case ":$PATH:" in
    *":$FALLBACK_DIR:"*) ;;
    *)
      for rc in "$HOME/.zshrc" "$HOME/.bashrc"; do
        if [ -f "$rc" ] || [ "$rc" = "$HOME/.zshrc" ]; then
          if ! grep -q "$FALLBACK_DIR" "$rc" 2>/dev/null; then
            printf '\nexport PATH="$PATH:%s"\n' "$FALLBACK_DIR" >> "$rc"
          fi
        fi
      done
      printf "\033[1;33m! Added %s to your shell PATH. Run: source ~/.zshrc\033[0m\n" "$FALLBACK_DIR"
      ;;
  esac
fi

printf "\n\033[1;32m✅ Successfully installed %s to %s!\033[0m\n" "$BIN_NAME" "$FINAL_PATH"
printf "   Run \033[1;36mlazygcp\033[0m to launch the live GCP cockpit, or \033[1;36mlazygcp --demo\033[0m for offline demo mode.\n\n"
