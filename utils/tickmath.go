package utils

import (
	"math/big"

	"github.com/pkg/errors"
)

const (
	// MinTick is the minimum tick that can be used on any pool.
	MinTick = -887272
	// MaxTick is the maximum tick that can be used on any pool.
	MaxTick = -MinTick
)

var (
	// MinSqrtPrice is the Q64.96 square root ratio corresponding to MinTick.
	MinSqrtPrice = big.NewInt(4295128739)

	// MaxSqrtPrice is the Q64.96 square root ratio corresponding to MaxTick.
	// Source: Uniswap v4 specification
	MaxSqrtPrice, _ = new(big.Int).SetString(
		"1461446703485210103287273052203988822378723970342", 10,
	)

	// MaxLiquidity is the maximum amount of liquidity that can be provided to a pool.
	MaxLiquidity = new(big.Int).Sub(Q128, big.NewInt(1))
)

var (
	ErrInvalidTick      = errors.New("invalid tick")
	ErrInvalidSqrtPrice = errors.New("invalid sqrt price")
	ErrTicksMisordered     = errors.New("ticks misordered")
    ErrTickLowerOutOfBounds = errors.New("tickLower out of bounds")
    ErrTickUpperOutOfBounds = errors.New("tickUpper out of bounds")
	ErrZeroTickSpacing = errors.New("tickSpacing cannot be zero")
	ErrInvalidTickRange = errors.New("invalid tick range")
)

// mulShift multiplies val by mulBy and then right-shifts the result by 128 bits.
// This is equivalent to fixed-point Q128.128 multiplication with truncation.
func mulShift(val *big.Int, mulBy *big.Int) *big.Int {
	return new(big.Int).Rsh(new(big.Int).Mul(val, mulBy), 128)
}

// Precomputed constants used in the GetSqrtPriceAtTick calculation.
// These constants mirror Uniswap v4's on-chain constants for fast sqrt price computation.
var (
	sqrtConst1, _  = new(big.Int).SetString("fffcb933bd6fad37aa2d162d1a594001", 16)
	sqrtConst2, _  = new(big.Int).SetString("100000000000000000000000000000000", 16)
	sqrtConst3, _  = new(big.Int).SetString("fff97272373d413259a46990580e213a", 16)
	sqrtConst4, _  = new(big.Int).SetString("fff2e50f5f656932ef12357cf3c7fdcc", 16)
	sqrtConst5, _  = new(big.Int).SetString("ffe5caca7e10e4e61c3624eaa0941cd0", 16)
	sqrtConst6, _  = new(big.Int).SetString("ffcb9843d60f6159c9db58835c926644", 16)
	sqrtConst7, _  = new(big.Int).SetString("ff973b41fa98c081472e6896dfb254c0", 16)
	sqrtConst8, _  = new(big.Int).SetString("ff2ea16466c96a3843ec78b326b52861", 16)
	sqrtConst9, _  = new(big.Int).SetString("fe5dee046a99a2a811c461f1969c3053", 16)
	sqrtConst10, _ = new(big.Int).SetString("fcbe86c7900a88aedcffc83b479aa3a4", 16)
	sqrtConst11, _ = new(big.Int).SetString("f987a7253ac413176f2b074cf7815e54", 16)
	sqrtConst12, _ = new(big.Int).SetString("f3392b0822b70005940c7a398e4b70f3", 16)
	sqrtConst13, _ = new(big.Int).SetString("e7159475a2c29b7443b29c7fa6e889d9", 16)
	sqrtConst14, _ = new(big.Int).SetString("d097f3bdfd2022b8845ad8f792aa5825", 16)
	sqrtConst15, _ = new(big.Int).SetString("a9f746462d870fdf8a65dc1f90e061e5", 16)
	sqrtConst16, _ = new(big.Int).SetString("70d869a156d2a1b890bb3df62baf32f7", 16)
	sqrtConst17, _ = new(big.Int).SetString("31be135f97d08fd981231505542fcfa6", 16)
	sqrtConst18, _ = new(big.Int).SetString("9aa508b5b7a84e1c677de54f3e99bc9", 16)
	sqrtConst19, _ = new(big.Int).SetString("5d6af8dedb81196699c329225ee604", 16)
	sqrtConst20, _ = new(big.Int).SetString("2216e584f5fa1ea926041bedfe98", 16)
	sqrtConst21, _ = new(big.Int).SetString("48a170391f7dc42444e8fa2", 16)
)

