package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/vitaliy-ukiru/test-bank/pkg/logging"
)

type BankAccount interface {
	Deposit(amount float64) error
	Withdraw(amount float64) error
	GetBalance() float64
}

type AccountProcessFunc func(account BankAccount) error

type Acquirer interface {
	Acquire(ctx context.Context, id int64, fn AccountProcessFunc) error
}

type Repository interface {
	NewAccount(ctx context.Context) (int64, error)
}

type AccountService struct {
	locker Acquirer
	repo   Repository
}

func NewAccountService(locker Acquirer, repo Repository) *AccountService {
	return &AccountService{locker: locker, repo: repo}
}

var ErrAccountNotFound = errors.New("account not found")

func (a *AccountService) CreateAccount(ctx context.Context) (accountId int64, err error) {
	const op = "CreateAccount"
	log := logging.FromContext(ctx)

	defer func() {
		if err != nil {
			log.Error(op, "fail create account", err)
			err = fmt.Errorf("%s: %w", op, err)
		} else {
			log.Info(op, "account created", logging.AccountId(accountId))
		}
	}()

	accountId, err = a.repo.NewAccount(ctx)
	if err != nil {
		return
	}
	return accountId, nil
}

func (a *AccountService) DepositBalance(ctx context.Context, cmd DepositBalanceCommand) (err error) {
	const op = "DepositBalance"
	log := logging.FromContext(ctx).With(logging.AccountId(cmd.AccountId))
	defer func() {
		if err != nil {
			log.Error(op, "fail deposit account", err)
			err = fmt.Errorf("%s: %w", op, err)
		} else {
			log.Info(op, "success deposit account")
		}

	}()

	err = a.locker.Acquire(ctx, cmd.AccountId, func(account BankAccount) error {
		time.Sleep(time.Second)
		return account.Deposit(cmd.Amount)
	})
	return

}

func (a *AccountService) WithdrawBalance(ctx context.Context, cmd WithdrawBalanceCommand) (err error) {
	const op = "WithdrawBalance"
	log := logging.FromContext(ctx).With(logging.AccountId(cmd.AccountId))
	defer func() {
		if err != nil {
			log.Error(op, "fail withdraw account", err, logging.AccountId(cmd.AccountId))
			err = fmt.Errorf("%s: %w", op, err)
		} else {
			log.Info(op, "success withdraw account", logging.AccountId(cmd.AccountId))
		}

	}()

	err = a.locker.Acquire(ctx, cmd.AccountId, func(account BankAccount) error {
		return account.Withdraw(cmd.Amount)
	})
	return
}

func (a *AccountService) GetBalance(ctx context.Context, cmd GetBalanceCommand) (balance float64, err error) {
	const op = "GetBalance"
	log := logging.FromContext(ctx).With(logging.AccountId(cmd.AccountId))
	defer func() {
		if err != nil {
			log.Error(op, "fail get balance account", err)
			err = fmt.Errorf("%s: %w", op, err)
		} else {
			log.Info(op, "success get balance account")
		}

	}()

	err = a.locker.Acquire(ctx, cmd.AccountId, func(account BankAccount) error {
		balance = account.GetBalance()
		return nil
	})
	return

}
