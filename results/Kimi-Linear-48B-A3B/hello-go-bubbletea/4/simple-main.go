package main

import (
	"fmt"
	"time"
)

// Simple version without bubbletea dependencies for testing
func main() {
	fmt.Println("🚀 Hello, World! - Bubbletea Edition")
	fmt.Println("")
	fmt.Println("✨ This is a colorful, animated 'Hello, World!' application")
	fmt.Println("   built with charmbracelet/bubbletea.")
	fmt.Println("")
	fmt.Println("📋 To run the full application:")
	fmt.Println("   1. Install dependencies: go get github.com/charmbracelet/bubbletea")
	fmt.Println("   2. Run: go run main.go")
	fmt.Println("")
	fmt.Println("🎮 Controls:")
	fmt.Println("   • Press 'q' to quit (always works!)")
	fmt.Println("   • Press ' ' to toggle animation")
	fmt.Println("   • Press 'r' to reset")
	fmt.Println("")
	fmt.Println("🎨 Features:")
	fmt.Println("   • Colorful text animation")
	fmt.Println("   • Particle effects")
	fmt.Println("   • Interactive controls")
	fmt.Println("   • Responsive design")
	fmt.Println("")

	// Show a simple animation simulation
	fmt.Println("🎬 Starting animation simulation...")
	for i := 0; i < 5; i++ {
		time.Sleep(500 * time.Millisecond)
		fmt.Printf("   Frame %d/5...\r", i+1)
	}
	fmt.Println("")
	fmt.Println("✅ Animation demo complete!")
	fmt.Println("")
	fmt.Println("🚀 Ready to launch the full bubbletea application!")
}