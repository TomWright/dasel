#!/bin/sh
set -e

REPO="TomWright/dasel"
BINDIR="${BINDIR:-/usr/local/bin}"
VERSION="${VERSION:-latest}"
DRY_RUN="${DRY_RUN:-0}"
OUTPUT_TARGET_ONLY=0

print_usage() {
    cat <<EOF
Usage: install.sh [options]

Install the dasel CLI binary for the detected OS and architecture.

Options:
  -b, --bindir DIR      Installation target directory (default: /usr/local/bin)
  -v, --version VER     Specific version or tag to install (default: latest)
  -d, --dry-run         Print detection and installation plan without downloading
  -t, --target          Print only the release asset target name and exit
  -h, --help            Show this help message
EOF
}

while [ $# -gt 0 ]; do
    case "$1" in
        -b|--bindir)
            BINDIR="$2"
            shift 2
            ;;
        -v|--version)
            VERSION="$2"
            shift 2
            ;;
        -d|--dry-run)
            DRY_RUN=1
            shift
            ;;
        -t|--target|--detect)
            OUTPUT_TARGET_ONLY=1
            DRY_RUN=1
            shift
            ;;
        -h|--help)
            print_usage
            exit 0
            ;;
        *)
            echo "Error: Unknown option: $1" >&2
            print_usage >&2
            exit 1
            ;;
    esac
done

detect_os() {
    raw_os="${TARGET_OS:-$(uname -s)}"
    case "$raw_os" in
        Linux*|linux*)
            echo "linux"
            ;;
        Darwin*|darwin*)
            echo "darwin"
            ;;
        CYGWIN*|MINGW*|MSYS*|Windows*|windows*)
            echo "windows"
            ;;
        *)
            echo "Unsupported operating system: $raw_os" >&2
            exit 1
            ;;
    esac
}

detect_arch() {
    raw_arch="${TARGET_ARCH:-$(uname -m)}"
    case "$raw_arch" in
        x86_64|amd64)
            echo "amd64"
            ;;
        aarch64|arm64)
            echo "arm64"
            ;;
        armv7*|armv6*|armhf|arm)
            echo "arm32"
            ;;
        i386|i686|x86|386)
            echo "386"
            ;;
        *)
            echo "Unsupported architecture: $raw_arch" >&2
            exit 1
            ;;
    esac
}

resolve_target() {
    os="$1"
    arch="$2"

    case "${os}_${arch}" in
        linux_amd64)
            echo "dasel_linux_amd64"
            ;;
        linux_arm64)
            echo "dasel_linux_arm64"
            ;;
        linux_arm32)
            echo "dasel_linux_arm32"
            ;;
        linux_386)
            echo "dasel_linux_386"
            ;;
        darwin_amd64)
            echo "dasel_darwin_amd64"
            ;;
        darwin_arm64)
            echo "dasel_darwin_arm64"
            ;;
        windows_amd64)
            echo "dasel_windows_amd64.exe"
            ;;
        windows_386)
            echo "dasel_windows_386.exe"
            ;;
        *)
            echo "No prebuilt release asset for ${os}/${arch}" >&2
            exit 1
            ;;
    esac
}

OS="$(detect_os)"
ARCH="$(detect_arch)"
TARGET="$(resolve_target "$OS" "$ARCH")"

if [ "$OUTPUT_TARGET_ONLY" -eq 1 ]; then
    echo "$TARGET"
    exit 0
fi

if [ "$VERSION" = "latest" ]; then
    DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${TARGET}"
else
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${TARGET}"
fi

if [ "$DRY_RUN" = "1" ] || [ "$DRY_RUN" = "true" ]; then
    echo "OS: ${OS}"
    echo "Architecture: ${ARCH}"
    echo "Target Asset: ${TARGET}"
    echo "Version: ${VERSION}"
    echo "Download URL: ${DOWNLOAD_URL}"
    echo "Install Directory: ${BINDIR}"
    exit 0
fi

# Download utility detection
if command -v curl >/dev/null 2>&1; then
    FETCH_CMD="curl -sSLf"
elif command -v wget >/dev/null 2>&1; then
    FETCH_CMD="wget -qO-"
else
    echo "Error: curl or wget required to download dasel" >&2
    exit 1
fi

TMP_DIR="$(mktemp -d 2>/dev/null || mktemp -d -t 'dasel')"
cleanup() {
    rm -rf "$TMP_DIR"
}
trap cleanup EXIT INT TERM

TMP_FILE="${TMP_DIR}/${TARGET}"

echo "Downloading ${TARGET} (${VERSION}) from ${DOWNLOAD_URL}..."
if [ "${FETCH_CMD% *}" = "curl" ]; then
    curl -sSLf -o "$TMP_FILE" "$DOWNLOAD_URL"
else
    wget -q -O "$TMP_FILE" "$DOWNLOAD_URL"
fi

chmod +x "$TMP_FILE"

BIN_NAME="dasel"
if [ "$OS" = "windows" ]; then
    BIN_NAME="dasel.exe"
fi

if [ ! -d "$BINDIR" ]; then
    mkdir -p "$BINDIR" 2>/dev/null || {
        echo "Creating ${BINDIR} requires administrative privileges. Running with sudo..."
        sudo mkdir -p "$BINDIR"
    }
fi

DEST="${BINDIR}/${BIN_NAME}"

echo "Installing to ${DEST}..."
if [ -w "$BINDIR" ]; then
    mv "$TMP_FILE" "$DEST"
else
    echo "Writing to ${BINDIR} requires administrative privileges. Running with sudo..."
    sudo mv "$TMP_FILE" "$DEST"
fi

echo "Installed dasel to ${DEST} successfully."
if command -v "$DEST" >/dev/null 2>&1; then
    "$DEST" --version || true
fi
