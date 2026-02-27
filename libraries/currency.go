package libraries

import (
	"strings"

	"github.com/pkg/errors"
)

var (
	// ErrInvalidAddressFormat is returned when a currency address is invalid.
	ErrInvalidAddressFormat = errors.New("invalid address format")

	// ErrDifferentChain is returned when two currencies belong to different chains.
	ErrDifferentChain = errors.New("different chain")

	// ErrSameAddress is returned when two currencies have the same address but ordering is required.
	ErrSameAddress = errors.New("same address")
)

var (
	// ZeroAddress is the canonical zero address in Ethereum.
	ZeroAddress = "0x0000000000000000000000000000000000000000"
)

// isNative checks if the given address represents a native token (ETH, BNB, etc.).
//
// Parameters:
//   - address: the token address to check
//
// Returns:
//   - true if the address is empty, 0x0, or the canonical zero address.
func isNative(address string) bool {
	addr := strings.TrimSpace(strings.ToLower(address))

	return addr == "" ||
		addr == "0x0" ||
		addr == ZeroAddress
}

// ICurrency defines the interface for a blockchain currency/token.
type ICurrency interface {
	// IsNative returns true if this currency is the native token of the chain.
	IsNative() bool

	// Address returns the canonical lowercase address of the token.
	Address() string

	// ChainId returns the chain ID of the currency.
	ChainId() uint

	// Decimals returns the number of decimals the token uses.
	Decimals() uint

	// Symbol returns the token symbol (e.g., "ETH", "USDT").
	Symbol() string

	// Name returns the token name (e.g., "Ethereum", "Tether").
	Name() string

	// Equal checks whether two currencies are identical (same chain and address).
	Equal(other ICurrency) bool

	// Lt compares two currencies by address lexicographically.
	// Returns an error if currencies are on different chains or have the same address.
	Lt(other ICurrency) (bool, error)
}

// Currency implements ICurrency and represents a blockchain token or native currency.
type Currency struct {
	isNative bool   // true if the currency is the native token of the chain
	address  string // canonical lowercase token address
	chainId  uint   // chain ID
	decimals uint   // number of token decimals (must be < 255)
	symbol   string // token symbol
	name     string // token name
}

// NewCurrency creates a new Currency instance.
//
// Parameters:
//   - chainID: the blockchain chain ID
//   - address: token address; empty or zero address indicates native token
//   - decimals: number of token decimals (must be < 255)
//   - symbol: token symbol
//   - name: token name
//
// Returns:
//   - ICurrency instance representing the token
//
// Panics:
//   - if decimals >= 255
func NewCurrency(chainID uint, address string, decimals uint, symbol string, name string) ICurrency {
	if decimals >= 255 {
		panic("Token currency decimals must be less than 255")
	}
	return &Currency{
		isNative: isNative(address),
		address:  strings.ToLower(strings.TrimSpace(address)),
		chainId:  chainID,
		decimals: decimals,
		symbol:   symbol,
		name:     name,
	}
}

// IsNative returns true if this currency is the native token of the chain.
func (c *Currency) IsNative() bool {
	return c.isNative
}

// Address returns the canonical lowercase address of the token.
func (c *Currency) Address() string {
	return c.address
}

// ChainId returns the chain ID of the currency.
func (c *Currency) ChainId() uint {
	return c.chainId
}

// Decimals returns the number of decimals used by the token.
func (c *Currency) Decimals() uint {
	return c.decimals
}

// Symbol returns the token symbol (e.g., "ETH").
func (c *Currency) Symbol() string {
	return c.symbol
}

// Name returns the token name (e.g., "Ethereum").
func (c *Currency) Name() string {
	return c.name
}

// Equal checks whether two currencies are identical (same chain ID and address).
func (c *Currency) Equal(other ICurrency) bool {
	return c.chainId == other.ChainId() && c.address == other.Address()
}

// Lt compares two currencies lexicographically by address.
//
// Returns:
//   - true if c < other
//   - ErrDifferentChain if currencies are on different chains
//   - ErrSameAddress if currencies have the same address
func (c *Currency) Lt(other ICurrency) (bool, error) {
	if c.chainId != other.ChainId() {
		return false, ErrDifferentChain
	}
	if c.address == other.Address() {
		return false, ErrSameAddress
	}
	return c.address < other.Address(), nil
}