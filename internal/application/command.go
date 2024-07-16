package application

type GetBalanceCommand struct {
	AccountId int64
}

type DepositBalanceCommand struct {
	AccountId int64
	Amount    float64
}

type WithdrawBalanceCommand struct {
	AccountId int64
	Amount    float64
}