// GetSqrtPriceAtTick returns the square root price (Q64.96) for a given tick.
// - tick: The tick index (must be between MinTick and MaxTick).
// - Returns a *big.Int representing sqrtPriceX96.
// This mirrors the Uniswap v4 on-chain logic for gas-efficient fixed-point computation.
func GetSqrtPriceAtTick(tick int) (*big.Int, error) {
	if tick < MinTick || tick > MaxTick {
		return nil, ErrInvalidTick
	}
	absTick := tick
	if tick < 0 {
		absTick = -tick
	}
	var ratio *big.Int
	if absTick&0x1 != 0 {
		ratio = sqrtConst1
	} else {
		ratio = sqrtConst2
	}
	if (absTick & 0x2) != 0 {
		ratio = mulShift(ratio, sqrtConst3)
	}
	if (absTick & 0x4) != 0 {
		ratio = mulShift(ratio, sqrtConst4)
	}
	if (absTick & 0x8) != 0 {
		ratio = mulShift(ratio, sqrtConst5)
	}
	if (absTick & 0x10) != 0 {
		ratio = mulShift(ratio, sqrtConst6)
	}
	if (absTick & 0x20) != 0 {
		ratio = mulShift(ratio, sqrtConst7)
	}
	if (absTick & 0x40) != 0 {
		ratio = mulShift(ratio, sqrtConst8)
	}
	if (absTick & 0x80) != 0 {
		ratio = mulShift(ratio, sqrtConst9)
	}
	if (absTick & 0x100) != 0 {
		ratio = mulShift(ratio, sqrtConst10)
	}
	if (absTick & 0x200) != 0 {
		ratio = mulShift(ratio, sqrtConst11)
	}
	if (absTick & 0x400) != 0 {
		ratio = mulShift(ratio, sqrtConst12)
	}
	if (absTick & 0x800) != 0 {
		ratio = mulShift(ratio, sqrtConst13)
	}
	if (absTick & 0x1000) != 0 {
		ratio = mulShift(ratio, sqrtConst14)
	}
	if (absTick & 0x2000) != 0 {
		ratio = mulShift(ratio, sqrtConst15)
	}
	if (absTick & 0x4000) != 0 {
		ratio = mulShift(ratio, sqrtConst16)
	}
	if (absTick & 0x8000) != 0 {
		ratio = mulShift(ratio, sqrtConst17)
	}
	if (absTick & 0x10000) != 0 {
		ratio = mulShift(ratio, sqrtConst18)
	}
	if (absTick & 0x20000) != 0 {
		ratio = mulShift(ratio, sqrtConst19)
	}
	if (absTick & 0x40000) != 0 {
		ratio = mulShift(ratio, sqrtConst20)
	}
	if (absTick & 0x80000) != 0 {
		ratio = mulShift(ratio, sqrtConst21)
	}
	// If tick > 0, invert the ratio
	if tick > 0 {
		ratio = new(big.Int).Div(MaxUint256, ratio)
	}

	// Round up and convert back to Q64.96
	if new(big.Int).Rem(ratio, Q32).Cmp(big.NewInt(0)) > 0 {
		return new(big.Int).Add((new(big.Int).Div(ratio, Q32)), big.NewInt(1)), nil
	} else {
		return new(big.Int).Div(ratio, Q32), nil
	}
}

var (
	magicSqrt10001, _ = new(big.Int).SetString("255738958999603826347141", 10)
	magicTickLow, _   = new(big.Int).SetString("3402992956809132418596140100660247210", 10)
	magicTickHigh, _  = new(big.Int).SetString("291339464771989622907027621153398088495", 10)
)

// GetTickAtSqrtPrice returns the nearest tick index corresponding to a given sqrtPriceX96.
// - sqrtPriceX96: The square root price in Q64.96 format.
// - Returns the tick index and error if input is out of bounds.
// This implements Uniswap v4's logarithmic approximation for gas-efficient tick calculation.
func GetTickAtSqrtPrice(sqrtPriceX96 *big.Int) (int, error) {
	if sqrtPriceX96.Cmp(MinSqrtPrice) < 0 || sqrtPriceX96.Cmp(MaxSqrtPrice) > 0 {
		return 0, ErrInvalidSqrtPrice
	}
	// Scale to Q128.128 for computation
	sqrtPriceX128 := new(big.Int).Lsh(sqrtPriceX96, 32)
	msb, err := MostSignificantBit(sqrtPriceX128)
	if err != nil {
		return 0, err
	}
	var r *big.Int
	if big.NewInt(int64(msb)).Cmp(big.NewInt(128)) >= 0 {
		r = new(big.Int).Rsh(sqrtPriceX128, uint(msb-127))
	} else {
		r = new(big.Int).Lsh(sqrtPriceX128, uint(127-msb))
	}

	log2 := new(big.Int).Lsh(new(big.Int).Sub(big.NewInt(int64(msb)), big.NewInt(128)), 64)

	for i := 0; i < 14; i++ {
		r = new(big.Int).Rsh(new(big.Int).Mul(r, r), 127)
		f := new(big.Int).Rsh(r, 128)
		log2 = new(big.Int).Or(log2, new(big.Int).Lsh(f, uint(63-i)))
		r = new(big.Int).Rsh(r, uint(f.Int64()))
	}

	logSqrt10001 := new(big.Int).Mul(log2, magicSqrt10001)

	tickLow := new(big.Int).Rsh(new(big.Int).Sub(logSqrt10001, magicTickLow), 128).Int64()
	tickHigh := new(big.Int).Rsh(new(big.Int).Add(logSqrt10001, magicTickHigh), 128).Int64()

	if tickLow == tickHigh {
		return int(tickLow), nil
	}

	sqrtPrice, err := GetSqrtPriceAtTick(int(tickHigh))
	if err != nil {
		return 0, err
	}
	if sqrtPrice.Cmp(sqrtPriceX96) <= 0 {
		return int(tickHigh), nil
	} else {
		return int(tickLow), nil
	}
}

