#!/bin/bash

# Demo script for Hello, World! Bubbletea Edition

echo "🚀 Starting Hello, World! Bubbletea Demo..."
echo ""
echo "This demo will show you the colorful, animated 'Hello, World!' application"
echo "built with charmbracelet/bubbletea."
echo ""
echo "Press Enter to start the demo, or Ctrl+C to cancel..."
read

echo ""
echo "🎨 Starting the application with animation enabled..."
echo ""
echo "Controls:"
echo "  • Press 'q' to quit anytime"
echo "  • Press ' ' (space) to toggle animation"
echo "  • Press 'r' to reset"
echo ""
echo "Enjoy the colorful display! ✨"
echo ""

# Start the application
exec go run main.go