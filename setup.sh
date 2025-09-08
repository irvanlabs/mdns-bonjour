#!/bin/bash
set -e

# --- CONFIG ---
REPO_URL="https://github.com/username/mdns-go.git"   # ganti dengan repo kamu
INSTALL_DIR="/opt/mdns"
SERVICE_NAME="mdns"

echo "📥 Clone/update source..."
if [ -d "$INSTALL_DIR/.git" ]; then
  cd "$INSTALL_DIR"
  git pull
else
  sudo rm -rf "$INSTALL_DIR"
  sudo git clone "$REPO_URL" "$INSTALL_DIR"
  cd "$INSTALL_DIR"
fi

echo "🔨 Build binary..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o mdns .

echo "⚙️  Setup env file..."
if [ ! -f "$INSTALL_DIR/$SERVICE_NAME.env" ]; then
cat <<EOF | sudo tee "$INSTALL_DIR/$SERVICE_NAME.env" > /dev/null
# contoh konfigurasi
SERVICE_NAME=vms-api
SERVICE_TYPE=_http._tcp
SERVICE_DOMAIN=local.
PORT=8080
EOF
fi

echo "📝 Create systemd unit..."
cat <<EOF | sudo tee /etc/systemd/system/$SERVICE_NAME.service > /dev/null
[Unit]
Description=mDNS Service
After=network-online.target
Wants=network-online.target

[Service]
ExecStart=$INSTALL_DIR/mdns
WorkingDirectory=$INSTALL_DIR
User=root
EnvironmentFile=$INSTALL_DIR/$SERVICE_NAME.env
Restart=on-failure

[Install]
WantedBy=multi-user.target
EOF

echo "🔒 Fix SELinux context..."
if command -v semanage >/dev/null 2>&1; then
  sudo semanage fcontext -a -t bin_t "$INSTALL_DIR(/.*)?"
fi
sudo restorecon -Rv "$INSTALL_DIR"

echo "🔄 Reload systemd & enable service..."
sudo systemctl daemon-reexec
sudo systemctl enable --now $SERVICE_NAME

echo "✅ Done! Service status:"
sudo systemctl status $SERVICE_NAME --no-pager -l
