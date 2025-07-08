#!/bin/bash

# DEPRECATED: This script has been replaced by Go-based tests.
# Please use the new test suite instead:
#   make test       # Run all tests
#   make test-unit  # Run unit tests only
#   make test-integration  # Run integration tests only
#   make test-functional   # Run functional tests only

# Set warna untuk output
GREEN="\033[0;32m"
RED="\033[0;31m"
YELLOW="\033[0;33m"
NC="\033[0m" # No Color

echo -e "${YELLOW}WARNING: This script is deprecated and will be removed in a future version.${NC}"
echo -e "${YELLOW}Please use the Go-based tests instead:${NC}"
echo -e "${GREEN}  make test${NC}            # Run all tests"
echo -e "${GREEN}  make test-unit${NC}       # Run unit tests only"
echo -e "${GREEN}  make test-integration${NC}# Run integration tests only"
echo -e "${GREEN}  make test-functional${NC} # Run functional tests only"
echo ""
read -p "Do you want to continue with this deprecated script anyway? (y/n) " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${RED}Exiting.${NC}"
    exit 1
fi
echo -e "${YELLOW}Continuing with deprecated script...${NC}"
echo ""

# Array untuk menyimpan hasil
results=()
failed_count=0
success_count=0

# Fungsi untuk menguji satu kombinasi
test_scaffold() {
  local app_type=$1
  local router_type=$2
  local db_type=$3
  local config_type=$4
  local features=$5
  local test_name="$app_type-$router_type-$db_type-$config_type"
  local test_dir="/tmp/test-scaffold-$test_name"
  
  echo -e "${YELLOW}Testing $test_name with features: $features${NC}"
  
  # Hapus direktori test jika sudah ada
  rm -rf "$test_dir"
  rm -f "$test_dir.zip" "$test_dir.zip.tmp"
  
  echo "Generating scaffold..."
  
  # Mode API atau CLI
  if [ "$test_mode" == "1" ]; then
    # Buat scaffold dengan curl ke API lokal
    response=$(curl -s -w "\n%{http_code}" -X POST "http://localhost:3030/api/v1/scaffold" \
      -H "Content-Type: application/json" \
      -d "{
        \"appType\": \"$app_type\",
        \"routerType\": \"$router_type\",
        \"databaseType\": \"$db_type\",
        \"configType\": \"$config_type\",
        \"logFormat\": \"json\",
        \"modulePath\": \"github.com/example/testapp\",
        \"features\": [$(echo "$features" | sed 's/,/","/g' | sed 's/^/"/' | sed 's/$/"/')]
      }" \
      -o "$test_dir.zip.tmp")
    
    # Ambil status code dari response
    http_code=$(echo "$response" | tail -n1)
    
    # Periksa status code
    if [ "$http_code" != "200" ]; then
      results+=("${RED}❌ $test_name with $features: API returned status $http_code${NC}")
      ((failed_count++))
      return
    fi
    
    # Pindahkan file temporary ke file zip final
    mv "$test_dir.zip.tmp" "$test_dir.zip"
  else
    # Gunakan CLI tool langsung
    feature_args=""
    IFS=',' read -ra FEATURE_ARRAY <<< "$features"
    for i in "${FEATURE_ARRAY[@]}"; do
      feature_args="$feature_args --feature $i"
    done
    
    $GO_SCAFFOLD_BIN generate \
      --app-type="$app_type" \
      --router-type="$router_type" \
      --database-type="$db_type" \
      --config-type="$config_type" \
      --log-format="json" \
      --module-path="github.com/example/testapp" \
      $feature_args \
      --output="$test_dir.zip"
    
    if [ $? -ne 0 ]; then
      results+=("${RED}❌ $test_name with $features: CLI command failed${NC}")
      ((failed_count++))
      return
    fi
  fi
  
  # Periksa apakah zip file valid
  if ! unzip -t "$test_dir.zip" > /dev/null 2>&1; then
    results+=("${RED}❌ $test_name with $features: Generated ZIP file is invalid${NC}")
    ((failed_count++))
    return
  fi
  
  # Buat direktori dan unzip
  mkdir -p "$test_dir"
  unzip -q "$test_dir.zip" -d "$test_dir"
  
  # Cek apakah direktori codebase ada
  if [ ! -d "$test_dir/codebase" ]; then
    echo "Directory structure not as expected. Checking for alternative paths..."
    
    # Temukan direktori yang mungkin berisi proyek
    project_dir=$(find "$test_dir" -type f -name "go.mod" | head -1 | xargs dirname 2>/dev/null)
    
    if [ -z "$project_dir" ]; then
      results+=("${RED}❌ $test_name with $features: Cannot find project directory${NC}")
      ((failed_count++))
      return
    fi
    
    echo "Found project at $project_dir"
    cd "$project_dir"
  else
    # Pindah ke direktori project seperti yang diharapkan
    cd "$test_dir/codebase"
  fi
  
  # Build project
  echo "Building project..."
  go mod tidy
  go build ./...
  build_status=$?
  
  if [ $build_status -eq 0 ]; then
    results+=("${GREEN}✅ $test_name with $features: BUILD SUCCESS${NC}")
    ((success_count++))
  else
    results+=("${RED}❌ $test_name with $features: BUILD FAILED${NC}")
    ((failed_count++))
  fi
}

