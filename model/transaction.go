package model

type Transaction struct {
	Base
	Amount float64 `json:"amount"`
	Status string  `json:"status"`
}

const Pending = "pending"
const Success = "success"
const Failed = "failed"
