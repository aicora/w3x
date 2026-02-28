package utils

import (
	"errors"
	"fmt"
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

func TestCheckTicks(t *testing.T) {
	tests := []struct {
		name      string
		tickLower int32
		tickUpper int32
		wantErr   error
	}{
		{
			name:      "valid range",
			tickLower: -100,
			tickUpper: 100,
			wantErr:   nil,
		},
		{
			name:      "ticks misordered",
			tickLower: 50,
			tickUpper: 50,
			wantErr:   ErrTicksMisordered,
		},
		{
			name:      "tickLower out of bounds",
			tickLower: MinTick - 1,
			tickUpper: 0,
			wantErr:   ErrTickLowerOutOfBounds,
		},
		{
			name:      "tickUpper out of bounds",
			tickLower: 0,
			tickUpper: MaxTick + 1,
			wantErr:   ErrTickUpperOutOfBounds,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckTicks(tt.tickLower, tt.tickUpper)
			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("expected nil error, got %v", err)
				}
			} else {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
			}
		})
	}
}

func TestTickSpacingToMaxLiquidityPerTick(t *testing.T) {
	tests := []struct {
		name        string
		tickSpacing int
		wantErr     error
	}{
		{"valid spacing 1", 1, nil},
		{"valid spacing 60", 60, nil},
		{"zero spacing", 0, ErrZeroTickSpacing},
		{"negative spacing", -5, ErrInvalidTickRange}, // still valid mathematically
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			maxLiquidity, err := TickSpacingToMaxLiquidityPerTick(tt.tickSpacing)
			fmt.Println("maxLiquidity", maxLiquidity)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Check that maxLiquidity is positive
			if maxLiquidity.Cmp(big.NewInt(0)) <= 0 {
				t.Fatalf("maxLiquidity should be positive, got %v", maxLiquidity)
			}
		})
	}
}