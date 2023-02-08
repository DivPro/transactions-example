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

func (r Transactions) Create(ctx context.Context, transaction entity.Transaction) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO transactions (id, user_id, target_id, amount) VALUES
	   ($5, $1, $2, $3),
       ($5, $2, $1, $4)
	`,
		transaction.UserID, transaction.TargetID, "-"+transaction.Amount, transaction.Amount, transaction.ID)
	if err != nil {
		return fmt.Errorf("create transaction %v: %w", transaction, err)
	}
	_, err = tx.Exec(`
		INSERT INTO balances (user_id, amount)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE SET amount = EXCLUDED.amount - $3
	`,
		transaction.UserID, "-"+transaction.Amount, transaction.Amount)
	if err != nil {
		return fmt.Errorf("create transaction %v: %w", transaction, err)
	}
	_, err = tx.Exec(`
		INSERT INTO balances (user_id, amount)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE SET amount = EXCLUDED.amount + $2
	`,
		transaction.TargetID, transaction.Amount)
	if err != nil {
		return fmt.Errorf("create transaction %v: %w", transaction, err)
	}

	return tx.Commit()
}

func (r Transactions) List(ctx context.Context) ([]entity.TransactionView, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, target_id, amount, created_at
		FROM transactions
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []entity.TransactionView
	for rows.Next() {
		var ent entity.TransactionView
		if err := rows.Scan(&ent.ID, &ent.UserID, &ent.TargetID, &ent.Amount, &ent.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, ent)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}
