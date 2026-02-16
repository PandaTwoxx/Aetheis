#!/bin/bash
set -e

APP_NAME="aetheis"
DOWNLOAD_URL="https://raw.githubusercontent.com/PandaTwoxx/Aetheis/refs/heads/main/cli/bin/aetheis"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m'

MODE="system"

# Parse arguments
for arg in "$@"; do
    case $arg in
        --local)
        MODE="local"
        shift
        ;;
        *)
        ;;
    esac
done

if [ "$MODE" == "local" ]; then
    INSTALL_DIR="$HOME/.aetheis/bin"
    echo -e "${YELLOW}👉 Running in LOCAL mode. Installing to $INSTALL_DIR${NC}"
else
    INSTALL_DIR="/usr/local/bin"
    echo -e "${YELLOW}👉 Running in SYSTEM mode. Installing to $INSTALL_DIR (requires sudo)${NC}"
fi

echo -e "🚀 Installing ${APP_NAME}..."

# Check curl
if ! command -v curl &> /dev/null; then
    echo -e "${RED}Error: curl is not installed!${NC}"
    exit 1
fi

# Create a temporary directory
TEMP_DIR=$(mktemp -d)
trap 'rm -rf "$TEMP_DIR"' EXIT

# Download
echo -e "⬇️  Downloading from $DOWNLOAD_URL..."
if ! curl -fsSL "$DOWNLOAD_URL" -o "$TEMP_DIR/$APP_NAME"; then
    echo -e "${RED}Download failed! Check your internet connection or the URL.${NC}"
    exit 1
fi

chmod +x "$TEMP_DIR/$APP_NAME"

# Install
echo -e "📂 Installing to $INSTALL_DIR..."
mkdir -p "$INSTALL_DIR"

TARGET_PATH="$INSTALL_DIR/$APP_NAME"

# If existing binary exists, replace it
if [ -f "$TARGET_PATH" ]; then
    echo -e "${YELLOW}♻️  Existing $APP_NAME found. Replacing with updated version...${NC}"
    if [ "$MODE" == "system" ] && [ ! -w "$INSTALL_DIR" ]; then
        sudo rm -f "$TARGET_PATH"
    else
        rm -f "$TARGET_PATH"
    fi
else
    echo -e "${GREEN}🆕 Fresh install detected.${NC}"
fi

if [ "$MODE" == "local" ]; then
    mv -f "$TEMP_DIR/$APP_NAME" "$TARGET_PATH"

    # Update zshrc
    RC_FILE="$HOME/.zshrc"
    EXPORT_CMD="export PATH=\"$HOME/.aetheis/bin:\$PATH\""

    if ! grep -Fq "$HOME/.aetheis/bin" "$RC_FILE"; then
        echo -e "${YELLOW}📝 Adding $INSTALL_DIR to $RC_FILE...${NC}"
        echo "" >> "$RC_FILE"
        echo "# Aetheis CLI" >> "$RC_FILE"
        echo "$EXPORT_CMD" >> "$RC_FILE"
        echo -e "${GREEN}✅ Added to PATH. Please restart your shell or run: source $RC_FILE${NC}"
    else
        echo -e "${GREEN}✅ PATH already configured in $RC_FILE.${NC}"
    fi
else
    if [ -w "$INSTALL_DIR" ]; then
        mv -f "$TEMP_DIR/$APP_NAME" "$TARGET_PATH"
    else
        echo "Requires sudo privileges to install to $INSTALL_DIR"
        sudo mv -f "$TEMP_DIR/$APP_NAME" "$TARGET_PATH"
    fi
fi


echo -e "${GREEN}✅ Successfully installed $APP_NAME!${NC}"
echo -e "Run '${APP_NAME} help' to get started."
