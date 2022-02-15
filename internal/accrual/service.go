package accrual

import (
	"context"
	"encoding/json"
	"github.com/magmel48/go-musthave-diploma/internal/logger"
	"go.uber.org/zap"
	"net/http"
)

// OrderResponse represents response from accrual service about order status and reward.
type OrderResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

//go:generate mockery --name=Service
type Service interface {
	GetOrder(ctx context.Context, orderNumber string) (*OrderResponse, error)
}

type ExternalService struct {
	baseURL string
}

func NewExternalService(baseURL string) *ExternalService {
	return &ExternalService{baseURL: baseURL}
}

func (service *ExternalService) GetOrder(ctx context.Context, orderNumber string) (*OrderResponse, error) {
	client := http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, service.baseURL+"/api/orders/"+orderNumber, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			logger.Error("get order accrual close body error", zap.Error(err))
		}
	}()

	var response OrderResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
