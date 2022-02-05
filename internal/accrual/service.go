package accrual

import (
	"encoding/json"
	"io/ioutil"
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response OrderResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
