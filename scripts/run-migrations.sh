#!/bin/bash
set -e

set -a
source .env
set +a

: "${DATABASE_URL:?Need to set DATABASE_URL}"

mkdir -p config

cat > config/dbconfig.yml <<EOF
development:
  dialect: postgres
  datasource: "${DATABASE_URL}"
  dir: db/migrations
EOF

echo "Generated config/dbconfig.yml:"
cat config/dbconfig.yml

$HOME/go/bin/sql-migrate up -config=config/dbconfig.yml