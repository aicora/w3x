package libraries

import (
	"math/big"
	"testing"
)


func newETH() ICurrency {
	return NewCurrency(1, "0xABCDEF1234567890abcdef1234567890abcdef13", 18, "ETH", "ETH Token")
}

func newUSDC() ICurrency {
	return NewCurrency(1, "0xABCDEF1234567890abcdef1234567890abcdef12", 6, "USDC", "USDC Token",)
}

func newDAI() ICurrency {
	return NewCurrency(1, "0xABCDEF1234567890abcdef1234567890abcdef11", 18, "DAI", "DAI Token")
}

func TestNewPrice(t *testing.T) {
	token0 := newETH()
	token1 := newUSDC()

	price := NewPrice(
		token0,
		token1,
		big.NewInt(1_000_000),
		big.NewInt(2_000_000_000_000_000_000),
	)

	if !price.Token0.Equal(token0) {
		t.Fatal("token0 mismatch")
	}

	if !price.Token1.Equal(token1) {
		t.Fatal("token1 mismatch")
	}

	if price.Denominator.Cmp(big.NewInt(1_000_000)) != 0 {
		t.Fatal("denominator incorrect")
	}

	if price.Numerator.Cmp(big.NewInt(2_000_000_000_000_000_000)) != 0 {
		t.Fatal("numerator incorrect")
	}
}

func TestPriceInvert(t *testing.T) {
	token0 := newETH()
	token1 := newUSDC()

	price := NewPrice(token0, token1, big.NewInt(2), big.NewInt(10))
	inv := price.Invert()

	if !inv.Token0.Equal(token1) {
		t.Fatal("invert token0 incorrect")
	}

	if !inv.Token1.Equal(token0) {
		t.Fatal("invert token1 incorrect")
	}

	if inv.Numerator.Cmp(big.NewInt(2)) != 0 {
		t.Fatal("invert numerator incorrect")
	}

	if inv.Denominator.Cmp(big.NewInt(10)) != 0 {
		t.Fatal("invert denominator incorrect")
	}
}

func TestPriceMultiply(t *testing.T) {
	ETH := newETH()
	USDC := newUSDC()

	DAI := newDAI()

	// ETH/USDC = 2
	p1 := NewPrice(ETH, USDC, big.NewInt(1), big.NewInt(2))

	// USDC/DAI = 3
	p2 := NewPrice(USDC, DAI, big.NewInt(1), big.NewInt(3))

	result, err := p1.Multiply(p2)
	if err != nil {
		t.Fatal(err)
	}

	// Expect ETH/DAI = 6
	if result.Numerator.Cmp(big.NewInt(6)) != 0 {
		t.Fatal("multiply numerator incorrect")
	}

	if result.Denominator.Cmp(big.NewInt(1)) != 0 {
		t.Fatal("multiply denominator incorrect")
	}
}

func TestPriceMultiplyCurrencyMismatch(t *testing.T) {
	ETH := newETH()
	USDC := newUSDC()

	// Wrong chain
	p1 := NewPrice(ETH, USDC, big.NewInt(1), big.NewInt(2))
	p2 := NewPrice(ETH, USDC, big.NewInt(1), big.NewInt(3))

	_, err := p1.Multiply(p2)
	if err == nil {
		t.Fatal("expected ErrDifferentCurrencies")
	}
}

func TestPriceQuote(t *testing.T) {
	ETH := newETH()
	USDC := newUSDC()

	// ETH/USDC = 2
	price := NewPrice(ETH, USDC, big.NewInt(1), big.NewInt(2))

	amountIn := newCurrencyAmount(ETH, big.NewInt(5), big.NewInt(1))

	out, err := price.Quote(amountIn)
	if err != nil {
		t.Fatal(err)
	}

	// Expect 10 USDC
	if out.Fraction.Numerator.Cmp(big.NewInt(10)) != 0 {
		t.Fatal("quote incorrect")
	}

	if !out.Currency.Equal(USDC) {
		t.Fatal("quote currency incorrect")
	}
}

func TestPriceQuoteCurrencyMismatch(t *testing.T) {
	ETH := newETH()
	USDC := newUSDC()

	price := NewPrice(ETH, USDC, big.NewInt(1), big.NewInt(2))

	amount := newCurrencyAmount(USDC, big.NewInt(1), big.NewInt(1))

	_, err := price.Quote(amount)
	if err == nil {
		t.Fatal("expected ErrDifferentCurrencies")
	}
}

func TestPriceAdjustedForDecimals(t *testing.T) {
	ETH := newETH()
	USDC := newUSDC()

	price := NewPrice(ETH, USDC, big.NewInt(1), big.NewInt(1))

	adj := price.adjustedForDecimals()

	// ETH decimals = 18
	// USDC decimals = 6
	// Scalar = 10^18 / 10^6 = 10^12

	expected18 := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	
	if adj.Numerator.Cmp(expected18) != 0 {
		t.Fatalf("adjustedForDecimals Numerator incorrect, got %s", adj.Numerator.String())
	}

	expected6 := new(big.Int).Exp(big.NewInt(10), big.NewInt(6), nil)

	if adj.Denominator.Cmp(expected6) != 0 {
		t.Fatalf("adjustedForDecimals Denominator incorrect, got %s", adj.Denominator.String())
	}
}

func TestPriceFormatting(t *testing.T) {
	ETH := newETH()
	USDC := newUSDC()

	price := NewPrice(ETH, USDC, big.NewInt(1e18), big.NewInt(2e6))

	fixed := price.ToFixed(2)
	if fixed != "2.00" {
		t.Fatal("ToFixed returned incorrect")
	}

	sig := price.ToSignificant(4)
	if sig != "2" {
		t.Fatal("ToSignificant returned incorrect")
	}

	
}