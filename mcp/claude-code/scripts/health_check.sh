#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Source .env if present
ENV_FILE="${SCRIPT_DIR}/../.env"
if [ -f "$ENV_FILE" ]; then
  # shellcheck disable=SC1090
  source "$ENV_FILE"
fi

# In production the domonda-web-server is mounted behind an /api path prefix,
# so the MCP endpoint is https://domonda.app/api/mcp/. When running the web
# server locally there is no /api prefix — use MCP_ENDPOINT to override, e.g.
#   MCP_ENDPOINT=http://localhost:5001/mcp/ ./scripts/health_check.sh
DOMONDA_URL="${DOMONDA_URL:-https://domonda.app}"
MCP_ENDPOINT="${MCP_ENDPOINT:-${DOMONDA_URL}/api/mcp/}"

if [ -z "${DOMONDA_API_KEY:-}" ]; then
  echo "Error: DOMONDA_API_KEY is not set."
  echo "Set it in your environment or in .env"
  exit 1
fi

echo "Checking ${MCP_ENDPOINT} ..."

HTTP_STATUS=$(curl -s -o /dev/null -w "%{http_code}" \
  -X POST \
  -H "Authorization: Bearer ${DOMONDA_API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"health-check","version":"1.0.0"}}}' \
  "${MCP_ENDPOINT}")

if [ "$HTTP_STATUS" -ge 200 ] && [ "$HTTP_STATUS" -lt 300 ]; then
  echo "OK (HTTP ${HTTP_STATUS})"
else
  echo "FAILED (HTTP ${HTTP_STATUS})"
  exit 1
fi
