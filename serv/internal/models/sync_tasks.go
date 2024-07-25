package models

const (
	SyncTaskInvalid = "invalid"
	SyncTaskPending = "pending"
	SyncTaskDone    = "done"
)

type SyncTask struct {
	Base
	ChainId     int64  `json:"chain_id"`
	LatestBlock int64  `json:"latest_block"`
	StartBlock  int64  `json:"start_block"`
	EndBlock    int64  `json:"end_block"`
	HandleNum   int64  `json:"handle_num"`
	Contracts   string `json:"contracts"`
	Status      string `json:"status"`
}

func (SyncTask) TableName() string {
	return "`sync_tasks`"
}
