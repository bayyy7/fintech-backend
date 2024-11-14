package model

type Deposit struct {
	ContractID string  `json:"contract_id"`
	DepositoID string  `json:"deposito_id"`
	Name       string  `json:"name"`
	AccountID  string  `json:"account_id"`
	MinMonth   int     `json:"min_month"`
	Amount     int     `json:"amount"`
	Bonus      float64 `json:"bonus"`
}
