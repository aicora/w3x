package libraries

import (
	"math/big"

	"github.com/openchat-im/go-uniswap/utils"
	"github.com/shopspring/decimal"
)

// CurrencyAmount represents an amount of a specific currency with arbitrary precision.
//
// This abstraction mirrors the design of the Uniswap SDK, where:
//
// - All values are stored as exact rational numbers (Fraction)
// - On-chain representation is integer-based (uint256)
// - Human-readable formatting is derived from token decimals
//
// Mathematical Model:
//
// 	rawAmount = numerator / denominator
// 	displayAmount = rawAmount / 10^decimals
//
// Invariants:
//
// 1. Quotient() MUST fit inside uint256.
// 2. DecimalScale = 10^currency.Decimals().
// 3. Arithmetic operations preserve full precision (no rounding).
//
// Security:
//
// Panics if the resulting integer quotient exceeds uint256 max.
// This ensures strict EVM compatibility.
//
// Immutability:
//
// All operations return new CurrencyAmount instances.
// The receiver is never mutated.
type CurrencyAmount struct {
	*Fraction

	// Currency metadata (symbol, decimals, chain info).
	Currency ICurrency

	// DecimalScale is 10^currency.Decimals().
	// Used to convert between raw integer and human-readable representation.
	DecimalScale *big.Int
}

// FromRawAmount constructs a CurrencyAmount from a raw on-chain integer amount.
//
// The raw amount MUST already be scaled by 10^decimals.
//
// Example:
// If token has 18 decimals:
// 	1 token = 1e18 raw units
//
// Equivalent to:
// 	newCurrencyAmount(currency, rawAmount, 1)
func FromRawAmount(currency ICurrency, rawAmount *big.Int) *CurrencyAmount {
	return newCurrencyAmount(currency, rawAmount, big.NewInt(1))
}


// FromFractionalAmount constructs a CurrencyAmount from a rational value.
//
// Used internally during intermediate calculations where
// denominator != 1.
//
// Example:
// (1 / 3) token before final rounding.
func FromFractionalAmount(currency ICurrency, numerator, denominator *big.Int) *CurrencyAmount {
	return newCurrencyAmount(currency, numerator, denominator)
}

// newCurrencyAmount is the internal constructor.
//
// Enforces uint256 upper bound constraint to ensure compatibility
// with Solidity arithmetic.
//
// Panics:
// If the integer quotient exceeds MaxUint256.
func newCurrencyAmount(currency ICurrency, numerator, denominator *big.Int) *CurrencyAmount {
	f := NewFraction(numerator, denominator)

	if f.Quotient().Cmp(utils.MaxUint256) > 0 {
		panic("CURRENCY_AMOUNT_OVERFLOW")
	}

	return &CurrencyAmount{
		Currency:     currency,
		Fraction:     f,
		DecimalScale: new(big.Int).Exp(
			big.NewInt(10),
			big.NewInt(int64(currency.Decimals())),
			nil,
		),
	}
}

// Add returns a new CurrencyAmount equal to:
//
// 	c + other
//
// Requirements:
// - Both amounts MUST reference the same currency.
// - Caller is responsible for currency equality validation.
//
// Precision:
// Exact rational arithmetic. No rounding occurs.
func (c *CurrencyAmount) Add(other *CurrencyAmount) *CurrencyAmount {
	added := c.Fraction.Add(other.Fraction)
	return FromFractionalAmount(c.Currency, added.Numerator, added.Denominator)
}

// Subtract returns:
//
// 	c - other
//
// Negative results are allowed at Fraction level.
// Higher-level logic should validate business constraints.
func (c *CurrencyAmount) Subtract(other *CurrencyAmount) *CurrencyAmount {
	subtracted := c.Fraction.Subtract(other.Fraction)
	return FromFractionalAmount(c.Currency, subtracted.Numerator, subtracted.Denominator)
}

// Multiply multiplies CurrencyAmount by a Fraction.
//
// Typical usage:
// - Applying price ratio
// - Applying slippage tolerance
// - Fee adjustment
//
// Preserves full rational precision.
func (c *CurrencyAmount) Multiply(other *Fraction) *CurrencyAmount {
	multiplied := c.Fraction.Multiply(other)
	return FromFractionalAmount(c.Currency, multiplied.Numerator, multiplied.Denominator)
}

// Divide divides CurrencyAmount by a Fraction.
//
// Typical usage:
// - Converting via price
// - Ratio-based calculations
func (c *CurrencyAmount) Divide(other *Fraction) *CurrencyAmount {
	divided := c.Fraction.Divide(other)
	return FromFractionalAmount(c.Currency, divided.Numerator, divided.Denominator)
}

// ToSignificant returns a string representation with the specified number
// of significant digits.
//
// Process:
// 1. Convert raw integer to display value.
// 2. Apply significant-digit rounding.
//
// Used for UI-friendly display.
func (c *CurrencyAmount) ToSignificant(significantDigits int32) string {
	return c.Fraction.
		Divide(NewFraction(c.DecimalScale, big.NewInt(1))).
		ToSignificant(significantDigits)
}

// ToFixed returns a fixed-decimal string representation.
//
// Panics:
// If decimalPlaces > currency.Decimals().
//
// Example:
// Token decimals = 6
// 	raw = 1234567
// 	ToFixed(2) => "1.23"
func (c *CurrencyAmount) ToFixed(decimalPlaces int32) string {
	if uint(decimalPlaces) > c.Currency.Decimals() {
		panic("DECIMAL_PLACES_EXCEED_TOKEN_DECIMALS")
	}

	return c.Fraction.
		Divide(NewFraction(c.DecimalScale, big.NewInt(1))).
		ToFixed(decimalPlaces)
}

// ToExact returns the fully precise decimal string representation.
//
// Unlike ToFixed or ToSignificant:
// - No rounding unless required by decimal precision.
// - Intended for debugging, logs, or precise export.
//
// Internally uses base-10 decimal arithmetic.
func (c *CurrencyAmount) ToExact() string {
	return decimal.
		NewFromBigInt(c.Quotient(), 0).
		Div(decimal.NewFromBigInt(c.DecimalScale, 0)).
		String()
}