#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
FRONTEND_DIR="${ROOT_DIR}/frontend"

if [ ! -d "${FRONTEND_DIR}/node_modules" ]; then
  echo "Önce bağımlılıkları yükleyin: npm --prefix frontend install"
  exit 1
fi

pick_port() {
  if [ -n "${DEV_PORT:-}" ]; then
    echo "${DEV_PORT}"
    return
  fi

  for port in $(seq 5173 5200); do
    if ! lsof -iTCP:${port} -sTCP:LISTEN >/dev/null 2>&1; then
      echo "${port}"
      return
    fi
  done

  echo "Uygun boş port bulunamadı (5173-5200 arası)." >&2
  exit 1
}

DEV_SERVER_PORT="$(pick_port)"

cd "${FRONTEND_DIR}"
exec node ./node_modules/vite/bin/vite.js --host 0.0.0.0 --port "${DEV_SERVER_PORT}" "$@"
