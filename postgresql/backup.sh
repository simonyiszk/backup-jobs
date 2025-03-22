export DUMP_NAME="./database.sql"
set -e

pg_dump -U "$PGUSER" "$DBNAME" -f "$DUMP_NAME"

./upload-dump
