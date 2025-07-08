#!/bin/bash

# Script sederhana untuk menguji satu kombinasi scaffold

if [ "$#" -lt 5 ]; then
    echo "Usage: $0 <app_type> <router_type> <database_type> <config_type> <feature1,feature2,...>"
    echo "Example: $0 api echo postgresql env basic-auth,sql-migrations"
    exit 1
fi

APP_TYPE=$1
ROUTER_TYPE=$2
DB_TYPE=$3
CONFIG_TYPE=$4
FEATURES=$5
OUTPUT_DIR="/tmp/scaffold-test"

# Hapus direktori uji jika sudah ada
rm -rf "$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"

echo "Generating scaffold for $APP_TYPE-$ROUTER_TYPE-$DB_TYPE with features: $FEATURES"

# Jalankan perintah scaffold
cd /home/witan/personal/scaffold/go-scaffold
./go-scaffold generate \
  --app-type="$APP_TYPE" \
  --router-type="$ROUTER_TYPE" \
  --database-type="$DB_TYPE" \
  --config-type="$CONFIG_TYPE" \
  --log-format="json" \
  --module-path="github.com/example/testproject" \
  --features="$FEATURES" \
  --output="$OUTPUT_DIR"

if [ $? -ne 0 ]; then
    echo "Failed to generate scaffold"
    exit 1
fi

echo "Scaffold generated successfully in $OUTPUT_DIR"

# Build the project
cd "$OUTPUT_DIR"
go mod tidy
go build ./...

if [ $? -eq 0 ]; then
    echo "Build successful!"
else
    echo "Build failed"
    exit 1
fi

echo "Test completed successfully!"
