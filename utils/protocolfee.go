package utils

import "github.com/pkg/errors"

// ProtocolFee represents a packed protocol fee for zero-for-one and one-for-zero swaps.
// It's a uint24 value with layout:
//   bits 0-11   -> zero-for-one fee (0..4095, max 1000 allowed)
//   bits 12-23  -> one-for-zero fee (0..4095, max 1000 allowed)
type ProtocolFee uint32

const (
	// MaxProtocolFee is the maximum allowed protocol fee (0.1% = 1000 pips)
	MaxProtocolFee uint16 = 1000

	// Fee thresholds for optimized bounds checking (mirrors Solidity assembly logic)
	fee0Threshold uint32 = 1001
	fee1Threshold uint32 = 1001 << 12

	// PipsDenominator defines the unit (1e6 = 100%)
	PipsDenominator uint32 = 1_000_000

	// Masks
	zeroForOneMask uint32 = 0xfff      // lower 12 bits
	oneForZeroMask uint32 = 0xfff000   // upper 12 bits
)

var (
	ErrInvalidProtocolFee = errors.New("protocol fee exceeds maximum allowed value")
)

// NewProtocolFee packs two fees into ProtocolFee.
//
// Each fee must be <= MaxProtocolFee.
func NewProtocolFee(zeroForOne, oneForZero uint16) (ProtocolFee, error) {
	if zeroForOne > MaxProtocolFee || oneForZero > MaxProtocolFee {
		return 0, ErrInvalidProtocolFee
	}
	return ProtocolFee(uint32(zeroForOne) | (uint32(oneForZero) << 12)), nil
}

// ZeroForOne extracts the zero-for-one fee (lower 12 bits)
func (f ProtocolFee) ZeroForOne() uint16 {
	return uint16(uint32(f) & zeroForOneMask)
}

// OneForZero extracts the one-for-zero fee (upper 12 bits)
func (f ProtocolFee) OneForZero() uint16 {
	return uint16((uint32(f) & oneForZeroMask) >> 12)
}

// IsValid checks whether both fees are <= MaxProtocolFee
func (f ProtocolFee) IsValid() bool {
	raw := uint32(f)
	isZeroForOneOk := (raw & zeroForOneMask) < fee0Threshold
	isOneForZeroOk := (raw & oneForZeroMask) < fee1Threshold
	return isZeroForOneOk && isOneForZeroOk
}

// CalculateSwapFee calculates the combined swap fee for a single direction
//
// swapFee = protocolFee + lpFee - floor(protocolFee * lpFee / 1_000_000)
func CalculateSwapFee(protocolFee uint16, lpFee LPFee) uint32 {
	protocolFeeMasked := uint32(protocolFee & 0xfff) // only lower 12 bits
	lpFeeMasked := lpFee.Value() & 0xffffff                   // uint24 mask

	numerator := protocolFeeMasked * lpFeeMasked
	return protocolFeeMasked + lpFeeMasked - (numerator / PipsDenominator)
}