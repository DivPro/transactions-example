package repository

import (
	"context"
	"database/sql"

	"github.com/divpro/transactions-example/pkg/entity"
)

type Users struct {
	db *sql.DB
}

func NewUsers(db *sql.DB) Users {
	return Users{db: db}
}

func (r Users) ListWithBalance(ctx context.Context) ([]entity.UserView, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, first_name, last_name, second_name, created_at, b.amount
		FROM users
		LEFT JOIN balances b on users.id = b.user_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []entity.UserView
	for rows.Next() {
		var ent entity.UserView
		if err := rows.Scan(&ent.UserID, &ent.FistName, &ent.LastName, &ent.SecondName, &ent.CreatedAt,
			&ent.Balance); err != nil {
			return nil, err
		}
		res = append(res, ent)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}
