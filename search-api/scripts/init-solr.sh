#!/usr/bin/env bash
set -euo pipefail

SOLR_URL=${SOLR_URL:-http://localhost:8983/solr}
CORE=${SOLR_CORE:-reservations}

echo "Init Solr core schema for $CORE"

schema="$SOLR_URL/$CORE/schema"

curl -fsS -X POST -H 'Content-Type: application/json' "$schema" -d '{
  "add-field": {"name":"id","type":"string","stored":true,"indexed":true}
}' || true

curl -fsS -X POST -H 'Content-Type: application/json' "$schema" -d '{
  "add-field": {"name":"owner_id","type":"string","stored":true,"indexed":true}
}' || true

curl -fsS -X POST -H 'Content-Type: application/json' "$schema" -d '{
  "add-field": {"name":"table_number","type":"pint","stored":true,"indexed":true}
}' || true

curl -fsS -X POST -H 'Content-Type: application/json' "$schema" -d '{
  "add-field": {"name":"guests","type":"pint","stored":true,"indexed":true}
}' || true

curl -fsS -X POST -H 'Content-Type: application/json' "$schema" -d '{
  "add-field": {"name":"date_time","type":"pdate","stored":true,"indexed":true}
}' || true

curl -fsS -X POST -H 'Content-Type: application/json' "$schema" -d '{
  "add-field": {"name":"meal_type","type":"string","stored":true,"indexed":true}
}' || true

curl -fsS -X POST -H 'Content-Type: application/json' "$schema" -d '{
  "add-field": {"name":"status","type":"string","stored":true,"indexed":true}
}' || true

curl -fsS -X POST -H 'Content-Type: application/json' "$schema" -d '{
  "add-field": {"name":"total_price","type":"pfloat","stored":true,"indexed":true}
}' || true

curl -fsS -X POST -H 'Content-Type: application/json' "$schema" -d '{
  "add-field": {"name":"created_at","type":"pdate","stored":true,"indexed":true}
}' || true

curl -fsS -X POST -H 'Content-Type: application/json' "$schema" -d '{
  "add-field": {"name":"updated_at","type":"pdate","stored":true,"indexed":true}
}' || true

echo "Schema init done."

