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
	return nil, nil
}
