#!/usr/bin/env bash

set -euo pipefail

SERVICE_USER="ddns-updater"
SERVICE_GROUP="ddns-updater"

CONFIG_DEST_DIR="/etc/ddns-updater"
CONFIG_DEST_PATH="${CONFIG_DEST_DIR}/config.yaml"

STATE_DIR="/var/lib/ddns-updater"

SYSTEMD_DIR="/etc/systemd/system"
SERVICE_DEST_PATH="${SYSTEMD_DIR}/ddns-updater.service"
TIMER_DEST_PATH="${SYSTEMD_DIR}/ddns-updater.timer"

BINARY_PATH="/usr/local/bin/ddns-updater"

usage() {
  cat <<'EOF'
Usage:
  sudo ./setup_ddns_updater.sh <config_yaml_path> <service_file_path> <timer_file_path>

Example:
  sudo ./setup_ddns_updater.sh \
    ./config.yaml \
    ./deploy/ddns-updater.service \
    ./deploy/ddns-updater.timer
EOF
}

require_root() {
  if [[ "${EUID}" -ne 0 ]]; then
    echo "This script must be run as root." >&2
    exit 1
  fi
}

require_file() {
  local path="$1"
  local label="$2"

  if [[ ! -f "${path}" ]]; then
    echo "Error: ${label} not found: ${path}" >&2
    exit 1
  fi
}

create_service_user() {
  if id -u "${SERVICE_USER}" >/dev/null 2>&1; then
    echo "User '${SERVICE_USER}' already exists, skipping creation."
    return
  fi

  echo "Creating system user '${SERVICE_USER}'..."
  useradd \
    --system \
    --no-create-home \
    --shell /usr/sbin/nologin \
    "${SERVICE_USER}"
}

prepare_directories() {
  echo "Creating directories..."
  mkdir -p "${CONFIG_DEST_DIR}"
  mkdir -p "${STATE_DIR}"

  chown -R "${SERVICE_USER}:${SERVICE_GROUP}" "${STATE_DIR}"
  chmod 755 "${CONFIG_DEST_DIR}"
  chmod 755 "${STATE_DIR}"
}

install_config() {
  local src_config="$1"

  echo "Installing config to ${CONFIG_DEST_PATH}..."
  cp "${src_config}" "${CONFIG_DEST_PATH}"
  chown root:"${SERVICE_GROUP}" "${CONFIG_DEST_PATH}"
  chmod 640 "${CONFIG_DEST_PATH}"
}

install_systemd_units() {
  local src_service="$1"
  local src_timer="$2"

  echo "Installing service file to ${SERVICE_DEST_PATH}..."
  cp "${src_service}" "${SERVICE_DEST_PATH}"
  chown root:root "${SERVICE_DEST_PATH}"
  chmod 644 "${SERVICE_DEST_PATH}"

  echo "Installing timer file to ${TIMER_DEST_PATH}..."
  cp "${src_timer}" "${TIMER_DEST_PATH}"
  chown root:root "${TIMER_DEST_PATH}"
  chmod 644 "${TIMER_DEST_PATH}"
}

check_binary() {
  if [[ ! -x "${BINARY_PATH}" ]]; then
    echo "Warning: binary not found or not executable at ${BINARY_PATH}" >&2
    echo "Make sure you already built and copied the binary before starting the service." >&2
  else
    echo "Found binary: ${BINARY_PATH}"
  fi
}

reload_and_enable_timer() {
  echo "Reloading systemd daemon..."
  systemctl daemon-reload

  echo "Enabling timer..."
  systemctl enable ddns-updater.timer

  echo "Starting timer..."
  systemctl restart ddns-updater.timer
}

print_status() {
  echo
  echo "Setup completed."
  echo
  echo "Useful commands:"
  echo "  systemctl status ddns-updater.timer"
  echo "  systemctl status ddns-updater.service"
  echo "  systemctl list-timers --all | grep ddns-updater"
  echo "  journalctl -u ddns-updater.service -n 50 --no-pager"
  echo "  systemctl start ddns-updater.service"
}

main() {
  if [[ $# -ne 3 ]]; then
    usage
    exit 1
  fi

  local config_src="$1"
  local service_src="$2"
  local timer_src="$3"

  require_root
  require_file "${config_src}" "Config file"
  require_file "${service_src}" "Service file"
  require_file "${timer_src}" "Timer file"

  create_service_user
  prepare_directories
  install_config "${config_src}"
  install_systemd_units "${service_src}" "${timer_src}"
  check_binary
  reload_and_enable_timer
  print_status
}

main "$@"