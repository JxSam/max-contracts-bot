package model

type UsersContracts struct {
	ChatID     int64 `json:"user_id"`
	ContractID int64 `json:"contract_id"`
	Active     bool  `json:"active"`
}
