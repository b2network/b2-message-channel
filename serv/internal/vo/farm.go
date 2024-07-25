package vo

import "github.com/shopspring/decimal"

type DepositRecordsRequest struct {
	//PoolId     int64  `json:"pool_id" form:"pool_id"`
	//StrategyId int64  `json:"strategy_id" form:"strategy_id"`
	Owner    string `json:"owner" form:"owner"`
	Page     int64  `json:"page" form:"page"`
	PageSize int64  `json:"page_size" form:"page_size"`
}

type DepositRecordsResponse struct {
	Code int                        `json:"code"`
	Msg  string                     `json:"msg"`
	Data DepositRecordsResponseData `json:"data"`
}

type DepositRecordsResponseData struct {
	Total       int64                        `json:"total"`
	TotalAmount decimal.Decimal              `json:"total_amount"`
	List        []DepositRecordsResponseList `json:"list"`
}

type DepositRecordsResponseList struct {
	Id         int64           `json:"id"`
	PoolId     int64           `json:"pool_id"`
	StrategyId int64           `json:"strategy_id"`
	Amount     decimal.Decimal `json:"amount"`
	Timestamp  int64           `json:"timestamp"`
	TxHash     string          `json:"tx_hash"`
}

///////////////////////////////////////////////////

type WithdrawRecordsRequest struct {
	//PoolId     int64  `json:"pool_id" form:"pool_id"`
	//StrategyId int64  `json:"strategy_id" form:"strategy_id"`
	Owner    string `json:"owner" form:"owner"`
	Page     int64  `json:"page" form:"page"`
	PageSize int64  `json:"page_size" form:"page_size"`
}

type WithdrawRecordsResponse struct {
	Code int                         `json:"code"`
	Msg  string                      `json:"msg"`
	Data WithdrawRecordsResponseData `json:"data"`
}

type WithdrawRecordsResponseData struct {
	Total       int64                         `json:"total"`
	TotalAmount decimal.Decimal               `json:"total_amount"`
	List        []WithdrawRecordsResponseList `json:"list"`
}

type WithdrawRecordsResponseList struct {
	Id         int64           `json:"id"`
	PoolId     int64           `json:"pool_id"`
	StrategyId int64           `json:"strategy_id"`
	Amount     decimal.Decimal `json:"amount"`
	Timestamp  int64           `json:"timestamp"`
	TxHash     string          `json:"tx_hash"`
}

///////////////////////////////////////////////////

type ClaimRecordsRequest struct {
	PoolId int64  `json:"pool_id" form:"pool_id"`
	Owner  string `json:"owner" form:"owner"`
}

type ClaimRecordsResponse struct {
	Code int                      `json:"code"`
	Msg  string                   `json:"msg"`
	Data ClaimRecordsResponseData `json:"data"`
}

type ClaimRecordsResponseData struct {
	List []ClaimRecordsResponseList `json:"list"`
}

type ClaimRecordsResponseList struct {
	Id         int64           `json:"id"`
	PoolId     int64           `json:"pool_id"`
	StrategyId int64           `json:"strategy_id"`
	Amount     decimal.Decimal `json:"amount"`
	Timestamp  int64           `json:"timestamp"`
	TxHash     string          `json:"tx_hash"`
}

///////////////////////////////////////////////////

type FarmStatsRequest struct {
	Owner string `json:"owner" form:"owner"`
}

type FarmStatsResponse struct {
	Code int                   `json:"code"`
	Msg  string                `json:"msg"`
	Data FarmStatsResponseData `json:"data"`
}

type FarmStatsResponseData struct {
	TotalDeposited decimal.Decimal         `json:"total_deposited"`
	TotalUserCount int64                   `json:"total_user_count"`
	List           []FarmStatsResponseList `json:"list"`
}

type FarmStatsResponseList struct {
	PoolId           int64           `json:"pool_id"`
	StrategyId       int64           `json:"strategy_id"`
	TotalDeposited   decimal.Decimal `json:"total_deposited"`
	CurrentDeposited decimal.Decimal `json:"current_deposited"`
	TVL              decimal.Decimal `json:"tvl"`
	TotalWithdraw    decimal.Decimal `json:"total_withdraw"`
	RequestWithdraw  decimal.Decimal `json:"request_withdraw"`
	TotalUserCount   int64           `json:"total_user_count"`
	CurrentUserCount int64           `json:"current_user_count"`
}

///////////////////////////////////////////////////

type UserDetailsRequest struct {
	Owner string `json:"owner" form:"owner"`
}

type UserDetailsResponse struct {
	Code int                     `json:"code"`
	Msg  string                  `json:"msg"`
	Data UserDetailsResponseData `json:"data"`
}

type UserDetailsResponseData struct {
	List []UserDetailsResponseList `json:"list"`
}

type UserDetailsResponseList struct {
	PoolId           int64           `json:"pool_id"`
	StrategyId       int64           `json:"strategy_id"`
	TotalDeposited   decimal.Decimal `json:"total_deposited"`
	CurrentDeposited decimal.Decimal `json:"current_deposited"`
	RequestWithdraw  decimal.Decimal `json:"request_withdraw"`
	TotalWithdraw    decimal.Decimal `json:"total_withdraw"`
}

///////////////////////////////////////////////////

type TransactionRecordsRequest struct {
	Owner    string `json:"owner" form:"owner"`
	Page     int64  `json:"page" form:"page"`
	PageSize int64  `json:"page_size" form:"page_size"`
}

type TransactionRecordsResponse struct {
	Code int                            `json:"code"`
	Msg  string                         `json:"msg"`
	Data TransactionRecordsResponseData `json:"data"`
}

type TransactionRecordsResponseData struct {
	Total int64                            `json:"total"`
	List  []TransactionRecordsResponseList `json:"list"`
}

type TransactionRecordsResponseList struct {
	Id         int64           `json:"id"`
	PoolId     int64           `json:"pool_id"`
	StrategyId int64           `json:"strategy_id"`
	Type       string          `json:"type"`
	Amount     decimal.Decimal `json:"amount"`
	Timestamp  int64           `json:"timestamp"`
	TxHash     string          `json:"tx_hash"`
}
