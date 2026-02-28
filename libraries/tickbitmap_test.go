package libraries

import (
	"testing"
)

func TestTickBitmap_FlipAndIsInitialized(t *testing.T) {
	tb := NewTickBitmap()
	tick := int(10)
	tickSpacing := int(1)

	if tb.IsInitialized(tick, tickSpacing) {
		t.Fatal("tick should not be initialized initially")
	}

	if err := tb.FlipTick(tick, tickSpacing); err != nil {
		t.Fatalf("FlipTick error: %v", err)
	}

	if !tb.IsInitialized(tick, tickSpacing) {
		t.Fatal("tick should be initialized after FlipTick")
	}

	if err := tb.FlipTick(tick, tickSpacing); err != nil {
		t.Fatalf("FlipTick error: %v", err)
	}
	if tb.IsInitialized(tick, tickSpacing) {
		t.Fatal("tick should be uninitialized after second FlipTick")
	}
}

func TestTickBitmap_FlipTick(t *testing.T) {
	tb := NewTickBitmap()
	tickSpacing := int(10)
	tick := int(20)

	// Flip tick to initialized
	if err := tb.FlipTick(tick, tickSpacing); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Flip tick again to uninitialized
	if err := tb.FlipTick(tick, tickSpacing); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check uninitialized
	compressed := tick / tickSpacing
	wordPos := int16(compressed >> 8)
	if word, ok := tb.words[wordPos]; ok {
		if word.Sign() != 0 {
			t.Errorf("expected word to be zero after flipping back, got %v", word)
		}
	}
}

func TestNextInitializedTickWithinOneWord(t *testing.T) {
	tb := NewTickBitmap()
	tickSpacing := int(10)

	// Initialize some ticks
	tb.FlipTick(20, tickSpacing)
	tb.FlipTick(50, tickSpacing)
	tb.FlipTick(200, tickSpacing)

	// Search lte
	next, initialized, err := tb.NextInitializedTickWithinOneWord(45, tickSpacing, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !initialized || next != 20 {
		t.Errorf("expected 20, got %d", next)
	}

	// Search gt
	next, initialized, err = tb.NextInitializedTickWithinOneWord(45, tickSpacing, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !initialized || next != 50 {
		t.Errorf("expected 50, got %d", next)
	}
}