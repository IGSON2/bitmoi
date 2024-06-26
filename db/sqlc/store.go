package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	SelectMinMaxTime(interval, name string, c context.Context) (int64, int64, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error)
	SpendTokenTx(ctx context.Context, arg SpendTokenTxParams) (SpendTokenTxResult, error)
	CheckAttendTx(ctx context.Context, arg CheckAttendTxParams) (float64, error)
	AppendPracBalanceTx(ctx context.Context, arg AppendPracBalanceTxParams) error
	SettleImdPracScoreTx(ctx context.Context, arg SettleImdScoreTxParams) (float64, error)
	RewardRecommenderTx(ctx context.Context, arg RewardRecommenderTxParams) (string, error)
}

type SqlStore struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new store
func NewStore(db *sql.DB) Store {
	return &SqlStore{
		Queries: New(db),
		db:      db,
	}
}

// ExecTx executes a function within a database transaction
func (store *SqlStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
