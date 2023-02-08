package handlers

import (
	"context"
	"time"

	"github.com/divpro/transactions-example/internal/service"
	"github.com/divpro/transactions-example/pkg/entity"
)

type Deposit struct {
	finance service.Finance
	timeout time.Duration
}

func NewDeposit(
	finance service.Finance,
) Deposit {
	return Deposit{
		finance: finance,
		timeout: time.Second * 3,
	}
}

func (h Deposit) Handle(v entity.Deposit) error {
	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	return h.finance.CreateDeposit(ctx, v)
}
