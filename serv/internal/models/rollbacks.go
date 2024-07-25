package models

type Rollback struct {
	Base
	ChainId int64 `json:"chain_id"`
	EventId int64 `json:"event_id"`
}

func (Rollback) TableName() string {
	return "`rollbacks`"
}
