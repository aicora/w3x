package libraries

import (
	"math/big"

	"github.com/aicora/go-uniswap/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

var (
	// ErrPoolAlreadyInitialized is returned if Initialize is called twice.
	ErrPoolAlreadyInitialized = errors.New("pool already initialized")

	// ErrPoolNotInitialized is returned if state is accessed before initialization.
	ErrPoolNotInitialized = errors.New("pool not initialized")

	// ErrNoLiquidity is returned when liquidity-dependent operations are executed
	// while active liquidity is zero.
	ErrNoLiquidity = errors.New("no liquidity")
)

// Slot0 contains the hot-path pool state.
//
// In Solidity this struct is tightly packed into a single storage slot.
// Here it models the same economic state in Go.
//
// Invariants:
//
//   - SqrtPriceX96 and Tick must be consistent via TickMath
//   - SqrtPriceX96 > 0 once initialized
//   - Fee configuration must respect protocol constraints
//
// SqrtPriceX96 uses Q64.96 fixed-point format.
type Slot0 struct {
    SqrtPriceX96 *big.Int
    Tick         int
    ProtocolFee  utils.ProtocolFee
    LPFee        utils.LPFee
}

// Pool implements the concentrated liquidity AMM state machine.
//
// It maintains:
//
//   1. Global fee accumulators (monotonic increasing)
//   2. Active in-range liquidity
//   3. Tick boundary fee accounting
//
// Financial invariants:
//
//   - feeGrowthGlobal{0,1}X128 are strictly non-decreasing
//   - liquidity >= 0
//   - liquidityGross(T) >= 0 for any tick T
//   - Crossing a tick applies liquidityNet exactly once
//
// Fee accounting follows the model introduced in Uniswap v4.
type Pool struct {
	slot0 Slot0

	// Monotonic global fee growth accumulators (Q128 precision).
	feeGrowthGlobal0X128 *big.Int
	feeGrowthGlobal1X128 *big.Int

	// Active liquidity within current price range.
	liquidity *big.Int

	// Externalized tick storage.
	tickManager ITickManager
}

// NewPool creates an uninitialized pool instance.
func NewPool(tickManager ITickManager) *Pool {
	return &Pool {
		feeGrowthGlobal0X128: big.NewInt(0),
		feeGrowthGlobal1X128: big.NewInt(0),
		liquidity: big.NewInt(0),
		tickManager: tickManager,
	}
}

// Initialize sets the initial price and LP fee.
//
// It computes Tick from sqrtPriceX96 and sets Slot0.
// Can only be executed once.
func (p *Pool) Initialize(sqrtPriceX96 *big.Int, lpFee utils.LPFee) (int, error) {
	if p.slot0.SqrtPriceX96 != nil && p.slot0.SqrtPriceX96.Cmp(big.NewInt(0)) != 0 {
		return 0, ErrPoolAlreadyInitialized
	}

	tick, err := utils.GetTickAtSqrtPrice(sqrtPriceX96)
	if err != nil {
		return 0, err
	}

	p.slot0 = Slot0{
		SqrtPriceX96: sqrtPriceX96,
		Tick:         tick,
		LPFee:        lpFee,
		ProtocolFee:  0,
	}

	return tick, nil
}

// CheckPoolInitialized verifies pool has been initialized.
func (p *Pool) CheckPoolInitialized() error {
	if p.slot0.SqrtPriceX96 == nil || p.slot0.SqrtPriceX96.Cmp(big.NewInt(0)) == 0 {
		return ErrPoolNotInitialized
	}
	return nil
}

// SetProtocolFee updates protocol fee configuration.
func (p *Pool) SetProtocolFee(protocolFee utils.ProtocolFee) error {
	if err := p.CheckPoolInitialized(); err != nil {
		return err
	}

	p.slot0.ProtocolFee = protocolFee
	return nil
}

// SetLPFee updates liquidity provider fee.
func (p *Pool) SetLPFee(lpFee utils.LPFee) error {
	if err := p.CheckPoolInitialized(); err != nil {
		return err
	}

	p.slot0.LPFee = lpFee
	return nil
}

// ClearTick removes tick storage when liquidityGross becomes zero.
func (p *Pool) ClearTick(tick int) {
	p.tickManager.Clear(tick)
}

// GetFeeGrowthInside computes cumulative fee growth inside [tickLower, tickUpper).
//
// Tick axis:
// 
// 0    tickLower        tickUpper     MaxTick
// |-------|-----------------|-------------|
//    10           20              70
//
// | tickCurrent position                 | formula                                         | conceptual result            
// | ------------------------------------ | ----------------------------------------------- | ------------------------- 
// | tickCurrent < tickLower              | lower.FeeGrowthOutside - upper.FeeGrowthOutside | 10 - 70 → negative (not started)   
// | tickLower <= tickCurrent < tickUpper | feeGrowthGlobal - lower - upper                 | 100 - 10 - 70 = 20 (internal accrual) 
// | tickCurrent >= tickUpper             | upper - lower                                   | 70 - 10 = 60 (fully accrued)   
func (p *Pool) GetFeeGrowthInside(tickLower, tickUpper int) (feeGrowthInside0X128, feeGrowthInside1X128 *big.Int) {
	lower := p.tickManager.Get(tickLower)
	upper := p.tickManager.Get(tickUpper)
	tickCurrent := p.slot0.Tick

	switch {
	case tickCurrent < tickLower:
		feeGrowthInside0X128 = new(big.Int).Sub(lower.FeeGrowthOutside0X128, upper.FeeGrowthOutside0X128)
		feeGrowthInside1X128 = new(big.Int).Sub(lower.FeeGrowthOutside1X128, upper.FeeGrowthOutside1X128)

	case tickCurrent >= tickUpper:
		feeGrowthInside0X128 = new(big.Int).Sub(upper.FeeGrowthOutside0X128, lower.FeeGrowthOutside0X128)
		feeGrowthInside1X128 = new(big.Int).Sub(upper.FeeGrowthOutside1X128, lower.FeeGrowthOutside1X128)

	default: // tickLower <= tickCurrent < tickUpper
		feeGrowthInside0X128 = new(big.Int).Sub(p.feeGrowthGlobal0X128, lower.FeeGrowthOutside0X128)
		feeGrowthInside0X128.Sub(feeGrowthInside0X128, upper.FeeGrowthOutside0X128)

		feeGrowthInside1X128 = new(big.Int).Sub(p.feeGrowthGlobal1X128, lower.FeeGrowthOutside1X128)
		feeGrowthInside1X128.Sub(feeGrowthInside1X128, upper.FeeGrowthOutside1X128)
	}

	return
}

// CrossTick executes boundary crossing logic.
//
// Transformation:
//
//   Fo(T) = G - Fo(T)
//
// This flips fee accounting perspective when
// price crosses the boundary.
//
// Returns liquidityNet to apply to active liquidity.
func (p *Pool) CrossTick(tick int, feeGrowthGlobal0X128, feeGrowthGlobal1X128 *big.Int) *big.Int {
    currentTick := p.tickManager.Get(tick)

    // feeGrowthOutside := feeGrowthGlobal - feeGrowthOutside
    currentTick.FeeGrowthOutside0X128 = new(big.Int).Sub(feeGrowthGlobal0X128, currentTick.FeeGrowthOutside0X128)
    currentTick.FeeGrowthOutside1X128 = new(big.Int).Sub(feeGrowthGlobal1X128, currentTick.FeeGrowthOutside1X128)

    return currentTick.LiquidityNet
}

// UpdateTick mutates liquidity state at a boundary.
//
// Steps:
//
// 1. Update liquidityGross
// 2. Detect flip (zero <-> non-zero)
// 3. Initialize feeGrowthOutside if first activation
// 4. Update liquidityNet
//
// liquidityNet semantics:
//   When crossing upward:
//       liquidity += liquidityNet
func (p *Pool) UpdateTick(tick int, liquidityDelta *big.Int, upper bool) (flipped bool, liquidityGrossAfter *big.Int) {
	currentTick := p.tickManager.Get(tick)

	liquidityGrossAfter = new(big.Int).Add(currentTick.LiquidityGross, liquidityDelta)
	flipped = (liquidityGrossAfter.Sign() == 0) != (currentTick.LiquidityGross.Sign() == 0)

	if currentTick.LiquidityGross.Sign() == 0 {
		if tick <= p.slot0.Tick {
			currentTick.FeeGrowthOutside0X128.Set(p.feeGrowthGlobal0X128)
			currentTick.FeeGrowthOutside1X128.Set(p.feeGrowthGlobal1X128)
		}
	}

	if upper {
		currentTick.LiquidityNet.Sub(currentTick.LiquidityNet, liquidityDelta)
	} else {
		currentTick.LiquidityNet.Add(currentTick.LiquidityNet, liquidityDelta)
	}

	currentTick.LiquidityGross.Set(liquidityGrossAfter)

	return
}

// Donate distributes tokens proportionally to all active liquidity.
//
// Formula:
//
//   feeGrowthGlobal += (amount * Q128) / liquidity
//
// Properties:
//
//   - Monotonic fee accumulator
//   - No price movement
//   - No tick mutation
func (p *Pool) Donate(amount0, amount1 *big.Int) (BalanceDelta, error) {
	if p.liquidity.Sign() == 0 {
		return ZeroBalanceDelta, ErrNoLiquidity
	}

	delta := BalanceDelta{
		Amount0: new(big.Int).Neg(amount0),
		Amount1: new(big.Int).Neg(amount1),
	}

	if amount0.Sign() > 0 {
		increment := new(big.Int).Mul(amount0, utils.Q128)
		increment.Div(increment, p.liquidity)
		p.feeGrowthGlobal0X128.Add(p.feeGrowthGlobal0X128, increment)
	}

	if amount1.Sign() > 0 {
		increment := new(big.Int).Mul(amount1, utils.Q128)
		increment.Div(increment, p.liquidity)
		p.feeGrowthGlobal1X128.Add(p.feeGrowthGlobal1X128, increment)
	}

	return delta, nil
}

// ModifyLiquidityParams defines position mutation input.
type ModifyLiquidityParams struct {
	Owner         common.Address
	TickLower     int  
	TickUpper     int   
	LiquidityDelta *big.Int 
	TickSpacing   int   
	Salt          [32]byte
}

// ModifyLiquidityState captures tick mutation results.
type ModifyLiquidityState struct {
	FlippedLower           bool   
	LiquidityGrossAfterLower *big.Int 
	FlippedUpper           bool   
	LiquidityGrossAfterUpper *big.Int
}