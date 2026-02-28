package libraries

import "math/big"

// TickInfo holds all relevant data for a specific tick in a concentrated liquidity AMM pool.
//
// A "tick" represents a specific price point in the pool. TickInfo stores liquidity and
// fee growth information that allows the pool to efficiently calculate liquidity and fees
// between ticks.
type TickInfo struct {
	// LiquidityGross is the total liquidity that is active at this tick.
	// This value is always >= 0.
	LiquidityGross *big.Int

	// LiquidityNet is the change in liquidity when crossing this tick.
	// Positive when liquidity is added, negative when removed.
	LiquidityNet *big.Int

	// FeeGrowthOutside0X128 is the fee growth of token0 outside (below) this tick,
	// scaled by 2^128. Used to calculate fees for positions that do not span this tick.
	FeeGrowthOutside0X128 *big.Int

	// FeeGrowthOutside1X128 is the fee growth of token1 outside (below) this tick,
	// scaled by 2^128. Used to calculate fees for positions that do not span this tick.
	FeeGrowthOutside1X128 *big.Int
}

// ITickManager defines the interface for a TickManager, which tracks tick data
// and provides efficient lookup and navigation of initialized ticks.
type ITickManager interface {
	// Get returns the TickInfo for a given tick.
	// If the tick does not exist, it initializes a new TickInfo with zero values.
	Get(tick int) *TickInfo

	// Clear clears the TickInfo for a given tick.
	Clear(tick int)

	// IsInitialized checks if the specified tick is initialized in the underlying bitmap.
	IsInitialized(tick int, tickSpacing int) bool

	// FlipTick flips the initialized state of the tick in the underlying bitmap.
	// If the tick was previously uninitialized, it becomes initialized, and vice versa.
	FlipTick(tick int, tickSpacing int) error

	// NextInitializedTickWithinOneWord finds the next initialized tick within a 256-bit word.
	//
	// Parameters:
	// - tick: starting tick
	// - tickSpacing: spacing between valid ticks
	// - lte: search direction; true searches for the next tick <= tick, false for >= tick
	//
	// Returns:
	// - the next initialized tick
	// - a boolean indicating if a tick was found
	// - an error if the search fails
	NextInitializedTickWithinOneWord(tick int, tickSpacing int, lte bool) (int, bool, error)
}

// TickManager manages all ticks in the pool and provides efficient operations for
// checking initialization, flipping ticks, and finding adjacent initialized ticks.
//
// Internally, it stores tick data in a map for fast access and maintains a bitmap
// for efficient navigation of initialized ticks.
type TickManager struct {
	// infos maps a tick index to its TickInfo data.
	infos map[int]*TickInfo

	// bitmap is a helper structure used to efficiently track initialized ticks
	// and navigate between them.
	bitmap *TickBitmap
}

// NewTickManager creates and returns a new TickManager instance.
func NewTickManager() *TickManager {
	return &TickManager{
		infos: make(map[int]*TickInfo),
		bitmap:    NewTickBitmap(),
	}
}

// Get returns the TickInfo for a given tick. If the tick does not exist, it
// creates a new TickInfo with all values initialized to zero.
//
// This ensures that any tick queried always returns a valid TickInfo struct.
func (tm *TickManager) Get(tick int) *TickInfo {
	info, exists := tm.infos[tick]
	if !exists {
		info = &TickInfo{
			LiquidityGross:        big.NewInt(0),
			LiquidityNet:          big.NewInt(0),
			FeeGrowthOutside0X128: big.NewInt(0),
			FeeGrowthOutside1X128: big.NewInt(0),
		}
		tm.infos[tick] = info
	}
	return info
}

// Clear deletes the TickInfo for a given tick.
func (t *TickManager) Clear(tick int) {
	delete(t.infos, tick)
}

// IsInitialized checks if the specified tick is initialized in the underlying bitmap.
func (tm *TickManager) IsInitialized(tick int, tickSpacing int) bool {
	return tm.bitmap.IsInitialized(tick, tickSpacing)
}

// FlipTick flips the initialized state of the specified tick in the underlying bitmap.
//
// A tick is considered "initialized" if it has any liquidity associated with it.
// Flipping changes its state from initialized to uninitialized or vice versa.
func (tm *TickManager) FlipTick(tick int, tickSpacing int) error {
	return tm.bitmap.FlipTick(tick, tickSpacing)
}

// NextInitializedTickWithinOneWord finds the next initialized tick within a single
// 256-bit word of the bitmap.
//
// This function is used to efficiently navigate to the next active tick in either
// direction without scanning every possible tick.
//
// Parameters:
// - tick: the current tick index
// - tickSpacing: the spacing between valid ticks
// - lte: search direction; true = search for next tick <= tick, false = >= tick
//
// Returns:
// - the next initialized tick index
// - a boolean indicating if a tick was found
// - an error if something went wrong
func (tm *TickManager) NextInitializedTickWithinOneWord(tick int, tickSpacing int, lte bool) (int, bool, error) {
	return tm.bitmap.NextInitializedTickWithinOneWord(tick, tickSpacing, lte)
}