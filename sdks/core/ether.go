package sdkcore

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

// Ether represents the native Ether currency on a specific Ethereum chain.
// It embeds Currency and provides a singleton instance per chain ID.
//
// Ether is used to represent the native currency (ETH) for mainnet and testnets,
// and provides helper methods such as Equals to compare currencies.
type Ether struct {
	*Currency
}

// etherCache stores a singleton Ether instance per chain ID.
var etherCache = make(map[uint]*Ether)
var etherCacheMu sync.RWMutex

// newEther creates a new Ether instance for the given chain ID.
//
// This is an internal constructor. Users should use OnChain() to obtain
// the singleton instance for a given chain.
func newEther(chainId uint) *Ether {
	return &Ether{
		Currency: NewCurrency(chainId, common.Address{}, 18, "ETH", "Ether"),
	}
}

// OnChain returns the singleton Ether instance for a given chain ID.
//
// If the instance does not exist, it will create a new one and store it
// in the etherCache. Subsequent calls with the same chain ID will return
// the same instance.
//
// Thread-safe: uses a read-write mutex to guard concurrent access.
func OnChain(chainId uint) *Ether {
	etherCacheMu.RLock()
	e, ok := etherCache[chainId]
	etherCacheMu.RUnlock()
	if ok {
		return e
	}

	etherCacheMu.Lock()
	defer etherCacheMu.Unlock()

	// Double-check to ensure another goroutine hasn't already created it
	if e, ok := etherCache[chainId]; ok {
		return e
	}

	e = newEther(chainId)
	etherCache[chainId] = e
	return e
}

// Equals returns true if the given Currency is the same native Ether
// on the same chain.
//
// This allows comparison between an Ether instance and other Currency objects.
func (e *Ether) Equals(other *Currency) bool {
	return other.IsNative() && other.ChainId() == e.ChainId()
}
