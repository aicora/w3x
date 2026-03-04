package utils

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func mustHex(s string) []byte {
	return common.FromHex(s)
}

func TestComputeZkSyncCreate2Address(t *testing.T) {
	sender := "0x1111111111111111111111111111111111111111"
	bytecodeHash := mustHex("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	salt := mustHex("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	input := mustHex("0x1234")

	addr, err := ComputeZkSyncCreate2Address(
		sender,
		bytecodeHash,
		salt,
		input,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := common.HexToAddress("0x3AE01EE74208A42a98Aebddf2f8dA6E2BB04cf63")

	if addr != expected {
		t.Fatalf("address mismatch\nexpected: %s\ngot:      %s",
			expected.Hex(),
			addr.Hex(),
		)
	}
}

func TestComputeZkSyncCreate2Address_NilInput(t *testing.T) {
	sender := "0x2222222222222222222222222222222222222222"
	bytecodeHash := mustHex("0xcccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc")
	salt := mustHex("0xdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd")

	addr, err := ComputeZkSyncCreate2Address(
		sender,
		bytecodeHash,
		salt,
		nil, // test nil input
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if addr == (common.Address{}) {
		t.Fatalf("should not return zero address")
	}
}