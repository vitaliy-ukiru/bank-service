package request

type DepositRequest struct {
	AccountId int64   `path:"id"`
	Amount    float64 `json:"amount"`
}

type WithdrawRequest struct {
	AccountId int64   `path:"id"`
	Amount    float64 `json:"amount"`
}

type GetBalanceRequest struct {
	AccountId int64 `path:"id"`
}
