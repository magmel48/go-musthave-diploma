package daemons

import (
	"context"
	"database/sql"
	"github.com/magmel48/go-musthave-diploma/internal/accrual"
	"github.com/magmel48/go-musthave-diploma/internal/config"
	"github.com/magmel48/go-musthave-diploma/internal/logger"
	"github.com/magmel48/go-musthave-diploma/internal/orders"
)

type AccrualJob interface {
	Start()
}

type ExternalAccrualJob struct {
	accrual accrual.Service
	orders  orders.Repository
}

func NewExternalAccrualJob(db *sql.DB) *ExternalAccrualJob {
	return &ExternalAccrualJob{
		accrual: accrual.NewExternalService(config.AccrualServiceAddress),
		orders:  orders.NewRepository(db),
	}
}

func (job *ExternalAccrualJob) Start() {
	records, err := job.orders.FindUnprocessedOrders(context.TODO())
	if err != nil {
		logger.Error(err.Error())
		return
	}

	for _, order := range records {
		response, err := job.accrual.GetOrder(order.Number)
		if err != nil {
			logger.Error(err.Error())
			continue
		}

		if order.Status != response.Status {
			// TODO save changes
		}
	}
}
