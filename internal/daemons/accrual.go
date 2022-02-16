package daemons

import (
	"context"
	"database/sql"
	"github.com/magmel48/go-musthave-diploma/internal/accrual"
	"github.com/magmel48/go-musthave-diploma/internal/balances"
	"github.com/magmel48/go-musthave-diploma/internal/config"
	"github.com/magmel48/go-musthave-diploma/internal/logger"
	"github.com/magmel48/go-musthave-diploma/internal/orders"
	"go.uber.org/zap"
	"net/http"
)

type AccrualJob interface {
	Start()
}

type ExternalAccrualJob struct {
	ctx      context.Context
	accrual  accrual.Service
	orders   orders.Repository
	balances balances.Repository
}

func NewExternalAccrualJob(ctx context.Context, db *sql.DB) *ExternalAccrualJob {
	return &ExternalAccrualJob{
		ctx:      ctx,
		accrual:  accrual.NewExternalService(&http.Client{}, config.AccrualServiceAddress),
		orders:   orders.NewPostgreSQLRepository(db),
		balances: balances.NewPostgreSQLRepository(db),
	}
}

func (job *ExternalAccrualJob) Start() {
	records, err := job.orders.FindUnprocessedOrders(job.ctx)
	if err != nil {
		logger.Error("orders db fetch error", zap.Error(err))
		return
	}

	for _, order := range records {
		response, err := job.accrual.GetOrder(job.ctx, order.Number)
		if err != nil {
			logger.Error("accrual service call error", zap.Error(err))
			continue
		}

		if order.Status != response.Status {
			logger.Info("status updated", zap.Int64("id", order.ID), zap.String("status", response.Status))

			if rowsAffected, err := job.orders.Update(job.ctx, orders.Order{
				ID:      order.ID,
				Accrual: response.Accrual,
				Status:  response.Status,
			}); err != nil {
				logger.Error("order update error", zap.Error(err))
			} else {
				if rowsAffected == 1 {
					if err = job.balances.Change(job.ctx, order.UserID, response.Accrual); err != nil {
						logger.Error("balance update error", zap.Error(err))
					}
				} else {
					logger.Error("more than one order affected while order status and accrual update")
				}
			}
		}
	}
}
