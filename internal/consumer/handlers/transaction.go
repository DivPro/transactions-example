package handlers

import (
	"context"
	"github.com/divpro/transactions-example/internal/service"
	"github.com/divpro/transactions-example/pkg/entity"
	"time"
)

type Transaction struct {
	finance service.Finance
	timeout time.Duration
}

func NewTransaction(
	finance service.Finance,
) Transaction {
	return Transaction{
		finance: finance,
		timeout: time.Second * 3,
	}
}

func (h Transaction) Handle(v entity.Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	return h.finance.CreateTransaction(ctx, v)
}
