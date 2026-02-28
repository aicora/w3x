package libraries

import (
	"errors"
	"math/big"
)

var (
	// ErrInt128Overflow is returned when a value exceeds the signed int128 range.
	//
	// Solidity's BalanceDelta packs two int128 values into a single int256.
	// To preserve exact EVM semantics, this implementation enforces the same
	// numeric bounds in Go.
	ErrInt128Overflow = errors.New("int128 overflow")
)

var (
	// int128Max represents the maximum value of a signed int128:
	//  2^127 - 1
	int128Max = new(big.Int).Sub(
		new(big.Int).Lsh(big.NewInt(1), 127),
		big.NewInt(1),
	)

	// int128Min represents the minimum value of a signed int128:
	// -2^127
	int128Min = new(big.Int).Neg(
		new(big.Int).Lsh(big.NewInt(1), 127),
	)
)

// BalanceDelta represents the net change in token balances for token0 and token1.
//
// Solidity Equivalent:
//
//	type BalanceDelta is int256;
//	// upper 128 bits = amount0
//	// lower 128 bits = amount1
//
// In the EVM implementation, two signed int128 values are bit-packed into
// a single int256 to save storage gas.
//
// In Go, we explicitly store the two components separately for clarity,
// safety, and maintainability.
//
// IMPORTANT:
// All values are constrained to the signed int128 range to guarantee
// byte-level compatibility with Solidity behavior.
type BalanceDelta struct {
	Amount0 *big.Int // delta of token0
	Amount1 *big.Int // delta of token1
}

// ZeroDelta represents a zero balance change.
var ZeroBalanceDelta = BalanceDelta{
	Amount0: big.NewInt(0),
	Amount1: big.NewInt(0),
}

// NewBalanceDelta constructs a new BalanceDelta while enforcing int128 bounds.
//
// Returns ErrInt128Overflow if either amount exceeds the signed int128 range.
//
// This ensures strict compatibility with Solidity's BalanceDelta type.
func NewBalanceDelta(amount0, amount1 *big.Int) (BalanceDelta, error) {
	if !fitsInt128(amount0) || !fitsInt128(amount1) {
		return ZeroBalanceDelta, ErrInt128Overflow
	}

	return BalanceDelta{
		Amount0: new(big.Int).Set(amount0),
		Amount1: new(big.Int).Set(amount1),
	}, nil
}

// Add returns the element-wise sum of two BalanceDelta values.
//
// Equivalent to Solidity's `add(BalanceDelta a, BalanceDelta b)`.
//
// The result is validated against int128 bounds.
func (a BalanceDelta) Add(b BalanceDelta) (BalanceDelta, error) {
	res0 := new(big.Int).Add(a.Amount0, b.Amount0)
	res1 := new(big.Int).Add(a.Amount1, b.Amount1)

	return NewBalanceDelta(res0, res1)
}

// Sub returns the element-wise difference of two BalanceDelta values.
//
// Equivalent to Solidity's `sub(BalanceDelta a, BalanceDelta b)`.
//
// The result is validated against int128 bounds.
func (a BalanceDelta) Sub(b BalanceDelta) (BalanceDelta, error) {
	res0 := new(big.Int).Sub(a.Amount0, b.Amount0)
	res1 := new(big.Int).Sub(a.Amount1, b.Amount1)

	return NewBalanceDelta(res0, res1)
}

// Equal returns true if both token deltas are equal.
func (a BalanceDelta) Equal(b BalanceDelta) bool {
	return a.Amount0.Cmp(b.Amount0) == 0 &&
		a.Amount1.Cmp(b.Amount1) == 0
}

// fitsInt128 returns true if x is within the signed int128 range.
//
// Range:
//   [-2^127, 2^127 - 1]
//
// This mirrors Solidity's int128 overflow semantics.
func fitsInt128(x *big.Int) bool {
	return x.Cmp(int128Max) <= 0 &&
		x.Cmp(int128Min) >= 0
}