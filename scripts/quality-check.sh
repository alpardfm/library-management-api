#!/bin/bash
# scripts/quality-check.sh

set -e

echo "🔍 STARTING QUALITY CHECK"
echo "========================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    if [ $2 -eq 0 ]; then
        echo -e "${GREEN}✓ $1${NC}"
    else
        echo -e "${RED}✗ $1${NC}"
        exit 1
    fi
}

# 1. Check if Go is installed
echo "1. Checking Go installation..."
go version
print_status "Go is installed" $?

# 2. Run tests with coverage
echo -e "\n2. Running tests with coverage..."
go test ./... -coverprofile=coverage.out -covermode=atomic
print_status "Tests passed" $?

# 3. Check coverage percentage
echo -e "\n3. Checking coverage percentage..."
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
echo "Current Coverage: ${COVERAGE}%"

if (( $(echo "$COVERAGE >= 90" | bc -l) )); then
    echo -e "${GREEN}✓ Coverage meets minimum 90%${NC}"
else
    echo -e "${RED}✗ Coverage below 90% (current: ${COVERAGE}%)${NC}"
    exit 1
fi

# 4. Run linter
echo -e "\n4. Running golangci-lint..."
if command -v golangci-lint &> /dev/null; then
    golangci-lint run
    print_status "Linting passed" $?
else
    echo -e "${YELLOW}⚠ golangci-lint not installed, skipping...${NC}"
fi

# 5. Check for race conditions
echo -e "\n5. Checking for race conditions..."
go test ./... -race
print_status "No race conditions detected" $?

# 6. Generate test report
echo -e "\n6. Generating test report..."
go test ./... -json > test-report.json 2>/dev/null || true
print_status "Test report generated" $?

# 7. Check for unused dependencies
echo -e "\n7. Checking for unused dependencies..."
go mod tidy
git diff --exit-code go.mod go.sum
print_status "No unused dependencies" $?

# 8. Security check (basic)
echo -e "\n8. Running basic security checks..."
# Check for hardcoded secrets (simplified)
if grep -r "password.*=.*\"[^\"]*[a-zA-Z0-9]{8,}\"" --include="*.go" . | grep -v "_test.go" | grep -v "test"; then
    echo -e "${RED}✗ Potential hardcoded password found${NC}"
    exit 1
fi
print_status "No obvious security issues found" 0

# 9. SonarQube scan (optional)
echo -e "\n9. Running SonarQube scan..."
if command -v sonar-scanner &> /dev/null; then
    sonar-scanner
    print_status "SonarQube scan completed" $?
else
    echo -e "${YELLOW}⚠ sonar-scanner not installed, skipping...${NC}"
fi

echo -e "\n${GREEN}🎉 ALL QUALITY CHECKS PASSED!${NC}"
echo "📊 Coverage: ${COVERAGE}%"
echo "🚀 Ready for deployment!"