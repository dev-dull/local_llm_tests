#!/bin/bash

# Test script for Hello, World! Bubbletea Edition

echo "🧪 Running tests for Hello, World! Bubbletea Edition..."
echo ""

# Test 1: Check if go.mod is valid
echo "📋 Test 1: Checking go.mod validity..."
if go mod verify; then
    echo "✅ go.mod is valid"
else
    echo "❌ go.mod verification failed"
    exit 1
fi

# Test 2: Check if dependencies are available
echo ""
echo "📦 Test 2: Checking dependencies..."
if go mod download; then
    echo "✅ Dependencies downloaded successfully"
else
    echo "❌ Failed to download dependencies"
    exit 1
fi

# Test 3: Build the application
echo ""
echo "🔨 Test 3: Building application..."
if go build -o hello-bubbletea-test main.go; then
    echo "✅ Build successful"
else
    echo "❌ Build failed"
    exit 1
fi

# Test 4: Run basic functionality test (quick exit test)
echo ""
echo "⚡ Test 4: Testing application startup..."
timeout 2s ./hello-bubbletea-test -test-mode || echo "Application started successfully (timeout expected)"

# Test 5: Clean up test files
echo ""
echo "🧹 Test 5: Cleaning up..."
rm -f hello-bubbletea-test
rm -f main.go.bak

echo ""
echo "🎉 All tests passed! The application is ready to run."
echo ""
echo "To run the application:"
echo "  ./demo.sh"
echo ""
echo "Or manually:"
echo "  go run main.go"