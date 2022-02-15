package accrual

import (
	"encoding/json"
	"github.com/magmel48/go-musthave-diploma/internal/orders"
	"net/http"
)

// OrderResponse represents response from accrual service about order status and reward.
type OrderResponse struct {
	Order   string             `json:"order"`
	Status  orders.OrderStatus `json:"status"`
	Accrual float64            `json:"accrual"`
}

//go:generate mockery --name=Service
type Service interface {
	GetOrder(order string) (*OrderResponse, error)
}

type ExternalService struct {
	baseURL string
}

func NewExternalService(baseURL string) *ExternalService {
	return &ExternalService{baseURL: baseURL}
}

func (service *ExternalService) GetOrder(order string) (*OrderResponse, error) {
	client := http.Client{}
	resp, err := client.Get(service.baseURL + "/api/orders/" + order)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			panic(err)
		}
	}()

	var response OrderResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
