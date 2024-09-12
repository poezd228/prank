package transactions

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/mvd-inc/anibliss/internal/db"
)

type Repository interface {
	StartTransaction(ctx context.Context) (Transaction, error)
	StartReadOnlyTransaction(ctx context.Context) (Transaction, error)
}

type Transaction interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context)
	Txm() pgx.Tx
}

type Tx struct {
	pgx.Tx
}

type txRepository struct {
	client *db.PostgresClient
}

func NewTxRepository(client *db.PostgresClient) Repository {
	return &txRepository{
		client: client,
	}
}

func (r *txRepository) StartTransaction(ctx context.Context) (Transaction, error) {
	tx, err := r.client.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &Tx{Tx: tx}, nil
}

func (r *txRepository) StartReadOnlyTransaction(ctx context.Context) (Transaction, error) {
	tx, err := r.client.ReadOnlyDB.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		return nil, err
	}
	return &Tx{Tx: tx}, nil
}

func (t *Tx) Rollback(ctx context.Context) {
	_ = t.Tx.Rollback(ctx)
}

func (t *Tx) Txm() pgx.Tx {
	return t.Tx
}
