package utils

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetSqrtPriceAtTick(t *testing.T) {
	sqrtPrice, err := GetSqrtPriceAtTick(MinTick)
	require.NoError(t, err)
	require.True(t, sqrtPrice.Cmp(MinSqrtPrice) >= 0)

	sqrtPrice, err = GetSqrtPriceAtTick(MaxTick)
	require.NoError(t, err)
	require.True(t, sqrtPrice.Cmp(MaxSqrtPrice) <= 0)

	sqrtPrice, err = GetSqrtPriceAtTick(0)
	require.NoError(t, err)
	require.NotNil(t, sqrtPrice)

	_, err = GetSqrtPriceAtTick(MinTick - 1)
	require.ErrorIs(t, err, ErrInvalidTick)

	_, err = GetSqrtPriceAtTick(MaxTick + 1)
	require.ErrorIs(t, err, ErrInvalidTick)
}

func TestGetTickAtSqrtPrice(t *testing.T) {
	tick, err := GetTickAtSqrtPrice(MinSqrtPrice)
	require.NoError(t, err)
	require.Equal(t, MinTick, tick)

	maxSqrtMinus1 := new(big.Int).Sub(MaxSqrtPrice, big.NewInt(1))
	tick, err = GetTickAtSqrtPrice(maxSqrtMinus1)
	require.NoError(t, err)
	require.True(t, tick <= MaxTick)

	midTick := 0
	sqrtPrice, err := GetSqrtPriceAtTick(midTick)
	require.NoError(t, err)
	tick, err = GetTickAtSqrtPrice(sqrtPrice)
	require.NoError(t, err)
	require.Equal(t, midTick, tick)

	tooLow := new(big.Int).Sub(MinSqrtPrice, big.NewInt(1))
	_, err = GetTickAtSqrtPrice(tooLow)
	require.ErrorIs(t, err, ErrInvalidSqrtPrice)

	tooHigh := new(big.Int).Add(MaxSqrtPrice, big.NewInt(1))
	_, err = GetTickAtSqrtPrice(tooHigh)
	require.ErrorIs(t, err, ErrInvalidSqrtPrice)
}

func TestSqrtPriceTickRoundTrip(t *testing.T) {
	testTicks := []int{MinTick, -500_000, -1, 0, 1, 500_000, MaxTick}
	for _, tick := range testTicks {
		sqrtPrice, err := GetSqrtPriceAtTick(tick)
		require.NoError(t, err)
		tick2, err := GetTickAtSqrtPrice(sqrtPrice)
		require.NoError(t, err)
		require.True(t, tick2 == tick || tick2 == tick-1 || tick2 == tick+1, "tick roundtrip mismatch")
	}
}