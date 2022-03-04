package model

type Wallet struct {
	Base
	Balance float64 `json:"balance"`
}
