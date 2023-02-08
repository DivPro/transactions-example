package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
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

	id := uuid.New().String()
	_, err = tx.Exec("INSERT INTO deposits (id, user_id, amount) VALUES ($1, $2, $3)",
		id, userID, amount)
	if err != nil {
		return fmt.Errorf("create deposit for %s: %w", userID, err)
	}
	_, err = tx.Exec(`
		INSERT INTO balances (user_id, amount)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE SET amount = balances.amount + excluded.amount
	`,
		userID, amount)
	if err != nil {
		return fmt.Errorf("add deposit to balance %s: %w", userID, err)
	}

	return tx.Commit()
}
