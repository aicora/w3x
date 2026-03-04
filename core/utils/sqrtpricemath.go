package utils

import (
	"math/big"

	"github.com/pkg/errors"
)

var (
	ErrInvariant             = errors.New("price invariant violation")
	ErrZeroLiquidity         = errors.New("liquidity is zero")
	ErrSqrtPriceLessThanZero = errors.New("sqrt price is less than zero")
	ErrLiquidityLessThanZero = errors.New("liquidity is less than zero")
)

// GetAmount0Delta calculates the amount of token0 required for a given
// liquidity between two sqrt price ratios, using big.Int arithmetic.
//
// This is equivalent to the Uniswap v4 formula:
//
//	amount0 = L * (sqrtPriceB - sqrtPriceA) / (sqrtPriceB * sqrtPriceA)
//
// Parameters:
//   - sqrtPriceAX96: lower sqrt price as Q64.96
//   - sqrtPriceBX96: upper sqrt price as Q64.96
//   - liquidity: liquidity amount
//   - roundUp: whether to round the result up
//
// Returns:
//   - *big.Int: token0 amount delta
//   - error: if any division by zero occurs
func GetAmount0Delta(sqrtPriceAX96, sqrtPriceBX96, liquidity *big.Int, roundUp bool) (*big.Int, error) {
	// Ensure sqrtPriceAX96 <= sqrtPriceBX96
	if sqrtPriceAX96.Cmp(sqrtPriceBX96) > 0 {
		sqrtPriceAX96, sqrtPriceBX96 = sqrtPriceBX96, sqrtPriceAX96
	}

	// numerator1 = liquidity << 96
	numerator1 := new(big.Int).Lsh(liquidity, 96)

	// numerator2 = sqrtPriceBX96 - sqrtPriceAX96
	numerator2 := new(big.Int).Sub(sqrtPriceBX96, sqrtPriceAX96)

	if roundUp {
		// ceil( (liquidity<<96 * (sqrtB - sqrtA)) / (sqrtB) )
		temp, err := MulDivRoundingUp(numerator1, numerator2, sqrtPriceBX96)
		if err != nil {
			return nil, err
		}

		// ceil( temp / sqrtA )
		result, err := MulDivRoundingUp(temp, big.NewInt(1), sqrtPriceAX96)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	// Without rounding: floor((liquidity<<96 * (sqrtB - sqrtA)) / (sqrtB * sqrtA))
	temp, err := MulDiv(numerator1, numerator2, sqrtPriceBX96)
	if err != nil {
		return nil, err
	}
	result, err := MulDiv(temp, big.NewInt(1), sqrtPriceAX96)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetAmount1Delta calculates the amount of token1 required for a given
// liquidity between two sqrt price ratios, using big.Int arithmetic.
//
// This is equivalent to the Uniswap v4 formula:
//
//	amount1 = L * (sqrtPriceB - sqrtPriceA) / 2^96
//
// Parameters:
//   - sqrtPriceAX96: lower sqrt price as Q64.96
//   - sqrtPriceBX96: upper sqrt price as Q64.96
//   - liquidity: liquidity amount
//   - roundUp: whether to round the result up
//
// Returns:
//   - *big.Int: token1 amount delta
//   - error: if any division by zero occurs
func GetAmount1Delta(sqrtPriceAX96, sqrtPriceBX96, liquidity *big.Int, roundUp bool) (*big.Int, error) {
	// Ensure sqrtPriceAX96 <= sqrtPriceBX96
	if sqrtPriceAX96.Cmp(sqrtPriceBX96) > 0 {
		sqrtPriceAX96, sqrtPriceBX96 = sqrtPriceBX96, sqrtPriceAX96
	}

	// delta = sqrtPriceBX96 - sqrtPriceAX96
	delta := new(big.Int).Sub(sqrtPriceBX96, sqrtPriceAX96)

	if roundUp {
		result, err := MulDivRoundingUp(liquidity, delta, Q96)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	// Without rounding: floor(liquidity * delta / Q96)
	result, err := MulDiv(liquidity, delta, Q96)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// getNextSqrtPriceFromAmount1RoundingDown computes the next sqrt price
// (Q64.96) given a change in token1 amount.
//
// This is equivalent to the Uniswap v4 formula:
//
//	add: sqrtPriceNext = sqrtPriceX96 + (amount1 / liquidity)
//	remove: sqrtPriceNext = sqrtPriceX96 - (amount1 / liquidity)
//
// Parameters:
//   - sqrtPriceX96: current sqrt price in Q64.96
//   - liquidity: current liquidity L
//   - amount: amount of token1 added/removed
//   - add: true if token1 is added, false if removed
//
// Returns:
//   - *big.Int: next sqrt price
//   - error: if liquidity = 0 or invariant would be violated
func GetNextSqrtPriceFromAmount1RoundingDown(sqrtPriceX96, liquidity, amount *big.Int, add bool) (*big.Int, error) {
	if liquidity.Sign() == 0 {
		return nil, ErrZeroLiquidity
	}

	if add {
		var quotient *big.Int
		// If amount is small, we can use left shift
		if amount.Cmp(MaxUint160) <= 0 {
			quotient = new(big.Int).Div(new(big.Int).Lsh(amount, 96), liquidity)
		} else {
			quotient = new(big.Int).Div(new(big.Int).Mul(amount, Q96), liquidity)
		}
		return new(big.Int).Add(sqrtPriceX96, quotient), nil
	}

	// Subtracting token1: price decreases
	quotient, err := MulDivRoundingUp(amount, Q96, liquidity)
	if err != nil {
		return nil, err
	}
	if sqrtPriceX96.Cmp(quotient) <= 0 {
		return nil, ErrInvariant
	}
	return new(big.Int).Sub(sqrtPriceX96, quotient), nil
}

// GetNextSqrtPriceFromAmount0RoundingUp calculates the next sqrt price
// after adding or removing a given amount of token0, using big.Int arithmetic.
//
// This implements the Uniswap v4 formula for token0:
//
//	add: sqrtPriceNext = (L * sqrtPriceX96) / (L + (amount0 * sqrtPriceX96))
//	remove: sqrtPriceNext = (L * sqrtPriceX96) / (L - (amount0 * sqrtPriceX96))
//
// - Adding token0 (add=true) decreases the sqrt price.
// - Removing token0 (add=false) increases the sqrt price.
// - Rounding up is applied to ensure the AMM invariant is preserved.
//
// Parameters:
//   - sqrtPX96: current sqrt price in Q64.96
//   - liquidity: liquidity in the tick range
//   - amount: token0 amount to add or remove
//   - add: true if token0 is added, false if removed
//
// Returns:
//   - *big.Int: next sqrt price in Q64.96
//   - error: ErrInvariant if the operation violates the AMM invariant
func GetNextSqrtPriceFromAmount0RoundingUp(sqrtPriceX96, liquidity, amount *big.Int, add bool) (*big.Int, error) {
	if amount.Sign() == 0 {
		return sqrtPriceX96, nil
	}

	// numerator1 = L * 2^96 (Q64.96 scaling)
	numerator1 := new(big.Int).Lsh(liquidity, 96)

	if add {
		product := new(big.Int).Mul(amount, sqrtPriceX96)
		product.And(product, MaxUint256)

		// Overflow check: product / amount == sqrtPriceX96
		if new(big.Int).Div(product, amount).Cmp(sqrtPriceX96) == 0 {
			denominator := new(big.Int).Add(numerator1, product)
			denominator.And(denominator, MaxUint256)

			if denominator.Cmp(numerator1) >= 0 {
				return MulDivRoundingUp(numerator1, sqrtPriceX96, denominator)
			}
		}

		// Fallback for large amounts
		return MulDivRoundingUp(numerator1, big.NewInt(1),
			new(big.Int).Add(
				new(big.Int).Div(numerator1, sqrtPriceX96),
				amount,
			),
		)
	} else {
		// Removing token0
		product := new(big.Int).Mul(amount, sqrtPriceX96)
		product.And(product, MaxUint256)

		// Overflow or insufficient liquidity check
		if new(big.Int).Div(product, amount).Cmp(sqrtPriceX96) != 0 {
			return nil, ErrInvariant
		}
		if numerator1.Cmp(product) <= 0 {
			return nil, ErrInvariant
		}

		denominator := new(big.Int).Sub(numerator1, product)
		return MulDivRoundingUp(numerator1, sqrtPriceX96, denominator)
	}
}

// GetNextSqrtPriceFromInput calculates the next sqrt price after a given
// input amount is swapped into the pool.
//
// This function handles both token0 → token1 swaps and token1 → token0 swaps
// depending on the `zeroForOne` flag.
//
// Parameters:
//   - sqrtPriceX96: current sqrt price in Q64.96 format (√(P) = √(token1/token0))
//   - liquidity: current active liquidity in the tick range
//   - amountIn: amount of token being sent into the pool
//   - zeroForOne: true if swapping token0 → token1, false for token1 → token0
//
// Returns:
//   - *big.Int: next sqrt price in Q64.96 after accounting for the input amount
//   - error: if sqrt price or liquidity is zero/negative, or underlying calculation fails
func GetNextSqrtPriceFromInput(sqrtPriceX96, liquidity, amountIn *big.Int, zeroForOne bool) (*big.Int, error) {
	if sqrtPriceX96.Cmp(big.NewInt(0)) <= 0 {
		return nil, ErrSqrtPriceLessThanZero
	}
	if liquidity.Cmp(big.NewInt(0)) <= 0 {
		return nil, ErrLiquidityLessThanZero
	}
	if zeroForOne {
		// Input is token0 → token1: increasing token0 decreases sqrt price
		return GetNextSqrtPriceFromAmount0RoundingUp(sqrtPriceX96, liquidity, amountIn, true)
	}
	// Input is token1 → token0: increasing token1 increases sqrt price
	return GetNextSqrtPriceFromAmount1RoundingDown(sqrtPriceX96, liquidity, amountIn, true)
}

// GetNextSqrtPriceFromOutput calculates the next sqrt price after a given
// output amount is swapped out from the pool.
//
// This function handles both token0 → token1 swaps and token1 → token0 swaps
// depending on the `zeroForOne` flag.
//
// Parameters:
//   - sqrtPriceX96: current sqrt price in Q64.96 format (√(P) = √(token1/token0))
//   - liquidity: current active liquidity in the tick range
//   - amountOut: amount of token being withdrawn from the pool
//   - zeroForOne: true if swapping token0 → token1, false for token1 → token0
//
// Returns:
//   - *big.Int: next sqrt price in Q64.96 after accounting for the output amount
//   - error: if sqrt price or liquidity is zero/negative, or underlying calculation fails
func GetNextSqrtPriceFromOutput(sqrtPriceX96, liquidity, amountOut *big.Int, zeroForOne bool) (*big.Int, error) {
	if sqrtPriceX96.Cmp(big.NewInt(0)) <= 0 {
		return nil, ErrSqrtPriceLessThanZero
	}
	if liquidity.Cmp(big.NewInt(0)) <= 0 {
		return nil, ErrLiquidityLessThanZero
	}
	if zeroForOne {
		// Output is token0 → token1: removing token1 increases sqrt price
		return GetNextSqrtPriceFromAmount1RoundingDown(sqrtPriceX96, liquidity, amountOut, false)
	}
	// Output is token1 → token0: removing token0 decreases sqrt price
	return GetNextSqrtPriceFromAmount0RoundingUp(sqrtPriceX96, liquidity, amountOut, false)
}
