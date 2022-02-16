package accrual

import (
	"context"
	"encoding/json"
	"github.com/magmel48/go-musthave-diploma/internal/logger"
	"go.uber.org/zap"
	"net/http"
	"net/url"
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
	Client  *http.Client
	baseURL string
}

func NewExternalService(client *http.Client, baseURL string) *ExternalService {
	if _, err := url.Parse(baseURL); err != nil {
		// there is assumption about no protocol presence
		baseURL = "http://" + baseURL
	}

	return &ExternalService{
		Client:  client,
		baseURL: baseURL,
	}
}

func (service *ExternalService) GetOrder(ctx context.Context, orderNumber string) (*OrderResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, service.baseURL+"/api/orders/"+orderNumber, nil)
	if err != nil {
		return nil, err
	}

	resp, err := service.Client.Do(req)
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
