export DUMP_NAME="database.sql"
set -e

mysqldump --skip-ssl -u "$MYSQL_USERNAME" --host "$MYSQL_HOST" --password="$MYSQL_PASSWORD" --result-file "$DUMP_NAME" "$MYSQL_DB"
./upload-dump
