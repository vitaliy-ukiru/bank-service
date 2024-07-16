package acquire

import (
	"context"
	"sync"

	"github.com/vitaliy-ukiru/bank-service/internal/application"
	"github.com/vitaliy-ukiru/bank-service/internal/domain/account"
)

type AccountStorage interface {
	GetAccountById(ctx context.Context, id int64) (account.Account, error)
	SaveAccount(ctx context.Context, acc account.Account) error
}

type InMemoryAcquirer struct {
	storage AccountStorage

	rw                  sync.RWMutex
	accountInProcessing map[int64]chan account.Account
}

func NewInMemoryAcquirer(storage AccountStorage) *InMemoryAcquirer {
	return &InMemoryAcquirer{
		storage:             storage,
		accountInProcessing: make(map[int64]chan account.Account),
	}
}

func (im *InMemoryAcquirer) addToProcessing(a account.Account) chan account.Account {
	ch := make(chan account.Account, 1) // buffer for stop processing without wait
	im.rw.Lock()
	im.accountInProcessing[a.Id()] = ch
	defer im.rw.Unlock()
	return ch
}

func (im *InMemoryAcquirer) Acquire(ctx context.Context, id int64, fn application.AccountProcessFunc) (err error) {
	im.rw.RLock()
	ch, ok := im.accountInProcessing[id]
	im.rw.RUnlock()
	var a account.Account
	if ok {
		// waiting for free account or exit as done context
		select {
		case a = <-ch:
		case <-ctx.Done():
			return ctx.Err()
		}
	} else {
		// get account from storage
		// and add to processing
		a, err = im.storage.GetAccountById(ctx, id)
		if err != nil {
			return err
		}
		ch = im.addToProcessing(a)
	}

	defer func() { ch <- a }() // push account as free in any way

	err = fn(&a)
	if err != nil {
		return err
	}
	return im.storage.SaveAccount(ctx, a)

}
