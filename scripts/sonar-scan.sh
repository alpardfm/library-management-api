#!/bin/bash
# scripts/sonar-scan.sh
#!/bin/bash

echo "🚀 Starting SonarQube Scan for Library Management API"
echo "===================================================="

# 1. Run all tests with coverage
echo "🧪 Running tests..."
go test ./... -v -coverprofile=coverage.out -covermode=atomic

# 2. Check coverage
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
echo "📈 Current Coverage: ${COVERAGE}%"

# 3. Generate test report
echo "📝 Generating test report..."
go test ./... -json > test-report.json 2>/dev/null || true

# 4. Run Sonar Scanner
echo "🔍 Running SonarQube scan..."
sonar-scanner

# 5. Show results
echo ""
echo "✅ Scan completed!"
echo "📊 Check results at your SonarQube server"
echo ""
echo "🎯 Quality Gate Targets:"
echo "   - Coverage: 100%"
echo "   - Duplications: 0%"
echo "   - Bugs: 0"
echo "   - Vulnerabilities: 0"
echo "   - Code Smells: 0"