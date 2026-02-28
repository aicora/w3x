package utils

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// CalculatePositionKey generates a unique 32-byte key for a liquidity position.
//
// The key is calculated using Ethereum's Keccak256 hash function, following the same
// packing rules as Solidity's `keccak256(abi.encodePacked(owner, tickLower, tickUpper, salt))`.
//
// Parameters:
//   - owner: the 20-byte Ethereum address of the position owner (common.Address).
//   - tickLower: the lower tick boundary of the position (int24 in Solidity, passed as int32).
//   - tickUpper: the upper tick boundary of the position (int24 in Solidity, passed as int32).
//   - salt: a 32-byte unique value used to differentiate multiple positions in the same range.
//
// Returns:
//   - A 32-byte array representing the keccak256 hash of the packed input parameters.
//
// Notes:
//   - The owner address must be exactly 20 bytes to match Solidity's address type.
//   - tickLower and tickUpper use only the lower 24 bits to emulate Solidity int24 encoding.
//   - The returned key is fully compatible with Uniswap v4 on-chain position mappings.
func CalculatePositionKey(owner common.Address, tickLower, tickUpper int32, salt [32]byte) [32]byte {
	buf := make([]byte, 0, 58)

	// Append owner: 20 bytes
	buf = append(buf, owner.Bytes()...)

	// Append tickLower: int24 big-endian
	buf = append(buf, byte(tickLower>>16), byte(tickLower>>8), byte(tickLower))

	// Append tickUpper: int24 big-endian
	buf = append(buf, byte(tickUpper>>16), byte(tickUpper>>8), byte(tickUpper))

	// Append salt: 32 bytes
	buf = append(buf, salt[:]...)

	// Compute Keccak256 hash of packed input
	return crypto.Keccak256Hash(buf)
}