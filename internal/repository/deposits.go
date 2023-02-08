package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/divpro/transactions-example/pkg/entity"
)

type Deposits struct {
	db *sql.DB
}

func NewDeposits(db *sql.DB) Deposits {
	return Deposits{db: db}
}

func (r Deposits) Create(ctx context.Context, deposit entity.Deposit) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO deposits (id, user_id, amount) VALUES ($1, $2, $3)",
		deposit.ID, deposit.UserID, deposit.Amount)
	if err != nil {
		return fmt.Errorf("create deposit %v: %w", deposit, err)
	}
	_, err = tx.Exec(`
		INSERT INTO balances (user_id, amount)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE SET amount = balances.amount + excluded.amount
	`,
		deposit.UserID, deposit.Amount)
	if err != nil {
		return fmt.Errorf("add deposit to balance %v: %w", deposit, err)
	}

	return tx.Commit()
}
