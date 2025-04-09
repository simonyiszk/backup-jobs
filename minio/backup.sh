export DUMP_NAME="./export/archive.tar.gz"
mkdir ./export
set -e

rm -rf ./export/*
mkdir ./export/minio-backup
mc alias set backup-source "$MINIO_HOSTNAME" "$MINIO_ACCESS_KEY" "$MINIO_SECRET_KEY"
mc mirror backup-source/ ./export/minio-backup
tar -zcf "$DUMP_NAME" ./export/minio-backup/

./upload-dump
