#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

cd "${ROOT_DIR}"

if ! command -v docker >/dev/null 2>&1; then
  echo "Docker bulunamadı. Lütfen Docker Desktop kurulu ve çalışır durumda olsun." >&2
  exit 1
fi

echo "Backend (Docker) başlatılıyor..."
docker compose up --build -d

cleanup() {
  echo "Backend durduruluyor..."
  docker compose down
}
trap cleanup EXIT

echo "Frontend başlatılıyor..."
exec "${ROOT_DIR}/scripts/dev.sh"
