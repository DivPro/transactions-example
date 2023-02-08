package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type Deposits struct {
	db *sql.DB
}

func NewDeposits(db *sql.DB) Deposits {
	return Deposits{db: db}
}

func (r Deposits) Create(ctx context.Context, userID string, amount string) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO deposits (user_id, amount) VALUES ($1, $2)",
		userID, amount)
	if err != nil {
		return fmt.Errorf("create deposit for %s: %w", userID, err)
	}
	_, err = tx.Exec(`
		INSERT INTO balances (user_id, amount)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE SET amount = EXCLUDED.amount + $2;
	`,
		userID, amount)
	if err != nil {
		return fmt.Errorf("add deposit to balance %s: %w", userID, err)
	}

	return tx.Commit()
}