# Tentukan mode pengujian: API atau CLI
echo -e "${YELLOW}Select test mode:${NC}"
echo "1. API mode (requires API server running on http://localhost:3030)"
echo "2. CLI mode (uses go-scaffold command line tool directly)"
read -p "Enter your choice (1 or 2): " test_mode

if [ "$test_mode" != "1" ] && [ "$test_mode" != "2" ]; then
  echo -e "${RED}Invalid choice. Exiting.${NC}"
  exit 1
fi

# Jika mode API, pastikan server berjalan
if [ "$test_mode" == "1" ]; then
  echo -e "${YELLOW}Checking if API server is running...${NC}"
  
  # Coba beberapa endpoint umum untuk memeriksa apakah server berjalan
  if curl -s -o /dev/null "http://localhost:3030/" 2>/dev/null || 
     curl -s -o /dev/null "http://localhost:3030/health" 2>/dev/null ||
     curl -s -o /dev/null "http://localhost:3030/api/health" 2>/dev/null ||
     curl -s -o /dev/null "http://localhost:3030/healthz" 2>/dev/null; then
    echo -e "${GREEN}API server appears to be running.${NC}"
  else
    echo -e "${YELLOW}Warning: Could not verify API server is running on http://localhost:3030.${NC}"
    echo -e "${YELLOW}Do you want to continue anyway? (y/n)${NC}"
    read -p "" continue_anyway
    if [ "$continue_anyway" != "y" ] && [ "$continue_anyway" != "Y" ]; then
      echo -e "${RED}Exiting.${NC}"
      exit 1
    fi
    echo -e "${YELLOW}Continuing with API mode...${NC}"
  fi
else
  echo -e "${YELLOW}Using CLI mode.${NC}"
  # Check if go-scaffold binary exists
  if ! command -v go-scaffold &> /dev/null; then
    echo -e "${RED}go-scaffold command not found. Please make sure it's in your PATH.${NC}"
    # Try to find it in the workspace
    if [ -x "../build/go-scaffold" ]; then
      echo -e "${YELLOW}Found go-scaffold in ../build, using that.${NC}"
      GO_SCAFFOLD_BIN="../build/go-scaffold"
    elif [ -x "./build/go-scaffold" ]; then
      echo -e "${YELLOW}Found go-scaffold in ./build, using that.${NC}"
      GO_SCAFFOLD_BIN="./build/go-scaffold"
    else
      echo -e "${RED}Could not find go-scaffold binary. Exiting.${NC}"
      exit 1
    fi
  else
    GO_SCAFFOLD_BIN="go-scaffold"
  fi
