package sdkcore

// ChainId defines supported blockchain network IDs.
type ChainId uint64

const (
	Mainnet         ChainId = 1
	Goerli          ChainId = 5
	Sepolia         ChainId = 11155111
	Optimism        ChainId = 10
	OptimismGoerli  ChainId = 420
	OptimismSepolia ChainId = 11155420
	ArbitrumOne     ChainId = 42161
	ArbitrumGoerli  ChainId = 421613
	ArbitrumSepolia ChainId = 421614
	Polygon         ChainId = 137
	PolygonMumbai   ChainId = 80001
	Celo            ChainId = 42220
	CeloAlfajores   ChainId = 44787
	Gnosis          ChainId = 100
	Moonbeam        ChainId = 1284
	Bnb             ChainId = 56
	Avalanche       ChainId = 43114
	BaseGoerli      ChainId = 84531
	BaseSepolia     ChainId = 84532
	Base            ChainId = 8453
	Zora            ChainId = 7777777
	ZoraSepolia     ChainId = 999999999
	Rootstock       ChainId = 30
	Blast           ChainId = 81457
	Zksync          ChainId = 324
	Worldchain      ChainId = 480
	UniChainSepolia ChainId = 1301
	UniChain        ChainId = 130
	MonadTestnet    ChainId = 10143
	Soneium         ChainId = 1868
	Monad           ChainId = 143
	XLayer          ChainId = 196
	Linea           ChainId = 59144
)

// SupportedChains lists all chains supported by the application.
var SupportedChains = []ChainId{
	Mainnet,
	Optimism,
	OptimismGoerli,
	OptimismSepolia,
	ArbitrumOne,
	ArbitrumGoerli,
	ArbitrumSepolia,
	Polygon,
	PolygonMumbai,
	Goerli,
	Sepolia,
	CeloAlfajores,
	Celo,
	Bnb,
	Avalanche,
	Base,
	BaseGoerli,
	BaseSepolia,
	Zora,
	ZoraSepolia,
	Rootstock,
	Blast,
	Zksync,
	Worldchain,
	UniChainSepolia,
	UniChain,
	MonadTestnet,
	Soneium,
	Monad,
	XLayer,
	Linea,
}

// NativeCurrencyName represents the symbol of native currencies per chain.
type NativeCurrencyName string

const (
	EtherNative     NativeCurrencyName = "ETH"
	MaticNative     NativeCurrencyName = "MATIC"
	CeloNative      NativeCurrencyName = "CELO"
	GnosisNative    NativeCurrencyName = "XDAI"
	MoonbeamNative  NativeCurrencyName = "GLMR"
	BnbNative       NativeCurrencyName = "BNB"
	AvaxNative      NativeCurrencyName = "AVAX"
	RootstockNative NativeCurrencyName = "RBTC"
)
