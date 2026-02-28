package libraries

import (
	"math/big"
	"testing"
)

func TestTickManager_Get(t *testing.T) {
	tm := NewTickManager()

	tick := int(10)
	info := tm.Get(tick)
	if info == nil {
		t.Fatal("Get returned nil")
	}
	if info.LiquidityGross.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("expected LiquidityGross 0, got %s", info.LiquidityGross.String())
	}
	if info.LiquidityNet.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("expected LiquidityNet 0, got %s", info.LiquidityNet.String())
	}
	if info.FeeGrowthOutside0X128.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("expected FeeGrowthOutside0X128 0, got %s", info.FeeGrowthOutside0X128.String())
	}
	if info.FeeGrowthOutside1X128.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("expected FeeGrowthOutside1X128 0, got %s", info.FeeGrowthOutside1X128.String())
	}

	info2 := tm.Get(tick)
	if info2 != info {
		t.Errorf("expected same TickInfo instance")
	}
}

func TestTickManager_IsInitialized(t *testing.T) {
	tm := NewTickManager()
	tick := int(10)
	tickSpacing := int(1)
	if tm.IsInitialized(tick, tickSpacing) {
		t.Fatal("tick should not be initialized initially")
	}
}

func TestTickManager_FlipTick(t *testing.T) {
	tm := NewTickManager()
	tick := int(10)
	tickSpacing := int(1)

	found := tm.IsInitialized(tick, tickSpacing)
	if found {
		t.Fatal("tick should not be initialized initially")
	}

	if err := tm.FlipTick(tick, tickSpacing); err != nil {
		t.Fatalf("FlipTick error: %v", err)
	}

	found = tm.IsInitialized(tick, tickSpacing)
	if !found {
		t.Fatal("tick should be initialized after FlipTick")
	}

	if err := tm.FlipTick(tick, tickSpacing); err != nil {
		t.Fatalf("FlipTick error: %v", err)
	}
	found = tm.IsInitialized(tick, tickSpacing)
	if found {
		t.Fatal("tick should be uninitialized after second FlipTick")
	}
}

func TestTickManager_NextInitializedTickWithinOneWord(t *testing.T) {
	tm := NewTickManager()
	tickSpacing := int(1)

	ticks := []int{5, 10, 15}
	for _, tk := range ticks {
		_ = tm.FlipTick(tk, tickSpacing)
	}

	next, found, err := tm.NextInitializedTickWithinOneWord(7, tickSpacing, false)
	if err != nil {
		t.Fatalf("NextInitializedTickWithinOneWord error: %v", err)
	}
	if !found || next != 10 {
		t.Errorf("expected next initialized tick 10, got %d", next)
	}

	prev, found, err := tm.NextInitializedTickWithinOneWord(12, tickSpacing, true)
	if err != nil {
		t.Fatalf("NextInitializedTickWithinOneWord error: %v", err)
	}
	if !found || prev != 10 {
		t.Errorf("expected previous initialized tick 10, got %d", prev)
	}

	_, found, err = tm.NextInitializedTickWithinOneWord(3, tickSpacing, true)
	if err != nil {
		t.Fatalf("NextInitializedTickWithinOneWord error: %v", err)
	}
	if found {
		t.Errorf("expected no initialized tick found")
	}
}