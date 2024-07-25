package models

type Validator struct {
	Base
	ChainId int64  `json:"chain_id"`
	Address string `json:"address"`
	Status  bool   `json:"status"`
}

func (Validator) TableName() string {
	return "`validators`"
}
