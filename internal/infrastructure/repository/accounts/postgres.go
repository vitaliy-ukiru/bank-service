package accounts

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgtype/pgxtype"
	"github.com/jackc/pgx/v4"
	"github.com/vitaliy-ukiru/bank-service/internal/application"
	"github.com/vitaliy-ukiru/bank-service/internal/domain/account"
)

type Connection interface {
	pgxtype.Querier
	BeginFunc(ctx context.Context, f func(pgx.Tx) error) (err error)
}

type Repository struct {
	conn Connection
}

func NewRepository(conn Connection) *Repository {
	return &Repository{conn: conn}
}

const opPrefix = "repo.Postgres."

func (r *Repository) NewAccount(ctx context.Context) (int64, error) {
	const op = opPrefix + "NewAccount"

	row := r.conn.QueryRow(ctx, `INSERT INTO accounts(balance) VALUES(0) RETURNING id`)
	var accountId int64
	if err := row.Scan(&accountId); err != nil {
		return 0, fmt.Errorf("%s:%w", op, err)
	}
	return accountId, nil
}

func (r *Repository) GetAccountById(ctx context.Context, id int64) (account.Account, error) {
	const op = opPrefix + "GetAccountById"

	row := r.conn.QueryRow(ctx, `SELECT id, balance FROM accounts WHERE id=$1`, id)

	var (
		accountId int64
		balance   float64
	)
	if err := row.Scan(&accountId, balance); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return account.Account{}, fmt.Errorf("%s:%w", op, application.ErrAccountNotFound)
		}

		return account.Account{}, fmt.Errorf("%s:%w", op, err)
	}

	return account.NewAccount(id, balance), nil

}
func (r *Repository) SaveAccount(ctx context.Context, acc account.Account) error {
	const op = opPrefix + "SaveAccount"

	_, err := r.conn.Exec(ctx, `UPDATE accounts SET balance=$2 WHERE id=$1`, acc.Id(), acc.GetBalance())
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	return nil
}

func (r *Repository) Acquire(ctx context.Context, id int64, fn application.AccountProcessFunc) (err error) {
	err = r.conn.BeginFunc(ctx, func(tx pgx.Tx) error {
		wrapped := r.with(tx)
		a, err := wrapped.GetAccountById(ctx, id)
		if err != nil {
			return err
		}
		if err := fn(&a); err != nil {
			return err
		}
		if err := wrapped.SaveAccount(ctx, a); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil

}

func (r *Repository) with(tx pgx.Tx) *Repository {
	return &Repository{conn: tx}
}
