#!/bin/bash
set -e

REPO_URL="https://github.com/username/mdns-go.git"
INSTALL_DIR="/opt/mdns"
SERVICE_NAME="mdns"
PORT=49221

if [ -f /etc/os-release ]; then
  . /etc/os-release
  OS=$ID
else
  echo "Tidak bisa mendeteksi OS"
  exit 1
fi

if ! command -v docker >/dev/null 2>&1; then
  echo "Docker tidak ditemukan, install terlebih dahulu."
  exit 1
fi

if [ -d "$INSTALL_DIR/.git" ]; then
  cd "$INSTALL_DIR" && git pull
else
  sudo rm -rf "$INSTALL_DIR"
  sudo git clone "$REPO_URL" "$INSTALL_DIR"
  cd "$INSTALL_DIR"
fi

docker run --rm -v "$INSTALL_DIR":/app -w /app golang:1.25.1-alpine3.22 \
bash -c "go mod tidy && go mod download && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags='-s -w' -o mdns ."

sudo chmod +x "$INSTALL_DIR/mdns"

if [ ! -f "$INSTALL_DIR/$SERVICE_NAME.env" ]; then
cat <<EOF | sudo tee "$INSTALL_DIR/$SERVICE_NAME.env" > /dev/null
SERVICE_NAME=vms-api
SERVICE_TYPE=_http._tcp
SERVICE_DOMAIN=local.
PORT=$PORT
EOF
fi

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

if command -v semanage >/dev/null 2>&1; then
  sudo semanage fcontext -a -t bin_t "$INSTALL_DIR(/.*)?"
fi
sudo restorecon -Rv "$INSTALL_DIR"

case "$OS" in
  ubuntu|debian)
    sudo ufw allow $PORT/tcp
    ;;
  fedora)
    sudo firewall-cmd --permanent --add-port=$PORT/tcp
    sudo firewall-cmd --reload
    ;;
esac

sudo systemctl daemon-reload
sudo systemctl enable --now $SERVICE_NAME
sudo systemctl status $SERVICE_NAME --no-pager -l

