#!/bin/bash

# Assertion Generator Script Wrapper
# Usage: ./scripts/generate_assertions.sh <json_file> [status_code]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Go is installed
if ! command -v go &> /dev/null; then
    print_error "Go is not installed. Please install Go to use this script."
    exit 1
fi

# Check arguments
if [ $# -lt 1 ]; then
    print_error "Usage: $0 <json_file> [status_code]"
    echo ""
    echo "Examples:"
    echo "  $0 response.json"
    echo "  $0 response.json 200"
    echo "  $0 scripts/example_response.json 200"
    exit 1
fi

JSON_FILE="$1"
STATUS_CODE="${2:-}"

# Check if JSON file exists
if [ ! -f "$JSON_FILE" ]; then
    print_error "JSON file '$JSON_FILE' not found."
    exit 1
fi

print_info "Generating assertions from '$JSON_FILE'"
if [ -n "$STATUS_CODE" ]; then
    print_info "Status code: $STATUS_CODE"
fi

echo ""

# Run the Go script
if [ -n "$STATUS_CODE" ]; then
    go run scripts/generate_assertions.go "$JSON_FILE" "$STATUS_CODE"
else
    go run scripts/generate_assertions.go "$JSON_FILE"
fi

echo ""
print_success "Assertion generation completed!"
print_info "Copy the JSON array above and use it in your test spec."
