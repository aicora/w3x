package utils

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// ComputeZkSyncCreate2Address derives the contract address produced by the
// zkSync-specific CREATE2 deployment scheme.
//
// Parameters:
//   sender       - Hex-encoded deployer address (0x-prefixed).
//   bytecodeHash - Keccak256 hash of the contract bytecode.
//   salt         - 32-byte CREATE2 salt.
//   input        - Constructor calldata (optional; pass nil for empty input).
//
// Returns:
//   The derived contract address as a checksummed common.Address.
//
// Notes:
//   - This function does not perform validation on salt or bytecodeHash length.
//   - The result is deterministic and fully compatible with zkSync CREATE2 logic.
//   - Behavior is aligned with the official ethers.js implementation.
func ComputeZkSyncCreate2Address(sender string, bytecodeHash []byte, salt []byte, input []byte) (common.Address, error) {
	if input == nil {
		input = []byte{}
	}

	// Compute the zkSync-specific CREATE2 domain separator.
	prefix := crypto.Keccak256([]byte("zksyncCreate2"))

	// Hash constructor input (empty input is valid).
	inputHash := crypto.Keccak256(input)

	// Convert sender to address and left-pad to 32 bytes.
	senderAddr := common.HexToAddress(sender)
	senderPadded := common.LeftPadBytes(senderAddr.Bytes(), 32)

	// Concatenate all components in the required order.
	combined := append(prefix, senderPadded...)
	combined = append(combined, salt...)
	combined = append(combined, bytecodeHash...)
	combined = append(combined, inputHash...)

	// Final keccak256 hash.
	finalHash := crypto.Keccak256(combined)

	// Extract the last 20 bytes (equivalent to JS .slice(26)).
	addressBytes := finalHash[12:]

	return common.BytesToAddress(addressBytes), nil
}