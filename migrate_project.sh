#!/bin/bash

# This script migrates the echo-scaffold project to the new directory structure
# Usage: ./migrate_project.sh

echo "Starting project structure migration..."

# Create backup of old main.go
echo "Creating backup of old main.go"
cp main.go main.go.backup
cp main_test.go main_test.go.backup

# Move Swagger docs to new location
echo "Moving swagger docs"
mkdir -p docs
cp internal/docs/swagger.go docs/swagger.go

# Update imports in code
echo "Updating imports in code"
find . -name "*.go" -type f -exec sed -i 's,github.com/regiwitanto/echo-scaffold/internal/docs,github.com/regiwitanto/echo-scaffold/docs,g' {} \;

# Verify the project builds
echo "Verifying project builds"
go build -o build/echo-scaffold ./cmd/scaffold

# Remove unnecessary files
echo "Removing unnecessary files"
rm -f main.go.backup main_test.go.backup
rm -rf internal/docs

echo "Migration complete!"
echo ""
echo "Testing the application:"
echo "  ./build/echo-scaffold"
echo ""
echo "The project now follows standard Go project layout!"

# Create git commit message
cat > commit_msg.txt << 'EOL'
refactor: reorganize project to follow Go standard layout

- Move main.go to cmd/scaffold/main.go
- Move main_test.go to cmd/scaffold/main_test.go
- Move Swagger docs from internal/docs to docs/
- Update Makefile to use new paths
- Update README.md with new project structure
- Clean up unnecessary files
- Improve project organization following Go standards
EOL

echo ""
echo "Commit message saved to commit_msg.txt. Use it with:"
echo "  git commit -F commit_msg.txt"
