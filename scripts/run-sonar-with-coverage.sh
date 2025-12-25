#!/bin/bash
# scripts/run-sonar-with-coverage.sh

#!/bin/bash

echo "🚀 SonarQube Scan with Coverage Report"
echo "======================================"

# Configuration
PROJECT_KEY="library-management-api"
PROJECT_NAME="Library Management API"
SONAR_HOST="http://localhost:9000"  # Change to your SonarQube URL
SONAR_TOKEN="sqa_0b87942d609943eb4732172669a9e6902401b71d"  # From environment

if [ -z "$SONAR_TOKEN" ]; then
    echo "❌ SONAR_TOKEN environment variable is required"
    exit 1
fi

# 1. Clean previous reports
echo "🧹 Cleaning previous reports..."
rm -f coverage.out coverage.html test-report.json

# 2. Run tests with coverage
echo "🧪 Running tests..."
go test ./... -v -coverprofile=coverage.out -covermode=atomic

# 3. Generate test report
echo "📝 Generating test report..."
go test ./... -json > test-report.json 2>/dev/null || true

# 4. Calculate coverage
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
echo "📈 Current Coverage: ${COVERAGE}%"

# 5. Run Sonar Scanner
echo "🔍 Running SonarQube scan..."
sonar-scanner \
  -Dsonar.projectKey=$PROJECT_KEY \
  -Dsonar.projectName="$PROJECT_NAME" \
  -Dsonar.host.url=$SONAR_HOST \
  -Dsonar.login=$SONAR_TOKEN \
  -Dsonar.go.coverage.reportPaths=coverage.out \
  -Dsonar.go.tests.reportPaths=test-report.json \
  -Dsonar.sources=. \
  -Dsonar.exclusions="**/*_test.go,**/vendor/**,**/mocks/**,**/testdata/**" \
  -Dsonar.tests=. \
  -Dsonar.test.inclusions="**/*_test.go" \
  -Dsonar.qualitygate.wait=true \
  -Dsonar.qualitygate.timeout=600

# 6. Display results
echo ""
echo "✅ Scan completed!"
echo "🌐 Open ${SONAR_HOST}/dashboard?id=${PROJECT_KEY} to view results"
echo ""
echo "🎯 Quality Gate Status:"
echo "   - Coverage: ${COVERAGE}%"
echo "   - Target: 100%"
echo ""
if (( $(echo "$COVERAGE >= 100" | bc -l) )); then
    echo "🎉 PERFECT COVERAGE ACHIEVED!"
else
    echo "📈 Need to improve coverage by $((100 - COVERAGE))%"
fi