// CheckTicks validates that a given tick range is within allowed bounds.
//
// In Uniswap-style AMMs, ticks define the price grid boundaries.
// This function ensures the following conditions:
//   1. tickLower must be less than tickUpper, otherwise ErrTicksMisordered is returned.
//   2. tickLower must not be below the minimum system tick (MinTick), otherwise ErrTickLowerOutOfBounds is returned.
//   3. tickUpper must not exceed the maximum system tick (MaxTick), otherwise ErrTickUpperOutOfBounds is returned.
//
// Errors are wrapped with context using errors.Wrapf, preserving the stack trace
// and allowing callers to use errors.Is for exact error type checks.
//
// Parameters:
//   - tickLower: lower bound of the tick range (int32)
//   - tickUpper: upper bound of the tick range (int32)
//
// Returns:
//   - error: a wrapped error if the tick range is invalid; otherwise nil
//
// Example:
//
//    err := CheckTicks(-100, 200)
//    if err != nil {
//        if errors.Is(err, ErrTicksMisordered) {
//            fmt.Println("tickLower >= tickUpper")
//        } else if errors.Is(err, ErrTickLowerOutOfBounds) {
//            fmt.Println("tickLower out of bounds")
//        } else if errors.Is(err, ErrTickUpperOutOfBounds) {
//            fmt.Println("tickUpper out of bounds")
//        }
//    }
func CheckTicks(tickLower, tickUpper int32) error {
    if tickLower >= tickUpper {
        return errors.Wrapf(ErrTicksMisordered, "tickLower=%d >= tickUpper=%d", tickLower, tickUpper)
    }
    if tickLower < MinTick {
        return errors.Wrapf(ErrTickLowerOutOfBounds, "tickLower=%d < MinTick=%d", tickLower, MinTick)
    }
    if tickUpper > MaxTick {
        return errors.Wrapf(ErrTickUpperOutOfBounds, "tickUpper=%d > MaxTick=%d", tickUpper, MaxTick)
    }
    return nil
}

// TickSpacingToMaxLiquidityPerTick calculates the maximum liquidity allowed per tick
// given a specific tick spacing. This is used when adding liquidity to a pool
// to ensure that each tick does not exceed its maximum allowed liquidity.
//
// In Uniswap v4, the total liquidity is conceptually spread across all valid ticks.
// The function divides the maximum total liquidity (`MaxLiquidity`) evenly among
// all ticks that are multiples of `tickSpacing`.
//
// Parameters:
//   - tickSpacing: the required separation between initialized ticks. For example,
//     a tickSpacing of 3 allows ticks to be initialized at ..., -6, -3, 0, 3, 6, ...
//
// Returns:
//   - *big.Int: the maximum liquidity that can be assigned to a single tick.
//   - error: if tickSpacing is zero or the tick range is invalid.
//
// Errors:
//   - ErrZeroTickSpacing: returned when tickSpacing is zero.
//   - ErrInvalidTickRange: returned if computed number of ticks is non-positive.
//
// Example usage:
//
//     maxLiquidityPerTick, err := TickSpacingToMaxLiquidityPerTick(60)
//     if err != nil {
//         log.Fatal(err)
//     }
//     fmt.Println("Max liquidity per tick:", maxLiquidityPerTick)
func TickSpacingToMaxLiquidityPerTick(tickSpacing int) (*big.Int, error) {
	if tickSpacing == 0 {
		return nil, ErrZeroTickSpacing
	}

	minTick := MinTick / tickSpacing
	if MinTick%tickSpacing != 0 {
		minTick -= 1
	}

	maxTick := MaxTick / tickSpacing

	numTicks := maxTick - minTick + 1
	if numTicks <= 0 {
		return nil, ErrInvalidTickRange
	}

	return new(big.Int).Div(MaxLiquidity, big.NewInt(int64(numTicks))), nil
}