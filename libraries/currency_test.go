package libraries

import (
	"testing"
)

func TestIsNative(t *testing.T) {
	cases := []struct {
		address string
		native  bool
	}{
		{"", true},
		{"0x0", true},
		{ZeroAddress, true},
		{"0x1234567890abcdef1234567890abcdef12345678", false},
	}

	for _, c := range cases {
		if isNative(c.address) != c.native {
			t.Errorf("isNative(%s) = %v; want %v", c.address, !c.native, c.native)
		}
	}
}

func TestNewCurrencyAndGetters(t *testing.T) {
	c := NewCurrency(1, "0xABCDEF1234567890abcdef1234567890abcdef12", 18, "TKN", "Token")

	if c.ChainId() != 1 {
		t.Errorf("expected chainId 1, got %d", c.ChainId())
	}
	if c.Decimals() != 18 {
		t.Errorf("expected decimals 18, got %d", c.Decimals())
	}
	if c.Symbol() != "TKN" {
		t.Errorf("expected symbol TKN, got %s", c.Symbol())
	}
	if c.Name() != "Token" {
		t.Errorf("expected name Token, got %s", c.Name())
	}
	if c.IsNative() {
		t.Errorf("expected non-native currency")
	}
}

func TestEqual(t *testing.T) {
	c1 := NewCurrency(1, "0xabc", 18, "A", "TokenA")
	c2 := NewCurrency(1, "0xabc", 18, "A", "TokenA")
	c3 := NewCurrency(1, "0xdef", 18, "B", "TokenB")
	c4 := NewCurrency(2, "0xabc", 18, "A", "TokenA")

	if !c1.Equal(c2) {
		t.Errorf("expected c1 equal to c2")
	}
	if c1.Equal(c3) {
		t.Errorf("expected c1 not equal to c3")
	}
	if c1.Equal(c4) {
		t.Errorf("expected c1 not equal to c4")
	}
}

func TestLt(t *testing.T) {
	c1 := NewCurrency(1, "0xaaa", 18, "A", "TokenA")
	c2 := NewCurrency(1, "0xbbb", 18, "B", "TokenB")
	c3 := NewCurrency(2, "0xaaa", 18, "A", "TokenA")
	c4 := NewCurrency(1, "0xaaa", 18, "A", "TokenA")

	lt, err := c1.Lt(c2)
	if err != nil || !lt {
		t.Errorf("expected c1 < c2, got %v, err %v", lt, err)
	}

	_, err = c1.Lt(c3)
	if err != ErrDifferentChain {
		t.Errorf("expected ErrDifferentChain, got %v", err)
	}

	_, err = c1.Lt(c4)
	if err != ErrSameAddress {
		t.Errorf("expected ErrSameAddress, got %v", err)
	}
}