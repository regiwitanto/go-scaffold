#!/bin/bash

# run-tests.sh - Script to run all or specific tests for go-scaffold
# This is a helper script that makes it easier to run tests.

# Set colors for output
GREEN="\033[0;32m"
RED="\033[0;31m"
YELLOW="\033[0;33m"
BLUE="\033[0;34m"
NC="\033[0m" # No Color

# Set default test directory
PROJECT_ROOT=$(pwd)
TEST_DIR="$PROJECT_ROOT/test"

# Function to print usage information
print_usage() {
    echo -e "${YELLOW}Go-Scaffold Test Runner${NC}"
    echo ""
    echo -e "Usage: $0 [options]"
    echo ""
    echo -e "Options:"
    echo -e "  ${BLUE}-a, --all${NC}            Run all tests"
    echo -e "  ${BLUE}-u, --unit${NC}           Run unit tests only"
    echo -e "  ${BLUE}-i, --integration${NC}    Run integration tests only"
    echo -e "  ${BLUE}-f, --functional${NC}     Run functional tests only"
    echo -e "  ${BLUE}-b, --benchmark${NC}      Run benchmark tests"
    echo -e "  ${BLUE}-c, --cover${NC}          Generate coverage report"
    echo -e "  ${BLUE}-v, --verbose${NC}        Run tests in verbose mode"
    echo -e "  ${BLUE}-h, --help${NC}           Show this help message"
    echo ""
    echo -e "Examples:"
    echo -e "  $0 --all           # Run all tests"
    echo -e "  $0 --unit          # Run unit tests only"
    echo -e "  $0 --cover --unit  # Run unit tests with coverage"
    echo ""
}

# Parse command line arguments
if [ $# -eq 0 ]; then
    print_usage
    exit 0
fi

# Default options
RUN_ALL=false
RUN_UNIT=false
RUN_INTEGRATION=false
RUN_FUNCTIONAL=false
RUN_BENCHMARK=false
GENERATE_COVERAGE=false
VERBOSE=false

# Process command line arguments
while [[ $# -gt 0 ]]; do
    key="$1"
    case $key in
        -a|--all)
            RUN_ALL=true
            shift
            ;;
        -u|--unit)
            RUN_UNIT=true
            shift
            ;;
        -i|--integration)
            RUN_INTEGRATION=true
            shift
            ;;
        -f|--functional)
            RUN_FUNCTIONAL=true
            shift
            ;;
        -b|--benchmark)
            RUN_BENCHMARK=true
            shift
            ;;
        -c|--cover)
            GENERATE_COVERAGE=true
            shift
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -h|--help)
            print_usage
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $key${NC}"
            print_usage
            exit 1
            ;;
    esac
done

# If no specific test types are selected but --all is not specified, show usage
if [ "$RUN_ALL" = false ] && [ "$RUN_UNIT" = false ] && [ "$RUN_INTEGRATION" = false ] && [ "$RUN_FUNCTIONAL" = false ] && [ "$RUN_BENCHMARK" = false ]; then
    echo -e "${RED}No test type selected.${NC}"
    print_usage
    exit 1
fi

# Set verbose flag
VERBOSE_FLAG=""
if [ "$VERBOSE" = true ]; then
    VERBOSE_FLAG="-v"
fi

# Set coverage flag
COVERAGE_FLAGS=""
COVER_PROFILE=""
if [ "$GENERATE_COVERAGE" = true ]; then
    echo -e "${YELLOW}Generating coverage report...${NC}"
    COVERAGE_DIR="$PROJECT_ROOT/coverage"
    mkdir -p "$COVERAGE_DIR"
    COVER_PROFILE="$COVERAGE_DIR/coverage.out"
    COVERAGE_FLAGS="-coverprofile=$COVER_PROFILE"
fi

