package sdkv3

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSqrtPriceX96(t *testing.T) {
	tests := []struct {
		name    string
		amount1 *big.Int
		amount0 *big.Int
		wantGT  int64
	}{
		{
			name:    "normal values",
			amount1: big.NewInt(2000_000_000),              // e.g., 2000 USDC (6 decimals)
			amount0: big.NewInt(1_000_000_000_000_000_000), // 1 ETH in wei
			wantGT:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := EncodeSqrtRatioX96(tt.amount1, tt.amount0)
			require.True(t, res.Cmp(big.NewInt(tt.wantGT)) > 0)
		})
	}
}
