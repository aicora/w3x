package libraries

import (
	"math/big"
	"testing"
)

func TestNewBalanceDelta_Success(t *testing.T) {
	a0 := big.NewInt(100)
	a1 := big.NewInt(-200)

	delta, err := NewBalanceDelta(a0, a1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if delta.Amount0.Cmp(a0) != 0 || delta.Amount1.Cmp(a1) != 0 {
		t.Fatalf("values not set correctly")
	}
}

func TestNewBalanceDelta_Overflow(t *testing.T) {
	overflow := new(big.Int).Add(int128Max, big.NewInt(1))

	_, err := NewBalanceDelta(overflow, big.NewInt(0))
	if err == nil {
		t.Fatalf("expected overflow error")
	}
}

func TestBalanceDelta_Add_Success(t *testing.T) {
	a, _ := NewBalanceDelta(big.NewInt(10), big.NewInt(20))
	b, _ := NewBalanceDelta(big.NewInt(5), big.NewInt(-10))

	res, err := a.Add(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected0 := big.NewInt(15)
	expected1 := big.NewInt(10)

	if res.Amount0.Cmp(expected0) != 0 ||
		res.Amount1.Cmp(expected1) != 0 {
		t.Fatalf("add result incorrect")
	}
}

func TestBalanceDelta_Add_Overflow(t *testing.T) {
	a, _ := NewBalanceDelta(int128Max, big.NewInt(0))
	b, _ := NewBalanceDelta(big.NewInt(1), big.NewInt(0))

	_, err := a.Add(b)
	if err == nil {
		t.Fatalf("expected overflow error")
	}
}

func TestBalanceDelta_Sub_Success(t *testing.T) {
	a, _ := NewBalanceDelta(big.NewInt(20), big.NewInt(10))
	b, _ := NewBalanceDelta(big.NewInt(5), big.NewInt(3))

	res, err := a.Sub(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected0 := big.NewInt(15)
	expected1 := big.NewInt(7)

	if res.Amount0.Cmp(expected0) != 0 ||
		res.Amount1.Cmp(expected1) != 0 {
		t.Fatalf("sub result incorrect")
	}
}

func TestBalanceDelta_Sub_Overflow(t *testing.T) {
	a, _ := NewBalanceDelta(int128Min, big.NewInt(0))
	b, _ := NewBalanceDelta(big.NewInt(1), big.NewInt(0))

	_, err := a.Sub(b)
	if err == nil {
		t.Fatalf("expected overflow error")
	}
}

func TestBalanceDelta_Equal(t *testing.T) {
	a, _ := NewBalanceDelta(big.NewInt(100), big.NewInt(200))
	b, _ := NewBalanceDelta(big.NewInt(100), big.NewInt(200))
	c, _ := NewBalanceDelta(big.NewInt(1), big.NewInt(2))

	if !a.Equal(b) {
		t.Fatalf("expected equal")
	}

	if a.Equal(c) {
		t.Fatalf("expected not equal")
	}
}

func TestInt128BoundaryValues(t *testing.T) {
	// max boundary
	_, err := NewBalanceDelta(int128Max, big.NewInt(0))
	if err != nil {
		t.Fatalf("int128Max should be valid")
	}

	// min boundary
	_, err = NewBalanceDelta(int128Min, big.NewInt(0))
	if err != nil {
		t.Fatalf("int128Min should be valid")
	}
}