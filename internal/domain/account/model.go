package account

import "errors"

type Money float64

type Account struct {
	id      int64
	balance float64
}

func NewAccount(id int64, balance float64) Account {
	return Account{id: id, balance: balance}
}

func (a *Account) Id() int64 {
	return a.id
}

var (
	ErrNegativeAmount   = errors.New("negative amount")
	ErrZeroAmount       = errors.New("zero amount")
	ErrNotEnoughBalance = errors.New("not enough balance")
)

func (a *Account) GetBalance() float64 {
	return a.balance
}

func (a *Account) Deposit(amount float64) error {
	if err := a.validateAmount(amount); err != nil {
		return err
	}

	a.balance += amount
	return nil
}

func (a *Account) Withdraw(amount float64) error {
	if err := a.validateAmount(amount); err != nil {
		return err
	}

	if a.balance < amount {
		return ErrNotEnoughBalance
	}
	a.balance -= amount
	return nil
}

func (a *Account) validateAmount(amount float64) error {
	if amount == 0 {
		return ErrZeroAmount
	}

	if amount < 0 {
		return ErrNegativeAmount
	}
	return nil
}
