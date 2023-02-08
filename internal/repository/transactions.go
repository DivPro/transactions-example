package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/divpro/transactions-example/pkg/entity"
)

type Transactions struct {
	db *sql.DB
}

func NewTransactions(db *sql.DB) Transactions {
	return Transactions{db: db}
}

func (r Transactions) Create(ctx context.Context, userID, targetID string, amount string) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO transactions (user_id, target_id, amount) VALUES
	   ($1, $2, -$3),
       ($2, $1, $3)
	`,
		userID, targetID, amount)
	if err != nil {
		return fmt.Errorf("create deposit for %s: %w", userID, err)
	}
	_, err = tx.Exec(`
		INSERT INTO balances (user_id, amount)
		VALUES ($1, -$2)
		ON CONFLICT (user_id) DO UPDATE SET amount = EXCLUDED.amount - $2
	`,
		userID, amount)
	if err != nil {
		return fmt.Errorf("create deposit for %s: %w", userID, err)
	}
	_, err = tx.Exec(`
		INSERT INTO balances (user_id, amount)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE SET amount = EXCLUDED.amount + $2
	`,
		targetID, amount)
	if err != nil {
		return fmt.Errorf("create deposit for %s: %w", userID, err)
	}

	return tx.Commit()
}

func (r Transactions) List(ctx context.Context) ([]entity.Transaction, error) {
	return nil, nil
}
