package utils

import (
	"math/big"
)

// MaxFee defines the maximum swap fee in hundredths of a bip (1e6 = 100%).
var (
	MaxFee = new(big.Int).Exp(big.NewInt(10), big.NewInt(6), nil)
)

// ComputeSwapStep calculates the result of a single swap step within a tick.
//
// This function is a Go translation of the Uniswap V3 swap step calculation,
// handling both exact input and exact output swaps, taking liquidity, sqrt prices,
// and swap fees into account.
//
// Parameters:
//   - sqrtPriceCurrentX96: current sqrt price in Q64.96 format
//   - sqrtPriceTargetX96: target sqrt price for this swap step (usually the next initialized tick)
//   - liquidity: current in-range liquidity
//   - amountRemaining: amount left to swap (positive for exact input, negative for exact output)
//   - feePips: swap fee in hundredths of a bip (1e6 = 100%)
//
// Returns:
//   - sqrtPriceNextX96: the sqrt price after this swap step
//   - amountIn: amount of input token consumed
//   - amountOut: amount of output token produced
//   - feeAmount: swap fee taken in input token units
//   - err: error if any computation fails
//
// Behavior:
//   - Determines the swap direction (zeroForOne) based on current and target sqrt prices.
//   - For exact input swaps:
//       * Deducts the fee from the remaining input amount
//       * Checks if target price can be reached with remaining input
//       * If not, computes next sqrt price using GetNextSqrtPriceFromInput
//   - For exact output swaps:
//       * Checks if target price can be reached with remaining output
//       * If not, computes next sqrt price using GetNextSqrtPriceFromOutput
//   - Computes the input/output amounts based on liquidity and swap direction
//   - Calculates the fee amount depending on whether target price was fully reached
//
// Notes:
//   - Uses big.Int arithmetic to handle Q64.96 fixed-point numbers and large liquidity values
//   - Implements rounding up for fee calculations via MulDivRoundingUp
//   - Mirrors Uniswap V3 core swap step logic
func ComputeSwapStep(sqrtPriceCurrentX96, sqrtPriceTargetX96, liquidity, amountRemaining *big.Int, feePips uint32) (sqrtPriceNextX96, amountIn, amountOut, feeAmount *big.Int, err error) {
	zeroForOne := sqrtPriceCurrentX96.Cmp(sqrtPriceTargetX96) >= 0
	exactIn := amountRemaining.Cmp(big.NewInt(0)) >= 0

	if exactIn {
		// deduct fee from remaining input
		amountRemainingLessFee := new(big.Int).Div(new(big.Int).Mul(amountRemaining, new(big.Int).Sub(MaxFee, big.NewInt(int64(feePips)))), MaxFee)
		if zeroForOne {
			amountIn, err = GetAmount0Delta(sqrtPriceTargetX96, sqrtPriceCurrentX96, liquidity, true)
			if err != nil {
				return
			}
		} else {
			amountIn, err = GetAmount1Delta(sqrtPriceCurrentX96, sqrtPriceTargetX96, liquidity, true)
			if err != nil {
				return
			}
		}
		if amountRemainingLessFee.Cmp(amountIn) >= 0 {
			sqrtPriceNextX96 = sqrtPriceTargetX96
		} else {
			sqrtPriceNextX96, err = GetNextSqrtPriceFromInput(sqrtPriceCurrentX96, liquidity, amountRemainingLessFee, zeroForOne)
			if err != nil {
				return
			}
		}
	} else {
		// exact output swap
		if zeroForOne {
			amountOut, err = GetAmount1Delta(sqrtPriceTargetX96, sqrtPriceCurrentX96, liquidity, false)
			if err != nil {
				return
			}
		} else {
			amountOut, err = GetAmount0Delta(sqrtPriceCurrentX96, sqrtPriceTargetX96, liquidity, false)
			if err != nil {
				return
			}
		}
		if new(big.Int).Mul(amountRemaining, big.NewInt(-1)).Cmp(amountOut) >= 0 {
			sqrtPriceNextX96 = sqrtPriceTargetX96
		} else {
			sqrtPriceNextX96, err = GetNextSqrtPriceFromOutput(sqrtPriceCurrentX96, liquidity, new(big.Int).Mul(amountRemaining, big.NewInt(-1)), zeroForOne)
			if err != nil {
				return
			}
		}
	}

	max := sqrtPriceTargetX96.Cmp(sqrtPriceNextX96) == 0

	// compute actual input/output deltas based on swap direction and whether max reached
	if zeroForOne {
		if !(max && exactIn) {
			amountIn, err = GetAmount0Delta(sqrtPriceNextX96, sqrtPriceCurrentX96, liquidity, true)
			if err != nil {
				return
			}
		}
		if !(max && !exactIn) {
			amountOut, err = GetAmount1Delta(sqrtPriceNextX96, sqrtPriceCurrentX96, liquidity, false)
			if err != nil {
				return
			}
		}
	} else {
		if !(max && exactIn) {
			amountIn, err = GetAmount1Delta(sqrtPriceCurrentX96, sqrtPriceNextX96, liquidity, true)
			if err != nil {
				return
			}
		}
		if !(max && !exactIn) {
			amountOut, err = GetAmount0Delta(sqrtPriceCurrentX96, sqrtPriceNextX96, liquidity, false)
			if err != nil {
				return
			}
		}
	}

	// clamp output for exact output swaps
	if !exactIn && amountOut.Cmp(new(big.Int).Mul(amountRemaining, big.NewInt(-1))) > 0 {
		amountOut = new(big.Int).Mul(amountRemaining, big.NewInt(-1))
	}

	// compute fee amount
	if exactIn && sqrtPriceNextX96.Cmp(sqrtPriceTargetX96) != 0 {
		// target not reached, remainder is fee
		feeAmount = new(big.Int).Sub(amountRemaining, amountIn)
	} else {
		feeAmount, err = MulDivRoundingUp(amountIn, big.NewInt(int64(feePips)), new(big.Int).Sub(MaxFee, big.NewInt(int64(feePips))))
		if err != nil {
			return
		}
	}

	return
}


// GetSqrtPriceTarget computes the next sqrt price for a swap step.
//
// zeroForOne: true if swapping 0→1, false if 1→0
// sqrtPriceNextX96: next initialized tick price (Q64.96)
// sqrtPriceLimitX96: price limit (Q64.96)
func GetSqrtPriceTarget(zeroForOne bool, sqrtPriceNextX96, sqrtPriceLimitX96 *big.Int) *big.Int {
	target := new(big.Int)
	if zeroForOne {
		// 0→1 swap, price cannot go below limit
		if sqrtPriceNextX96.Cmp(sqrtPriceLimitX96) < 0 {
			target.Set(sqrtPriceLimitX96)
		} else {
			target.Set(sqrtPriceNextX96)
		}
	} else {
		// 1→0 swap, price cannot go above limit
		if sqrtPriceNextX96.Cmp(sqrtPriceLimitX96) > 0 {
			target.Set(sqrtPriceLimitX96)
		} else {
			target.Set(sqrtPriceNextX96)
		}
	}
	return target
}