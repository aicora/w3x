package utils

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetSqrtPriceTarget(t *testing.T) {
	sqrtPriceNext := big.NewInt(100)
	sqrtPriceLimit := big.NewInt(90)

	// zeroForOne swap
	target := GetSqrtPriceTarget(true, sqrtPriceNext, sqrtPriceLimit)
	require.Equal(t, sqrtPriceNext, target)

	sqrtPriceNext = big.NewInt(80)
	target = GetSqrtPriceTarget(true, sqrtPriceNext, sqrtPriceLimit)
	require.Equal(t, sqrtPriceLimit, target)

	// oneForZero swap
	sqrtPriceNext = big.NewInt(100)
	target = GetSqrtPriceTarget(false, sqrtPriceNext, sqrtPriceLimit)
	require.Equal(t, sqrtPriceLimit, target)

	sqrtPriceNext = big.NewInt(95)
	target = GetSqrtPriceTarget(false, sqrtPriceNext, sqrtPriceLimit)
	require.Equal(t, sqrtPriceLimit, target)
}

func TestComputeSwapStepExactInZeroForOne(t *testing.T) {
	liquidity := big.NewInt(1_000_000)
	sqrtPriceCurrent := big.NewInt(100)
	sqrtPriceTarget := big.NewInt(120)
	amountRemaining := big.NewInt(10_000)
	feePips := uint32(3000)

	sqrtNext, amountIn, amountOut, feeAmount, err := ComputeSwapStep(sqrtPriceCurrent, sqrtPriceTarget, liquidity, amountRemaining, feePips)
	require.NoError(t, err)
	require.NotNil(t, sqrtNext)
	require.NotNil(t, amountIn)
	require.NotNil(t, amountOut)
	require.NotNil(t, feeAmount)
}

func TestComputeSwapStepExactOutOneForZero(t *testing.T) {
	liquidity := big.NewInt(1_000_000)
	sqrtPriceCurrent := big.NewInt(120)
	sqrtPriceTarget := big.NewInt(100)
	amountRemaining := big.NewInt(-5_000) // exact output
	feePips := uint32(3000)

	sqrtNext, amountIn, amountOut, feeAmount, err := ComputeSwapStep(sqrtPriceCurrent, sqrtPriceTarget, liquidity, amountRemaining, feePips)
	require.NoError(t, err)
	require.NotNil(t, sqrtNext)
	require.NotNil(t, amountIn)
	require.NotNil(t, amountOut)
	require.NotNil(t, feeAmount)
}

func TestComputeSwapStepReachesTarget(t *testing.T) {
	liquidity := big.NewInt(1_000_000)
	sqrtPriceCurrent := big.NewInt(100)
	sqrtPriceTarget := big.NewInt(100)
	amountRemaining := big.NewInt(10_000)
	feePips := uint32(3000)

	sqrtNext, amountIn, amountOut, feeAmount, err := ComputeSwapStep(sqrtPriceCurrent, sqrtPriceTarget, liquidity, amountRemaining, feePips)
	require.NoError(t, err)
	require.Equal(t, sqrtPriceTarget, sqrtNext)
	require.True(t, amountIn.Cmp(big.NewInt(0)) >= 0)
	require.True(t, amountOut.Cmp(big.NewInt(0)) >= 0)
	require.True(t, feeAmount.Cmp(big.NewInt(0)) >= 0)
}