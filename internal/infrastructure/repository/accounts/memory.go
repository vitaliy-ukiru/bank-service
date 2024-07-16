package accounts

import (
	"context"
	"sync"

	"github.com/vitaliy-ukiru/bank-service/internal/application"
	"github.com/vitaliy-ukiru/bank-service/internal/domain/account"
)

type AccountStorage struct {
	rw       sync.RWMutex
	id       int64
	accounts map[int64]float64
}

func NewInMemory() *AccountStorage {
	return &AccountStorage{
		accounts: make(map[int64]float64),
	}
}

func (a *AccountStorage) GetAccountById(ctx context.Context, id int64) (account.Account, error) {
	a.rw.RLock()
	defer a.rw.RUnlock()
	balance, ok := a.accounts[id]
	if !ok {
		return account.Account{}, application.ErrAccountNotFound
	}

	return account.NewAccount(id, balance), nil
}

func (a *AccountStorage) SaveAccount(ctx context.Context, acc account.Account) error {
	a.rw.Lock()
	a.rw.Unlock()
	a.accounts[acc.Id()] = acc.GetBalance()
	return nil
}

func (a *AccountStorage) NewAccount(ctx context.Context) (int64, error) {
	a.rw.Lock()
	a.rw.Unlock()
	a.id++
	id := a.id
	a.accounts[id] = 0
	return id, nil
}
