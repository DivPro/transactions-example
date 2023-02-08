package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/divpro/transactions-example/pkg/entity"
	"github.com/google/uuid"
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
		INSERT INTO transactions (id, user_id, target_id, amount) VALUES
	   ($5, $1, $2, $3),
       ($6, $2, $1, $4)
	`,
		userID, targetID, "-"+amount, amount, uuid.New().String(), uuid.New().String())
	if err != nil {
		return fmt.Errorf("create transaction for %s: %w", userID, err)
	}
	_, err = tx.Exec(`
		INSERT INTO balances (user_id, amount)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE SET amount = EXCLUDED.amount - $3
	`,
		userID, "-"+amount, amount)
	if err != nil {
		return fmt.Errorf("create transaction for %s: %w", userID, err)
	}
	_, err = tx.Exec(`
		INSERT INTO balances (user_id, amount)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE SET amount = EXCLUDED.amount + $2
	`,
		targetID, amount)
	if err != nil {
		return fmt.Errorf("create transaction for %s: %w", userID, err)
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
