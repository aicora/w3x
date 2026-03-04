package utils

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetAmount0Delta(t *testing.T) {
	liq := big.NewInt(1000)
	sqrtA := new(big.Int).Set(Q96)
	sqrtB := new(big.Int).Mul(Q96, big.NewInt(4))

	tests := []struct {
		name    string
		roundUp bool
	}{
		{"round down", false},
		{"round up", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := GetAmount0Delta(sqrtA, sqrtB, liq, tt.roundUp)
			require.NoError(t, err)
			require.True(t, res.Cmp(big.NewInt(0)) > 0)
		})
	}
}

func TestGetAmount1Delta(t *testing.T) {
	liq := big.NewInt(1000)
	sqrtA := new(big.Int).Set(Q96)
	sqrtB := new(big.Int).Mul(Q96, big.NewInt(4))

	tests := []struct {
		name    string
		roundUp bool
	}{
		{"round down", false},
		{"round up", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := GetAmount1Delta(sqrtA, sqrtB, liq, tt.roundUp)
			require.NoError(t, err)
			require.True(t, res.Cmp(big.NewInt(0)) > 0)
		})
	}
}

func TestGetNextSqrtPriceFromAmount1RoundingDown(t *testing.T) {
	liq := new(big.Int).SetInt64(1000)
	sqrtPX96 := new(big.Int).Set(Q96) // current sqrt price = 2^96
	amountSmall := big.NewInt(500)
	amountLarge := new(big.Int).Lsh(big.NewInt(1), 160) // 超过 MaxUint160

	tests := []struct {
		name      string
		sqrtPX96  *big.Int
		liquidity *big.Int
		amount    *big.Int
		add       bool
		expectErr bool
	}{
		{"add small amount", sqrtPX96, liq, amountSmall, true, false},
		{"add large amount", sqrtPX96, liq, amountLarge, true, false},
		{"subtract small amount", sqrtPX96, liq, amountSmall, false, false},
		{"subtract too much amount", sqrtPX96, liq, new(big.Int).Mul(sqrtPX96, liq), false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextPrice, err := GetNextSqrtPriceFromAmount1RoundingDown(tt.sqrtPX96, tt.liquidity, tt.amount, tt.add)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, nextPrice)
				// 对 add/subtract 做基本验证
				if tt.add {
					require.True(t, nextPrice.Cmp(tt.sqrtPX96) > 0)
				} else if !tt.expectErr {
					require.True(t, nextPrice.Cmp(tt.sqrtPX96) < 0)
				}
			}
		})
	}
}

func TestGetNextSqrtPriceFromAmount0RoundingUp(t *testing.T) {
	liquidity := big.NewInt(1000)
	sqrtPriceX96 := new(big.Int).Set(Q96)

	tests := []struct {
		name      string
		amount0   *big.Int
		add       bool
		wantError bool
	}{
		{"add small amount0", big.NewInt(10), true, false},
		{"remove small amount0", big.NewInt(10), false, false},
		{"add zero amount0", big.NewInt(0), true, false},
		{"remove zero amount0", big.NewInt(0), false, false},
		{"remove too much amount0", big.NewInt(1e18), false, true}, // should fail
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextPrice, err := GetNextSqrtPriceFromAmount0RoundingUp(sqrtPriceX96, liquidity, tt.amount0, tt.add)
			if (err != nil) != tt.wantError {
				t.Fatalf("unexpected error: %v", err)
			}
			if err != nil {
				return
			}

			// 处理 amount0 = 0 的特殊情况
			if tt.amount0.Sign() == 0 {
				if nextPrice.Cmp(sqrtPriceX96) != 0 {
					t.Errorf("amount0=0: expected nextPrice == currentPrice")
				}
				return
			}

			// 检查价格变化方向
			if tt.add && nextPrice.Cmp(sqrtPriceX96) >= 0 {
				t.Errorf("add: expected nextPrice < currentPrice")
			}
			if !tt.add && nextPrice.Cmp(sqrtPriceX96) <= 0 {
				t.Errorf("remove: expected nextPrice > currentPrice")
			}
		})
	}
}

func TestGetNextSqrtPriceFromInput(t *testing.T) {
	liquidity := big.NewInt(1000)
	sqrtPriceX96 := new(big.Int).Set(Q96)

	tests := []struct {
		name       string
		amountIn   *big.Int
		zeroForOne bool
		wantErr    bool
	}{
		{"token0->token1 small", big.NewInt(10), true, false},
		{"token1->token0 small", big.NewInt(10), false, false},
		{"zero amount", big.NewInt(0), true, false},
		{"zero amount", big.NewInt(0), false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextPrice, err := GetNextSqrtPriceFromInput(sqrtPriceX96, liquidity, tt.amountIn, tt.zeroForOne)
			if (err != nil) != tt.wantErr {
				t.Fatalf("unexpected error: %v", err)
			}
			if err != nil {
				return
			}

			if tt.amountIn.Sign() == 0 && nextPrice.Cmp(sqrtPriceX96) != 0 {
				t.Errorf("amount=0: expected nextPrice == currentPrice")
			}

			if tt.zeroForOne && tt.amountIn.Sign() > 0 && nextPrice.Cmp(sqrtPriceX96) >= 0 {
				t.Errorf("token0->token1: expected nextPrice < currentPrice")
			}
			if !tt.zeroForOne && tt.amountIn.Sign() > 0 && nextPrice.Cmp(sqrtPriceX96) <= 0 {
				t.Errorf("token1->token0: expected nextPrice > currentPrice")
			}
		})
	}
}

func TestGetNextSqrtPriceFromOutput(t *testing.T) {
	liquidity := big.NewInt(1000)
	sqrtPriceX96 := new(big.Int).Set(Q96)

	tests := []struct {
		name       string
		amountOut  *big.Int
		zeroForOne bool
		wantErr    bool
	}{
		{"token0->token1 small", big.NewInt(10), true, false},
		{"token1->token0 small", big.NewInt(10), false, false},
		{"zero amount", big.NewInt(0), true, false},
		{"zero amount", big.NewInt(0), false, false},
		{"too much amount", big.NewInt(1e18), true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextPrice, err := GetNextSqrtPriceFromOutput(sqrtPriceX96, liquidity, tt.amountOut, tt.zeroForOne)
			if (err != nil) != tt.wantErr {
				t.Fatalf("unexpected error: %v", err)
			}
			if err != nil {
				return
			}

			if tt.amountOut.Sign() == 0 && nextPrice.Cmp(sqrtPriceX96) != 0 {
				t.Errorf("amount=0: expected nextPrice == currentPrice")
			}
			if tt.zeroForOne && tt.amountOut.Sign() > 0 && nextPrice.Cmp(sqrtPriceX96) >= 0 {
				t.Errorf("token0->token1 remove: expected nextPrice < currentPrice")
			}
			if !tt.zeroForOne && tt.amountOut.Sign() > 0 && nextPrice.Cmp(sqrtPriceX96) <= 0 {
				t.Errorf("token1->token0 remove: expected nextPrice > currentPrice")
			}
		})
	}
}
