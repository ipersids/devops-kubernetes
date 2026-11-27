#!/usr/bin/env bash
set -eu

: "${PGHOST:?Missing PGHOST}"
: "${PGPORT:?Missing PGPORT}"
: "${PGDATABASE:?Missing PGDATABASE}"
: "${PGUSER:?Missing PGUSER}"
: "${PGPASSWORD:?Missing PGPASSWORD}"

pg_dump -v > /tmp/backup.sql

backup_name="backup-$(date +%s).sql"
gcloud storage cp /tmp/backup.sql "gs://project-db-backups/$backup_name"