fi

# Periksa apakah ada parameter command line untuk menjalankan test tunggal
if [ "$#" -eq 5 ]; then
  echo -e "${YELLOW}Running single test with specified parameters${NC}"
  test_mode=2  # Default ke CLI mode untuk test tunggal
  
  if [ "$1" == "--api" ]; then
    test_mode=1  # Gunakan API mode jika diminta
    shift
  fi
  
  if [ "$#" -ne 4 ]; then
    echo -e "${RED}Invalid parameters. Usage: $0 [--api] <app_type> <router_type> <database_type> <config_type> <features>${NC}"
    exit 1
  fi
  
  app_type=$1
  router_type=$2
  database_type=$3
  config_type=$4
  features=$5
  
  if [ "$test_mode" == "2" ]; then
    # Setup CLI mode
    if ! command -v go-scaffold &> /dev/null; then
      if [ -x "./build/go-scaffold" ]; then
        GO_SCAFFOLD_BIN="./build/go-scaffold"
      elif [ -x "../build/go-scaffold" ]; then
        GO_SCAFFOLD_BIN="../build/go-scaffold"
      else
        echo -e "${RED}go-scaffold binary not found. Please build it first.${NC}"
        exit 1
      fi
    else
      GO_SCAFFOLD_BIN="go-scaffold"
    fi
  fi
  
  # Jalankan test tunggal
  test_scaffold "$app_type" "$router_type" "$database_type" "$config_type" "$features"
  
  # Tampilkan hasil
  echo -e "\n${YELLOW}=== TEST RESULTS ===${NC}\n"
  for result in "${results[@]}"; do
    echo -e "$result"
  done
  
  echo -e "\n${YELLOW}=== SUMMARY ===${NC}"
  echo -e "${GREEN}Successful tests: $success_count${NC}"
  echo -e "${RED}Failed tests: $failed_count${NC}"
  
  exit 0
fi

echo "Starting tests..."

# Pengujian berbasis matriks (subset)
# 1. API templates
test_scaffold "api" "echo" "postgresql" "env" "basic-auth,sql-migrations"
test_scaffold "api" "chi" "mysql" "env" "basic-auth,email"
test_scaffold "api" "gin" "sqlite" "flags" "access-logging,sql-migrations"
test_scaffold "api" "standard" "none" "flags" "gitignore,admin-makefile"

# 2. Webapp templates
test_scaffold "webapp" "echo" "postgresql" "env" "sql-migrations,email"
test_scaffold "webapp" "chi" "mysql" "flags" "access-logging,basic-auth"
test_scaffold "webapp" "gin" "sqlite" "env" "live-reload,secure-cookies"
test_scaffold "webapp" "standard" "none" "flags" "automatic-versioning"

# 3. Database specific tests
test_scaffold "api" "echo" "postgresql" "env" "sql-migrations"
test_scaffold "api" "echo" "mysql" "env" "sql-migrations"
test_scaffold "api" "echo" "sqlite" "env" "sql-migrations"

# 4. Feature combination tests
test_scaffold "api" "chi" "postgresql" "flags" "email,error-notifications"
test_scaffold "webapp" "echo" "mysql" "env" "basic-auth,secure-cookies"

# Tampilkan hasil
echo -e "\n${YELLOW}=== TEST RESULTS ===${NC}\n"
for result in "${results[@]}"; do
  echo -e "$result"
done

echo -e "\n${YELLOW}=== SUMMARY ===${NC}"
echo -e "${GREEN}Successful tests: $success_count${NC}"
echo -e "${RED}Failed tests: $failed_count${NC}"

# Hapus file zip sementara
find /tmp -name "test-scaffold-*.zip" -delete

echo -e "\n${YELLOW}Test complete.${NC}"
