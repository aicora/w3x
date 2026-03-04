package sdkcore

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestEther_OnChainSingleton(t *testing.T) {
	eth1 := OnChain(1)
	eth2 := OnChain(1)

	if eth1 != eth2 {
		t.Errorf("OnChain should return the same instance for the same chainId")
	}

	eth3 := OnChain(3)
	if eth1 == eth3 {
		t.Errorf("OnChain should return different instances for different chainIds")
	}
}

func TestEther_Equals(t *testing.T) {
	ethMain := OnChain(1)
	otherETH := NewCurrency(1, common.Address{}, 18, "ETH", "Ether")
	otherToken := NewCurrency(1, common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"), 18, "WETH", "Wrapped Ether")

	if !ethMain.Equals(otherETH) {
		t.Errorf("Ether.Equals should return true for the same native currency on the same chain")
	}

	if ethMain.Equals(otherToken) {
		t.Errorf("Ether.Equals should return false for a non-native currency")
	}
}
