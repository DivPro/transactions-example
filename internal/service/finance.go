package service

import (
	"context"

	"github.com/divpro/transactions-example/pkg/entity"
)

type Deposits interface {
	Create(ctx context.Context, deposit entity.Deposit) error
}

type Transactions interface {
	Create(ctx context.Context, transaction entity.Transaction) error
	List(ctx context.Context) ([]entity.TransactionView, error)
}

type Users interface {
	ListWithBalance(ctx context.Context) ([]entity.UserView, error)
}

type Finance struct {
	deposits     Deposits
	transactions Transactions
	users        Users
}

func NewFinance(deposits Deposits,
	transactions Transactions,
	users Users,
) Finance {
	return Finance{
		deposits:     deposits,
		transactions: transactions,
		users:        users,
	}
}

func (f Finance) CreateDeposit(ctx context.Context, deposit entity.Deposit) error {
	return f.deposits.Create(ctx, deposit)
}

func (f Finance) CreateTransaction(ctx context.Context, transaction entity.Transaction) error {
	return f.transactions.Create(ctx, transaction)
}

func (f Finance) ListTransactions(ctx context.Context) ([]entity.TransactionView, error) {
	return f.transactions.List(ctx)
}

func (f Finance) ListUsers(ctx context.Context) ([]entity.UserView, error) {
	return f.users.ListWithBalance(ctx)
}
