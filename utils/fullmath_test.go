package utils

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMulDiv(t *testing.T) {
	tests := []struct {
		a, b, den int64
		want      int64
		wantErr   bool
	}{
		{a: 10, b: 20, den: 5, want: 40},
		{a: 1, b: 1, den: 2, want: 0}, // floor(0.5)
		{a: 2, b: 3, den: 0, want: 0, wantErr: true}, // zero denominator
		{a: 0, b: 100, den: 5, want: 0}, // zero multiplicand
	}

	for _, tt := range tests {
		a := big.NewInt(tt.a)
		b := big.NewInt(tt.b)
		den := big.NewInt(tt.den)

		got, err := MulDiv(a, b, den)
		if tt.wantErr {
			require.Error(t, err)
			continue
		}
		require.NoError(t, err)
		gotint64 := got.Int64()
		require.Equal(t, tt.want, gotint64)
	}
}

func TestMulDivRoundingUp(t *testing.T) {
	tests := []struct {
		a, b, den int64
		want      int64
		wantErr   bool
	}{
		{a: 10, b: 20, den: 5, want: 40},
		{a: 1, b: 1, den: 2, want: 1}, // ceil(0.5)
		{a: 2, b: 3, den: 0, want: 0, wantErr: true},
		{a: 0, b: 100, den: 5, want: 0},
	}

	for _, tt := range tests {
		a := big.NewInt(tt.a)
		b := big.NewInt(tt.b)
		den := big.NewInt(tt.den)

		got, err := MulDivRoundingUp(a, b, den)
		if tt.wantErr {
			require.Error(t, err)
			continue
		}
		require.NoError(t, err)
		gotint64 := got.Int64()
		require.Equal(t, tt.want, gotint64)
	}
}

func TestMulMod(t *testing.T) {
	tests := []struct {
		a, b, m int64
		want    int64
	}{
		{a: 10, b: 20, m: 6, want: 2},
		{a: 123, b: 456, m: 1000, want: 88},
		{a: 0, b: 100, m: 7, want: 0},
	}

	for _, tt := range tests {
		a := big.NewInt(tt.a)
		b := big.NewInt(tt.b)
		m := big.NewInt(tt.m)

		got := MulMod(a, b, m)
		gotint64 := got.Int64()
		require.Equal(t, tt.want, gotint64)
	}
}

func TestAbsBigInt(t *testing.T) {
	tests := []struct {
		val  int64
		want int64
	}{
		{val: 10, want: 10},
		{val: -10, want: 10},
		{val: 0, want: 0},
	}

	for _, tt := range tests {
		x := big.NewInt(tt.val)
		got := AbsBigInt(x)
		require.Equal(t, tt.want, got.Int64())
	}
}