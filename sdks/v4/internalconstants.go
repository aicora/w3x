package sdkv4

import (
	"math/big"

	sdkv3 "github.com/aicora/go-uniswap/sdks/v3"
	"github.com/ethereum/go-ethereum/common"
)

// AddressZero represents the canonical zero address in Ethereum (0x000...000).
var AddressZero = common.Address{}

// BigInt constants commonly used in arithmetic and token math.
var (
	NegativeOne = big.NewInt(-1)                                        // -1 as a big.Int
	Zero        = big.NewInt(0)                                         // 0 as a big.Int
	One         = big.NewInt(1)                                         // 1 as a big.Int
	OneEther    = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil) // 10^18, representing 1 ETH in wei
)

// EmptyBytes is a zero-length hex string placeholder.
const EmptyBytes = "0x"

// Q96 and Q192 are used in fixed-point arithmetic for liquidity calculations.
// Q96 = 2^96, Q192 = 2^192 (used in Uniswap V3 sqrt price computations).
var (
	Q96  = new(big.Int).Exp(big.NewInt(2), big.NewInt(96), nil)
	Q192 = new(big.Int).Exp(Q96, big.NewInt(2), nil)
)

// Pool fee tiers and default tick spacings.
const (
	FeeAmountLow     = 100   // Low fee tier (0.01%)
	FeeAmountMedium  = 3000  // Medium fee tier (0.3%)
	FeeAmountHighest = 10000 // Highest fee tier (1%)

	TickSpacingTen   = 10 // Tick spacing for low fees
	TickSpacingSixty = 60 // Tick spacing for medium fees
)

// MinSlippageDecrease represents the minimum slippage allowed when decreasing positions.
const MinSlippageDecrease = 0

// OpenDelta is used when unwrapping WETH to ETH in position manager logic.
var OpenDelta = big.NewInt(0)

// SqrtPrice1_1 is the default square root price (1:1) encoded in Q64.96 format.
// Utilizes the EncodeSqrtRatioX96 helper from the v3 SDK.
var SqrtPrice1_1 = sdkv3.EncodeSqrtRatioX96(big.NewInt(1), big.NewInt(1))

// EmptyHook represents a default no-op hook address.
var EmptyHook = common.HexToAddress("0x0000000000000000000000000000000000000000")

// Standard error constants used throughout the SDK.
const (
	ErrNativeNotSet  = "NATIVE_NOT_SET" // Native token not configured
	ErrZeroLiquidity = "ZERO_LIQUIDITY" // Operation would result in zero liquidity
	ErrNoSqrtPrice   = "NO_SQRT_PRICE"  // Pool has no initialized price
	ErrCannotBurn    = "CANNOT_BURN"    // Attempt to burn non-existent position
)

// PositionFunction represents the function selectors for the PositionManager contract.
type PositionFunction string

const (
	InitializePool    PositionFunction = "initializePool"    // Initialize a new pool with initial liquidity
	ModifyLiquidities PositionFunction = "modifyLiquidities" // Increase or decrease liquidity for a position

	// Inherited from PermitForwarder
	PermitBatch PositionFunction = "0x002a3e3a" // permitBatch(address,((address,uint160,uint48,uint48)[],address,uint256),bytes)

	// Inherited from ERC721Permit
	ERC721Permit PositionFunction = "0x0f5730f1" // permit(address,uint256,uint256,uint256,bytes)
)

// FeeAmount is an enum type representing standard Uniswap fee tiers.
type FeeAmount int

const (
	FeeLowest FeeAmount = 100   // 0.01%
	FeeLow    FeeAmount = 500   // 0.05%
	FeeMedium FeeAmount = 3000  // 0.3%
	FeeHigh   FeeAmount = 10000 // 1%
)

// TickSpacings maps fee tiers to their default tick spacing values.
var TickSpacings = map[FeeAmount]int{
	FeeLowest: 1,   // 1 tick spacing for lowest fee tier
	FeeLow:    10,  // 10 ticks for low fee tier
	FeeMedium: 60,  // 60 ticks for medium fee tier
	FeeHigh:   200, // 200 ticks for high fee tier
}
