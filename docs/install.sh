#!/bin/sh
set -e

REPO="viptony/miroxy"
INSTALL_DIR="/usr/local/bin"
SERVICE_NAME="miroxy"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    armv6l|armv7l) ARCH="armv6" ;;
    aarch64|arm64) ARCH="arm64" ;;
    x86_64)        ARCH="amd64" ;;
    *)             echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

BINARY="${SERVICE_NAME}-${OS}-${ARCH}"

echo "=> Installing miroxy (${OS}/${ARCH})..."

# Get latest release tag
LATEST=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | cut -d'"' -f4)
if [ -z "$LATEST" ]; then
    echo "Failed to fetch latest release"
    exit 1
fi

echo "=> Downloading ${LATEST}..."
curl -fsSL "https://github.com/${REPO}/releases/download/${LATEST}/${BINARY}" -o /tmp/miroxy
chmod +x /tmp/miroxy

# Install binary
if [ "$(id -u)" -eq 0 ]; then
    mv /tmp/miroxy "${INSTALL_DIR}/miroxy"
else
    sudo mv /tmp/miroxy "${INSTALL_DIR}/miroxy"
fi

echo "=> Installed miroxy to ${INSTALL_DIR}/miroxy"

# Interactive setup
echo ""
printf "  Which port should miroxy listen on? [8080]: "
read -r PORT
PORT=${PORT:-8080}

printf "  Set a bearer token for auth? (leave empty to skip): "
read -r TOKEN

echo ""

# Set up systemd service on Linux
if [ "$OS" = "linux" ] && command -v systemctl >/dev/null 2>&1; then
    echo "=> Setting up systemd service on port ${PORT}..."

    SERVICE_FILE="/etc/systemd/system/${SERVICE_NAME}.service"
    if [ "$(id -u)" -eq 0 ]; then
        SUDO=""
    else
        SUDO="sudo"
    fi

    TOKEN_LINE=""
    if [ -n "$TOKEN" ]; then
        TOKEN_LINE="Environment=MIROXY_TOKEN=${TOKEN}"
    fi

    $SUDO tee "$SERVICE_FILE" >/dev/null <<EOF
[Unit]
Description=Miroxy Proxy Relay
After=network.target

[Service]
Type=simple
ExecStart=${INSTALL_DIR}/miroxy --port ${PORT}
Restart=always
RestartSec=5
${TOKEN_LINE}

[Install]
WantedBy=multi-user.target
EOF

    $SUDO systemctl daemon-reload
    $SUDO systemctl enable "${SERVICE_NAME}"
    $SUDO systemctl start "${SERVICE_NAME}"

    echo "=> miroxy service started on port ${PORT}"
fi

echo ""
echo "  miroxy is ready on port ${PORT}!"
echo ""
echo "  Usage:"
echo "    curl -X POST http://localhost:${PORT}/proxy \\"
echo "      -H 'Content-Type: application/json' \\"
echo "      -d '{\"method\":\"GET\",\"url\":\"https://httpbin.org/get\"}'"
echo ""
echo "  Health check:"
echo "    curl http://localhost:${PORT}/health"
echo ""
echo "  Run manually:"
echo "    miroxy --port ${PORT}"
echo ""
