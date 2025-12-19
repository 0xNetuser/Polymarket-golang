package polymarket

// getContractConfig 获取链的合约配置
func getContractConfig(chainID int, negRisk bool) *ContractConfig {
	// 标准配置
	config := map[int]*ContractConfig{
		137: { // Polygon
			Exchange:         "0x4bFb41d5B3570DeFd03C39a9A4D8dE6Bd8B8982E",
			Collateral:       "0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174",
			ConditionalTokens: "0x4D97DCd97eC945f40cF65F87097ACe5EA0476045",
		},
		80002: { // Amoy
			Exchange:         "0xdFE02Eb6733538f8Ea35D585af8DE5958AD99E40",
			Collateral:       "0x9c4e1703476e875070ee25b56a58b008cfb8fa78",
			ConditionalTokens: "0x69308FB512518e39F9b16112fA8d994F4e2Bf8bB",
		},
	}

	// 负风险配置
	negRiskConfig := map[int]*ContractConfig{
		137: { // Polygon
			Exchange:         "0xC5d563A36AE78145C45a50134d48A1215220f80a",
			Collateral:       "0x2791bca1f2de4661ed88a30c99a7a9449aa84174",
			ConditionalTokens: "0x4D97DCd97eC945f40cF65F87097ACe5EA0476045",
		},
		80002: { // Amoy
			Exchange:         "0xd91E80cF2E7be2e162c6513ceD06f1dD0dA35296",
			Collateral:       "0x9c4e1703476e875070ee25b56a58b008cfb8fa78",
			ConditionalTokens: "0x69308FB512518e39F9b16112fA8d994F4e2Bf8bB",
		},
	}

	if negRisk {
		cfg := negRiskConfig[chainID]
		if cfg == nil {
			panic("Invalid chainID for neg risk: " + string(rune(chainID)))
		}
		return cfg
	}

	cfg := config[chainID]
	if cfg == nil {
		panic("Invalid chainID: " + string(rune(chainID)))
	}
	return cfg
}

