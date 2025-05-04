export DUMP_NAME="./export/archive.tar.gz"
mkdir ./export
set -e

rm -rf ./export/*

cd /mounted-folder/
tar -zcf "$OLDPWD/$DUMP_NAME" ./*

cd -
./upload-dump
