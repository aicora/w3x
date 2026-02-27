package libraries

import (
	"math/big"
	"testing"

	"github.com/openchat-im/go-uniswap/utils"
)

func TestFromRawAmount(t *testing.T) {
	token := NewCurrency(1, "0xABCDEF1234567890abcdef1234567890abcdef12", 18, "TKN", "Token")

	raw := big.NewInt(1e18)

	amount := FromRawAmount(token, raw)

	if amount.Quotient().Cmp(raw) != 0 {
		t.Fatal("raw amount mismatch")
	}

	if amount.ToExact() != "1" {
		t.Fatalf("expected 1, got %s", amount.ToExact())
	}
}

func TestOverflowPanics(t *testing.T) {
	token := NewCurrency(1, "0xABCDEF1234567890abcdef1234567890abcdef12", 18, "TKN", "Token")

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for overflow")
		}
	}()

	overflow := new(big.Int).Add(utils.MaxUint256, big.NewInt(1))
	FromRawAmount(token, overflow)
}

func TestAdd(t *testing.T) {
	token := NewCurrency(1, "0xABCDEF1234567890abcdef1234567890abcdef12", 6, "TKN", "Token")

	a := FromRawAmount(token, big.NewInt(1_000_000)) // 1
	b := FromRawAmount(token, big.NewInt(2_000_000)) // 2

	result := a.Add(b)

	if result.ToExact() != "3" {
		t.Fatalf("expected 3, got %s", result.ToExact())
	}
}

func TestSubtract(t *testing.T) {
	token := NewCurrency(1, "0xABCDEF1234567890abcdef1234567890abcdef12", 6, "TKN", "Token")

	a := FromRawAmount(token, big.NewInt(5_000_000))
	b := FromRawAmount(token, big.NewInt(2_000_000))

	result := a.Subtract(b)

	if result.ToExact() != "3" {
		t.Fatalf("expected 3, got %s", result.ToExact())
	}
}

func TestMultiply(t *testing.T) {
	token := NewCurrency(1, "0xABCDEF1234567890abcdef1234567890abcdef12", 6, "TKN", "Token")

	a := FromRawAmount(token, big.NewInt(1_000_000)) // 1

	fraction := NewFraction(big.NewInt(2), big.NewInt(1))

	result := a.Multiply(fraction)

	if result.ToExact() != "2" {
		t.Fatalf("expected 2, got %s", result.ToExact())
	}
}

func TestDivide(t *testing.T) {
	token := NewCurrency(1, "0xABCDEF1234567890abcdef1234567890abcdef12", 6, "TKN", "Token")

	a := FromRawAmount(token, big.NewInt(2_000_000)) // 2

	fraction := NewFraction(big.NewInt(2), big.NewInt(1))

	result := a.Divide(fraction)

	if result.ToExact() != "1" {
		t.Fatalf("expected 1, got %s", result.ToExact())
	}
}

func TestToFixed(t *testing.T) {
	token := NewCurrency(1, "0xABCDEF1234567890abcdef1234567890abcdef12", 6, "TKN", "Token")

	a := FromRawAmount(token, big.NewInt(1234567)) // 1.234567

	if a.ToFixed(2) != "1.23" {
		t.Fatalf("expected 1.23, got %s", a.ToFixed(2))
	}
}

func TestToSignificant(t *testing.T) {
	token := NewCurrency(1, "0xABCDEF1234567890abcdef1234567890abcdef12", 6, "TKN", "Token")

	a := FromRawAmount(token, big.NewInt(1234567)) // 1.234567

	val := a.ToSignificant(3)

	if val != "1.23" && val != "1.24" {
		t.Fatalf("unexpected significant result: %s", val)
	}
}

func TestToExactHighPrecision(t *testing.T) {
	token := NewCurrency(1, "0xABCDEF1234567890abcdef1234567890abcdef12", 6, "TKN", "Token")

	raw := new(big.Int)
	raw.SetString("1000000000000000001", 10)

	a := FromRawAmount(token, raw)

	expected := "1000000000000.000001"

	if a.ToExact() != expected {
		t.Fatalf("expected %s, got %s", expected, a.ToExact())
	}
}

func TestDeterministicArithmetic(t *testing.T) {
	token := NewCurrency(1, "0xABCDEF1234567890abcdef1234567890abcdef12", 6, "TKN", "Token")

	a := FromRawAmount(token, big.NewInt(1_000_000))
	b := FromRawAmount(token, big.NewInt(2_000_000))

	r1 := a.Add(b)
	r2 := a.Add(b)

	if r1.ToExact() != r2.ToExact() {
		t.Fatal("non-deterministic result")
	}
}