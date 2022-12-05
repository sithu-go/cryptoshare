package dto

type ERC20BalanceResp struct {
	USDTBalance float64 `json:"usdt_balance"`
	ETHBalance  float64 `json:"eth_balance"`
}

type ERC20TransStatusResp struct {
	State        int64  `json:"state"`
	StateMessage string `json:"state_message"`
}
