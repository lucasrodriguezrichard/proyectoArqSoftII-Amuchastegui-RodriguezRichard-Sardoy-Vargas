#!/usr/bin/env bash
set -euo pipefail

# Simple E2E test for Reservations API
# Requires: curl, jq (optional but recommended)

USERS_API_URL=${USERS_API_URL:-http://localhost:8080}
RES_API_URL=${RES_API_URL:-http://localhost:8081}
OWNER_ID=${OWNER_ID:-1}
TABLE_NUMBER=${TABLE_NUMBER:-5}
GUESTS=${GUESTS:-4}
MEAL_TYPE=${MEAL_TYPE:-dinner}
DATETIME=${DATETIME:-"2025-11-20T20:00:00Z"}

echo "[1/7] Health checks"
curl -fsS "$USERS_API_URL/health" | jq -r .status || true
curl -fsS "$RES_API_URL/health" | jq -r .status || true

echo "[2/7] Ensure user exists: $USERS_API_URL/api/users/$OWNER_ID"
curl -fsS -o /dev/null -w "%{http_code}\n" "$USERS_API_URL/api/users/$OWNER_ID" || true

echo "[3/7] Create reservation"
create_payload=$(cat <<JSON
{
  "owner_id": "$OWNER_ID",
  "table_number": $TABLE_NUMBER,
  "guests": $GUESTS,
  "date_time": $DATETIME,
  "meal_type": "$MEAL_TYPE",
  "special_requests": "Mesa junto a la ventana"
}
JSON
)

create_resp=$(curl -fsS -H 'Content-Type: application/json' \
  -d "$create_payload" \
  "$RES_API_URL/api/reservations")
echo "$create_resp" | sed 's/.*/[created]/'

if command -v jq >/dev/null 2>&1; then
  RES_ID=$(echo "$create_resp" | jq -r .id)
else
  RES_ID=$(echo "$create_resp" | sed -n 's/.*"id"\s*:\s*"\([^"]\+\)".*/\1/p')
fi

if [ -z "${RES_ID:-}" ]; then
  echo "Failed to parse reservation id" >&2
  exit 1
fi

echo "[4/7] Get reservation by id: $RES_ID"
curl -fsS "$RES_API_URL/api/reservations/$RES_ID" | sed 's/.*/[got]/'

echo "[5/7] List reservations (limit=5)"
curl -fsS "$RES_API_URL/api/reservations?limit=5&offset=0" > /dev/null
echo "[list ok]"

echo "[6/7] Confirm reservation"
curl -fsS -X POST -H 'Content-Type: application/json' \
  -d '{}' "$RES_API_URL/api/reservations/$RES_ID/confirm" > /dev/null
echo "[confirm ok]"

echo "[7/7] Update reservation (guests=6)"
curl -fsS -X PUT -H 'Content-Type: application/json' \
  -d '{"guests":6}' "$RES_API_URL/api/reservations/$RES_ID" > /dev/null
echo "[update ok]"

echo "E2E OK - reservation id: $RES_ID"