# Run tests
if [ "$RUN_ALL" = true ] || [ "$RUN_UNIT" = true ]; then
    echo -e "${YELLOW}Running unit tests...${NC}"
    go test $VERBOSE_FLAG $COVERAGE_FLAGS "$TEST_DIR/unit/..."
    UNIT_EXIT_CODE=$?
    if [ $UNIT_EXIT_CODE -ne 0 ]; then
        echo -e "${RED}Unit tests failed with exit code $UNIT_EXIT_CODE${NC}"
    else
        echo -e "${GREEN}Unit tests passed!${NC}"
    fi
fi

if [ "$RUN_ALL" = true ] || [ "$RUN_INTEGRATION" = true ]; then
    echo -e "${YELLOW}Running integration tests...${NC}"
    echo -e "${BLUE}Note: Integration tests are currently skipped.${NC}"
    echo -e "${BLUE}See /test/integration/README.md for more information.${NC}"
    go test $VERBOSE_FLAG $COVERAGE_FLAGS "$TEST_DIR/integration/..."
    INTEGRATION_EXIT_CODE=$?
    if [ $INTEGRATION_EXIT_CODE -ne 0 ]; then
        echo -e "${RED}Integration tests failed with exit code $INTEGRATION_EXIT_CODE${NC}"
    else
        echo -e "${GREEN}Integration tests passed!${NC}"
    fi
fi

if [ "$RUN_ALL" = true ] || [ "$RUN_FUNCTIONAL" = true ]; then
    echo -e "${YELLOW}Running functional tests...${NC}"
    
    # Check if port 8081 is already in use
    if (echo > /dev/tcp/127.0.0.1/8081) 2>/dev/null; then
        echo -e "${RED}Port 8081 is already in use. Functional tests will fail.${NC}"
        echo -e "${BLUE}Please free port 8081 before running functional tests.${NC}"
        if [ "$RUN_ALL" = true ]; then
            echo -e "${BLUE}Skipping functional tests...${NC}"
            FUNCTIONAL_EXIT_CODE=0
        else
            # Only show the warning but continue if specifically running functional tests
            echo -e "${BLUE}Continuing anyway since you specifically requested functional tests...${NC}"
            go test $VERBOSE_FLAG $COVERAGE_FLAGS "$TEST_DIR/functional/..."
            FUNCTIONAL_EXIT_CODE=$?
        fi
    else
        go test $VERBOSE_FLAG $COVERAGE_FLAGS "$TEST_DIR/functional/..."
        FUNCTIONAL_EXIT_CODE=$?
    fi
    
    if [ $FUNCTIONAL_EXIT_CODE -ne 0 ]; then
        echo -e "${RED}Functional tests failed with exit code $FUNCTIONAL_EXIT_CODE${NC}"
    else
        echo -e "${GREEN}Functional tests passed!${NC}"
    fi
fi

if [ "$RUN_BENCHMARK" = true ]; then
    echo -e "${YELLOW}Running benchmark tests...${NC}"
    go test -bench=. "$TEST_DIR/unit/benchmark_test.go"
    BENCHMARK_EXIT_CODE=$?
    if [ $BENCHMARK_EXIT_CODE -ne 0 ]; then
        echo -e "${RED}Benchmark tests failed with exit code $BENCHMARK_EXIT_CODE${NC}"
    else
        echo -e "${GREEN}Benchmark tests complete!${NC}"
    fi
fi

# Generate HTML coverage report if coverage was requested
if [ "$GENERATE_COVERAGE" = true ] && [ -f "$COVER_PROFILE" ]; then
    COVERAGE_HTML="$COVERAGE_DIR/coverage.html"
    echo -e "${YELLOW}Generating HTML coverage report at $COVERAGE_HTML${NC}"
    go tool cover -html="$COVER_PROFILE" -o "$COVERAGE_HTML"
    echo -e "${GREEN}Coverage report generated.${NC}"
    
    # Display coverage percentage
    COVERAGE=$(go tool cover -func="$COVER_PROFILE" | grep total | awk '{print $3}')
    echo -e "${BLUE}Total coverage: $COVERAGE${NC}"
fi

echo -e "${GREEN}All tests complete!${NC}"
