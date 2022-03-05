package model

import "errors"

type Wallet struct {
	Base
	Balance float64 `json:"balance"`
}

var ErrInsufficientBalance = errors.New("insufficient balance")
