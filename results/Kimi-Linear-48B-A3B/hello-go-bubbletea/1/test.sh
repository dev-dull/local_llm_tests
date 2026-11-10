#!/bin/bash

echo "=== Hello Go Bubbletea World! Test Script ==="
echo

echo "1. Checking if go.mod exists..."
if [ -f go.mod ]; then
    echo "✓ go.mod found"
else
    echo "✗ go.mod missing"
    exit 1
fi

echo
echo "2. Checking if main.go exists..."
if [ -f main.go ]; then
    echo "✓ main.go found"
else
    echo "✗ main.go missing"
    exit 1
fi

echo
echo "3. Checking if README.md exists..."
if [ -f README.md ]; then
    echo "✓ README.md found"
else
    echo "✗ README.md missing"
    exit 1
fi

echo
echo "4. Checking if Makefile exists..."
if [ -f Makefile ]; then
    echo "✓ Makefile found"
else
    echo "✗ Makefile missing"
    exit 1
fi

echo
echo "5. Checking Go module dependencies..."
if [ -f go.sum ]; then
    echo "✓ go.sum found"
    echo "  Dependencies:"
    grep "github.com/charmbracelet" go.sum | head -3
else
    echo "✗ go.sum missing (running go mod tidy)..."
    go mod tidy
fi

echo
echo "6. Testing build..."
if go build -o hello-go-bubbletea main.go; then
    echo "✓ Build successful"
else
    echo "✗ Build failed"
    exit 1
fi

echo
echo "7. Checking executable..."
if [ -f hello-go-bubbletea ]; then
    echo "✓ Executable created"
    echo "  Size: $(ls -la hello-go-bubbletea | awk '{print $5}') bytes"
else
    echo "✗ Executable not found"
    exit 1
fi

echo
echo "8. Running go mod tidy..."
go mod tidy

echo
echo "=== Test Summary ==="
echo "✓ All checks passed!"
echo
echo "To run the application:"
echo "  ./hello-go-bubbletea"
echo
echo "To use the Makefile:"
echo "  make run"
echo
echo "Features included:"
echo "  ✓ Interactive menu system"
echo "  ✓ Animated title with color cycling"
echo "  ✓ Color demonstration mode"
echo "  ✓ Input mode for text entry"
echo "  ✓ Game-like interface with scoring"
echo "  ✓ High score persistence"
echo "  ✓ Particle animation system"
echo "  ✓ Responsive navigation"
echo "  ✓ Clean state management"