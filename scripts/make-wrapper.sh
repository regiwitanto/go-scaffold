#!/bin/bash
# Helper script to transition from Makefile to build.go
# This script detects when users try to use 'make' and suggests using the build.go script instead

# Check if a command was provided
if [ "$#" -eq 0 ]; then
  echo "Error: No command provided"
  echo "Instead of using 'make', please use the Go build system:"
  echo "  go run build.go help"
  exit 1
fi

# Map make commands to build.go commands
command="$1"
build_cmd=""

case "$command" in
  all|build|clean|test|test-unit|test-integration|test-functional|test-cover|run|run-dev|api-docs)
    build_cmd="$command"
    ;;
  *)
    echo "Unknown command: $command"
    echo "Please use the Go build system instead:"
    echo "  go run build.go help"
    exit 1
    ;;
esac

# Execute the corresponding build.go command
echo "Using Go build system instead of Make..."
echo "Executing: go run build.go $build_cmd"
go run build.go "$build_cmd"
