package main

import (
	"testing"
)

// TestInitialModel tests that the initial model is created correctly
func TestInitialModel(t *testing.T) {
	model := InitialModel()

	if model.state != MenuState {
		t.Errorf("Expected initial state to be MenuState, got %v", model.state)
	}

	if model.score != 0 {
		t.Errorf("Expected initial score to be 0, got %d", model.score)
	}

	if model.text != "Hello, World!" {
		t.Errorf("Expected text to be 'Hello, World!', got %s", model.text)
	}

	if model.message != "Press SPACE to start, Q to quit" {
		t.Errorf("Expected message to be 'Press SPACE to start, Q to quit', got %s", model.message)
	}
}

// TestHelperFunctions tests the helper functions
func TestHelperFunctions(t *testing.T) {
	// Test min function
	if min(5, 10) != 5 {
		t.Errorf("min(5, 10) should return 5")
	}

	if min(10, 5) != 5 {
		t.Errorf("min(10, 5) should return 5")
	}

	// Test max function
	if max(5, 10) != 10 {
		t.Errorf("max(5, 10) should return 10")
	}

	if max(10, 5) != 10 {
		t.Errorf("max(10, 5) should return 10")
	}

	// Test repeatString function
	result := repeatString("a", 3)
	if result != "aaa" {
		t.Errorf("repeatString('a', 3) should return 'aaa', got %s", result)
	}
}

// TestModelMethods tests the model methods
func TestModelMethods(t *testing.T) {
	model := InitialModel()

	// Test updateWindowSize with valid bounds
	originalX, originalY := model.cursorX, model.cursorY
	model.width = 100
	model.height = 50
	model.updateWindowSize()

	if model.cursorX > model.width-1 || model.cursorY > model.height-1 {
		t.Errorf("cursor should be within window bounds after update")
	}

	// Test that cursor position doesn't go negative
	model.cursorX = -10
	model.cursorY = -10
	model.updateWindowSize()

	if model.cursorX < 0 || model.cursorY < 0 {
		t.Errorf("cursor should not go negative after update")
	}
}

// TestParticleCreation tests particle creation
func TestParticleCreation(t *testing.T) {
	model := InitialModel()
	model.width = 100
	model.height = 50

	// Test creating a particle burst
	originalCount := len(model.particles)
	model.createParticleBurst(10, 10)

	if len(model.particles) != originalCount+10 {
		t.Errorf("Expected 10 new particles, got %d", len(model.particles)-originalCount)
	}
}

// TestColorParticleCreation tests color particle creation
func TestColorParticleCreation(t *testing.T) {
	model := InitialModel()

	originalCount := len(model.background)
	newParticles := createBackgroundParticles()

	if len(newParticles) != 50 {
		t.Errorf("Expected 50 background particles, got %d", len(newParticles))
	}

	// Check that all particles have valid colors
	for _, p := range newParticles {
		if p.color == "" {
			t.Errorf("Particle has empty color")
		}
		if p.size <= 0 {
			t.Errorf("Particle has invalid size: %f", p.size)
		}
	}
}