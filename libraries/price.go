package libraries

import (
	"errors"
	"math/big"
)

var (
	// ErrDifferentCurrencies is returned when attempting to operate on two
	// Price or CurrencyAmount values whose currencies do not match the
	// required relationship.
	//
	// Example:
	//   - Multiplying two prices whose base/quote do not align
	//   - Quoting with a mismatched base currency
	ErrDifferentCurrencies = errors.New("different currencies")
)

// Price represents the price of Token0 in terms of Token1.
//
// Conceptually:
//
//   Price = Token1 / Token0
//
// Internally, Price embeds a Fraction where:
//
//   Numerator   = amount of Token1
//   Denominator = amount of Token0
//
// IMPORTANT:
//
//   - Token0 is the base currency (denominator)
//   - Token1 is the quote currency (numerator)
//
// The embedded Fraction stores the raw on-chain ratio.
// The Scalar field adjusts the price according to token decimals
// for proper human-readable formatting.
//
// Mathematical definition:
//
//   rawPrice        = numerator / denominator
//   adjustedPrice   = rawPrice × 10^(decimals(Token0)) / 10^(decimals(Token1))
//
// This mirrors how price scaling works in AMM protocols such as Uniswap.
//
// Price is immutable. All arithmetic operations return new instances.
type Price struct {
	*Fraction

	// Token0 is the base currency (input).
	// This corresponds to the denominator of the price fraction.
	Token0 ICurrency

	// Token1 is the quote currency (output).
	// This corresponds to the numerator of the price fraction.
	Token1 ICurrency

	// Scalar is a precomputed scaling factor used to adjust the raw
	// fraction according to token decimals for display purposes.
	//
	// Scalar = 10^decimals(Token0) / 10^decimals(Token1)
	Scalar *Fraction
}

// NewPrice constructs a new Price instance from a raw fraction.
//
// Parameters:
//   - token0: base currency (denominator)
//   - token1: quote currency (numerator)
//   - denominator: raw denominator amount (Token0)
//   - numerator: raw numerator amount (Token1)
//
// The resulting Price represents:
//
//   numerator / denominator
//
// The Scalar is automatically computed from token decimals.
//
// This constructor does NOT normalize or reduce precision beyond
// Fraction behavior.
func NewPrice(token0, token1 ICurrency, denominator, numerator *big.Int) *Price {
	return &Price{
		Fraction: NewFraction(numerator, denominator),
		Token0:   token0,
		Token1:   token1,
		Scalar: NewFraction(
			new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(token0.Decimals())), nil),
			new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(token1.Decimals())), nil),
		),
	}
}

// Invert returns the reciprocal of the current Price.
//
// Mathematically:
//
//   (Token1 / Token0)^-1 = Token0 / Token1
//
// The base and quote currencies are swapped accordingly.
//
// The original Price is not modified.
func (p *Price) Invert() *Price {
	return NewPrice(p.Token1, p.Token0, p.Numerator, p.Denominator)
}

// Multiply multiplies the current price by another price.
//
// Currency constraint:
//
//   other.Token0 must equal p.Token1
//
// Meaning:
//
//   (A/B) × (B/C) = A/C
//
// If currencies do not align, ErrDifferentCurrencies is returned.
//
// Returns a new Price instance; does not mutate the original.
func (p *Price) Multiply(other *Price) (*Price, error) {
	if !other.Token0.Equal(p.Token1) {
		return nil, ErrDifferentCurrencies
	}

	fraction := p.Fraction.Multiply(other.Fraction)
	return NewPrice(p.Token0, other.Token1, fraction.Denominator, fraction.Numerator), nil
}

// Quote returns the corresponding amount of Token1 for a given amount of Token0.
//
// Equivalent to:
//
//   amountOut = price × amountIn
//
// Currency constraint:
//
//   currencyAmount.Currency must equal p.Token0
//
// Returns ErrDifferentCurrencies if mismatched.
//
// The returned CurrencyAmount is denominated in Token1.
func (p *Price) Quote(currencyAmount *CurrencyAmount) (*CurrencyAmount, error) {
	if !currencyAmount.Currency.Equal(p.Token0) {
		return nil, ErrDifferentCurrencies
	}

	result := p.Fraction.Multiply(currencyAmount.Fraction)
	return newCurrencyAmount(p.Token1, result.Numerator, result.Denominator), nil
}

// adjustedForDecimals returns the price adjusted for token decimals.
//
// This is intended for formatting only.
//
// DO NOT use this value for on-chain or internal mathematical logic.
// Always use the raw Fraction for protocol-level calculations.
func (p *Price) adjustedForDecimals() *Fraction {
	return p.Fraction.Multiply(p.Scalar)
}

// ToSignificant formats the price to a string with a specified number
// of significant digits.
//
// This method applies decimal scaling before formatting.
func (p *Price) ToSignificant(significantDigits int32) string {
	return p.adjustedForDecimals().ToSignificant(significantDigits)
}

// ToFixed formats the price to a fixed number of decimal places.
//
// This method applies decimal scaling before formatting.
func (p *Price) ToFixed(decimalPlaces int32) string {
	return p.adjustedForDecimals().ToFixed(decimalPlaces)
}