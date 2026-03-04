package sdkv3

import (
	"math/big"

	"github.com/pkg/errors"
)

var (
	// ErrZeroDenominator is returned when a division operation is attempted with a zero denominator.
	ErrZeroDenominator = errors.New("denominator is zero")
)

// EncodeSqrtRatioX96 calculates the square root of a token1/token0 ratio
// and returns it in Q64.96 fixed-point format using big.Int arithmetic.
//
// This is commonly used in AMM pools (e.g., Uniswap v4) to encode the
// price ratio in a compact high-precision format suitable for tick calculations.
//
// Formula:
//
//	sqrtPriceX96 = sqrt(amount1 / amount0) * 2^96
//
// Parameters:
//   - amount1: numerator amount (token1, in smallest unit, e.g., wei)
//   - amount0: denominator amount (token0, in smallest unit)
//
// Returns:
//   - *big.Int: sqrt(amount1/amount0) scaled by 2^96 (Q64.96 format)
//   - error: if amount0 is zero
func EncodeSqrtRatioX96(amount1, amount0 *big.Int) *big.Int {
	// Step 1: numerator * 2^192 to prepare for Q64.96 sqrt scaling
	num := new(big.Int).Lsh(amount1, 192) // amount1 * 2^192

	// Step 2: division by denominator
	ratioX192 := new(big.Int).Div(num, amount0) // ratio * 2^192

	// Step 3: square root
	sqrtRatioX96 := new(big.Int).Sqrt(ratioX192) // sqrt(ratio * 2^192) = sqrt(ratio) * 2^96

	return sqrtRatioX96
}
